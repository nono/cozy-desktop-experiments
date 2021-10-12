package state

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"syscall"

	"github.com/nono/cozy-desktop-experiments/ng/state/local"
)

type CmdStat struct {
	Path string
}

func (cmd CmdStat) Exec(platform Platform) {
	info, err := platform.FS().Stat(cmd.Path)
	platform.Notify(EventStatDone{Cmd: cmd, Info: info, Error: err})
}

type EventStatDone struct {
	Cmd   CmdStat
	Info  fs.FileInfo
	Error error
}

func (e EventStatDone) Update(state *State) []Command {
	fmt.Printf("Update %#v\n", e.Info)
	if e.Cmd.Path == "." && e.Error == nil && e.Info.IsDir() {
		node := state.Local.Root()
		node.Ino = getIno(e.Info)
		state.Local.Upsert(node)
		state.Local.ScansInProgress++
		return []Command{CmdScan{"."}}
	}
	return []Command{CmdStop{}}
}

type CmdScan struct {
	Path string
}

func (cmd CmdScan) Exec(platform Platform) {
	entries, err := platform.FS().ReadDir(cmd.Path)
	platform.Notify(EventScanDone{Cmd: cmd, Path: cmd.Path, Entries: entries, Error: err})
}

type EventScanDone struct {
	Cmd     CmdScan
	Path    string
	Entries []fs.DirEntry
	Error   error
}

func (e EventScanDone) Update(state *State) []Command {
	state.Local.ScansInProgress--
	fmt.Printf("Update\n")
	cmds := []Command{}
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
			cmds = append(cmds, CmdScan{path})
			node.Type = local.DirType
		}
		state.Local.Upsert(node)
	}
	if state.Local.ScansInProgress == 0 {
		state.Local.PrintTree()
		return []Command{CmdStop{}}
	}
	return cmds
}

func getIno(info fs.FileInfo) uint64 {
	return info.Sys().(*syscall.Stat_t).Ino
}
