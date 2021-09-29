package state

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"syscall"
)

type LocalState struct {
	ByID            map[LocalID]*LocalNode
	ByIno           map[uint64]*LocalNode
	ScansInProgress int
}

type LocalNode struct {
	ID       LocalID
	Ino      uint64 // 0 means unknown
	ParentID LocalID
	Name     string
	Kind     Kind
}

type LocalID uint64
type Kind int

const (
	UnknownKind Kind = iota
	FileKind
	DirKind
)

var nextLocalID LocalID = 2 // 0 = unknown, and 1 is reserved for the root
const RootID LocalID = 1

func NewLocalState() *LocalState {
	return &LocalState{
		ByID:            make(map[LocalID]*LocalNode),
		ByIno:           make(map[uint64]*LocalNode),
		ScansInProgress: 0,
	}
}

func (nodes *LocalState) Upsert(n *LocalNode) {
	if n.Ino != 0 {
		if was, ok := nodes.ByIno[n.Ino]; ok {
			n.ID = was.ID
		}
		nodes.ByIno[n.Ino] = n
	}
	if n.ID == 0 {
		n.ID = nextLocalID
		nextLocalID++
	}
	nodes.ByID[n.ID] = n
	fmt.Printf("Upsert %#v\n", n)
}

func (nodes *LocalState) Root() *LocalNode {
	return nodes.ByID[RootID]
}

func (nodes *LocalState) ByPath(path string) (*LocalNode, error) {
	parts := strings.Split(path, string(filepath.Separator))
	node := nodes.Root()
	for {
		if node == nil {
			return nil, errors.New("Not found")
		}
		if len(parts) == 0 {
			return node, nil
		}
		part := parts[0]
		parts = parts[1:]
		if part == "." {
			continue
		}
		// TODO optimize me
		parentID := node.ID
		node = nil
		for _, n := range nodes.ByID {
			if n.Name == part && n.ParentID == parentID {
				node = n
			}
		}
	}
}

func (nodes *LocalState) PrintTree(node *LocalNode, indent int) {
	for i := 0; i < indent; i++ {
		fmt.Print("  ")
	}
	fmt.Printf("- %#v\n", node)
	for _, n := range nodes.ByID {
		if n.ParentID == node.ID && n != node { // n != node is needed to avoid looping on the root
			nodes.PrintTree(n, indent+1)
		}
	}
}

type OpStat struct {
	Path string
}

func (o OpStat) Go(platform Platform) {
	go func() {
		info, err := platform.FS().Stat(o.Path)
		platform.Notify(EventStatDone{Op: o, Info: info, Error: err})
	}()
}

type EventStatDone struct {
	Op    OpStat
	Info  fs.FileInfo
	Error error
}

func (e EventStatDone) Update(state *State) []Operation {
	fmt.Printf("Update %#v\n", e.Info)
	if e.Op.Path == "." && e.Error == nil && e.Info.IsDir() {
		node := &LocalNode{ID: RootID, ParentID: RootID, Name: "", Kind: DirKind}
		node.Ino = getIno(e.Info)
		state.local.Upsert(node)
		state.local.ScansInProgress++
		return []Operation{OpScan{"."}}
	}
	return []Operation{OpStop{}}
}

type OpScan struct {
	Path string
}

func (o OpScan) Go(platform Platform) {
	go func() {
		entries, err := platform.FS().ReadDir(o.Path)
		platform.Notify(EventScanDone{Op: o, Path: o.Path, Entries: entries, Error: err})
	}()
}

type EventScanDone struct {
	Op      OpScan
	Path    string
	Entries []fs.DirEntry
	Error   error
}

func (e EventScanDone) Update(state *State) []Operation {
	state.local.ScansInProgress--
	fmt.Printf("Update\n")
	ops := []Operation{}
	var parentID LocalID
	if len(e.Entries) > 0 {
		if parent, err := state.local.ByPath(e.Path); err == nil {
			parentID = parent.ID
		}
	}
	for _, entry := range e.Entries {
		node := &LocalNode{ParentID: parentID, Name: entry.Name(), Kind: FileKind}
		if info, err := entry.Info(); err == nil {
			node.Ino = getIno(info)
		}
		if entry.IsDir() {
			state.local.ScansInProgress++
			path := filepath.Join(e.Path, node.Name)
			ops = append(ops, OpScan{path})
			node.Kind = DirKind
		}
		state.local.Upsert(node)
	}
	if state.local.ScansInProgress == 0 {
		fmt.Printf("---\n")
		state.local.PrintTree(state.local.Root(), 0)
		return []Operation{OpStop{}}
	}
	return ops
}

func getIno(info fs.FileInfo) uint64 {
	return info.Sys().(*syscall.Stat_t).Ino
}
