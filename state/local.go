package state

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"syscall"

	"github.com/nono/cozy-desktop-experiments/state/common"
	"github.com/nono/cozy-desktop-experiments/state/local"
	"github.com/nono/cozy-desktop-experiments/state/remote"
	"github.com/nono/cozy-desktop-experiments/state/types"
)

// CmdStat is a command for making a stat call on a file. It allows to know if
// it is a file or a directory, the size, the inode number, etc.
type CmdStat struct {
	Path string
}

// Exec is required by Command interface.
func (cmd CmdStat) Exec(platform Platform) {
	info, err := platform.FS().Stat(cmd.Path)
	platform.Notify(EventStatDone{Cmd: cmd, Info: info, Error: err})
}

// EventStatDone is an event notified after a stat call was made to send back
// the result.
type EventStatDone struct {
	Cmd   CmdStat
	Info  fs.FileInfo
	Error error
}

// Update is required by Event interface.
func (e EventStatDone) Update(state *State) []Command {
	if e.Error != nil {
		panic(fmt.Errorf("EventStatDone: %s\n", e.Error)) // FIXME
	}
	if e.Cmd.Path == "." && e.Error == nil && e.Info.IsDir() {
		node := state.Nodes.Root()
		node.Ino = getIno(e.Info)
		node.Status = local.ScanningStatus
		state.Nodes.Upsert(node)
		return []Command{CmdScan{"."}}
	}
	return []Command{CmdStop{}} // FIXME
}

// CmdScan is a command to list files and directories inside a directory.
type CmdScan struct {
	Path string
}

// Exec is required by Command interface.
func (cmd CmdScan) Exec(platform Platform) {
	entries, err := platform.FS().ReadDir(cmd.Path)
	platform.Notify(EventScanDone{Cmd: cmd, Path: cmd.Path, Entries: entries, Error: err})
}

// EventScanDone is notified after the scan has been done to send back the
// result, a list of DirEntry.
type EventScanDone struct {
	Cmd     CmdScan
	Path    string
	Entries []fs.DirEntry
	Error   error
}

// Update is required by Event interface.
func (e EventScanDone) Update(state *State) []Command {
	if e.Error != nil {
		panic(fmt.Errorf("EventScanDone: %s\n", e.Error)) // FIXME
	}
	cmds := []Command{}
	parent, err := state.Nodes.ByPath(e.Path)
	if err != nil {
		panic(fmt.Errorf("EventScanDone: %s\n", err)) // FIXME
	}
	for _, entry := range e.Entries {
		node := &local.Node{
			ParentID: parent.ID,
			Name:     entry.Name(),
			Type:     types.FileType,
			Status:   local.InitialStatus,
		}
		if info, err := entry.Info(); err == nil {
			node.Ino = getIno(info)
		}
		if entry.IsDir() {
			node.Status = local.ScanningStatus
			// TODO we should limit the number of scans in parallel
			path := filepath.Join(e.Path, node.Name)
			cmds = append(cmds, CmdScan{path})
			node.Type = types.DirType
		}
		state.Nodes.Upsert(node)
	}
	state.Nodes.MarkAsScanned(parent)
	if len(cmds) > 0 {
		return cmds
	}
	return state.findNextCommand()
}

// CmdMkdir is a command for creating a directory on the local file system.
type CmdMkdir struct {
	Path         string
	RemoteID     remote.ID
	ParentLinkID common.ID
}

// Exec is required by Command interface.
func (cmd CmdMkdir) Exec(platform Platform) {
	var info fs.FileInfo
	localFS := platform.FS()
	err := localFS.Mkdir(cmd.Path)
	if err == nil {
		info, err = localFS.Stat(cmd.Path)
	}
	platform.Notify(EventMkdirDone{Cmd: cmd, Info: info, Error: err})
}

// EventMkdirDone is notified when a directory has been by the desktop client
// on the local file system.
type EventMkdirDone struct {
	Cmd   CmdMkdir
	Info  fs.FileInfo
	Error error
}

// Update is required by Event interface.
func (e EventMkdirDone) Update(state *State) []Command {
	if e.Error != nil {
		panic(fmt.Errorf("EventMkdirDone: %s\n", e.Error)) // FIXME
	}
	parent, err := state.Nodes.ByPath(filepath.Dir(e.Cmd.Path))
	if err != nil || !e.Info.IsDir() {
		panic(fmt.Errorf("EventMkdirDone: %s\n", err)) // FIXME
	}
	node := &local.Node{
		ParentID: parent.ID,
		Name:     e.Info.Name(),
		Type:     types.DirType,
		Status:   local.StableStatus,
		Ino:      getIno(e.Info),
	}
	state.Nodes.Upsert(node)

	link := &common.Link{
		LocalID:  node.ID,
		RemoteID: e.Cmd.RemoteID,
		ParentID: e.Cmd.ParentLinkID,
		Name:     node.Name,
		Type:     types.DirType,
	}
	state.Links.Add(link)
	return state.findNextCommand()
}

// getIno is a small helper function to get the inode number from a fs.FileInfo
// (Linux only).
func getIno(info fs.FileInfo) uint64 {
	return info.Sys().(*syscall.Stat_t).Ino
}
