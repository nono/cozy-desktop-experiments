package state

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"syscall"
)

type LocalState struct {
	ByID            map[LocalID]*LocalNode
	ByIno           map[uint64]*LocalNode
	ScansInProgress int
}

type LocalNode struct {
	ID   LocalID
	Ino  uint64 // 0 means unknown
	Name string
	Kind Kind
}

type LocalID int
type Kind int

const (
	UnknownKind Kind = iota
	FileKind
	DirKind
)

var nextLocalID LocalID = 0

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
		nextLocalID++
		n.ID = nextLocalID
	}
	nodes.ByID[n.ID] = n
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
		fmt.Printf("inode number = %v\n", e.Info.Sys().(*syscall.Stat_t).Ino)
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
	for _, entry := range e.Entries {
		fmt.Printf("* %#v\n", entry)
		node := &LocalNode{Name: entry.Name(), Kind: FileKind} // TODO ino
		if entry.IsDir() {
			state.local.ScansInProgress++
			path := filepath.Join(e.Path, node.Name)
			ops = append(ops, OpScan{path})
			node.Kind = DirKind
		}
		state.local.Upsert(node)
	}
	if state.local.ScansInProgress == 0 {
		return []Operation{OpStop{}}
	}
	return ops
}
