package state

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"syscall"

	"github.com/nono/cozy-desktop-experiments/ng/state/local"
)

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
		node := state.Local.Root()
		node.Ino = getIno(e.Info)
		state.Local.Upsert(node)
		state.Local.ScansInProgress++
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
	state.Local.ScansInProgress--
	fmt.Printf("Update\n")
	ops := []Operation{}
	var parentID local.ID
	if len(e.Entries) > 0 {
		if parent, err := state.Local.ByPath(e.Path); err == nil {
			parentID = parent.ID
		}
	}
	for _, entry := range e.Entries {
		node := &local.Node{ParentID: parentID, Name: entry.Name(), Type: local.FileType}
		if info, err := entry.Info(); err == nil {
			node.Ino = getIno(info)
		}
		if entry.IsDir() {
			state.Local.ScansInProgress++
			path := filepath.Join(e.Path, node.Name)
			ops = append(ops, OpScan{path})
			node.Type = local.DirType
		}
		state.Local.Upsert(node)
	}
	if state.Local.ScansInProgress == 0 {
		state.Local.PrintTree()
		return []Operation{OpStop{}}
	}
	return ops
}

func getIno(info fs.FileInfo) uint64 {
	return info.Sys().(*syscall.Stat_t).Ino
}
