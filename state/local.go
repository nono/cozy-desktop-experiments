package state

import (
	"io/fs"
	"path/filepath"
	"syscall"

	"github.com/nono/cozy-desktop-experiments/state/local"
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
	if e.Cmd.Path == "." && e.Error == nil && e.Info.IsDir() {
		node := state.Nodes.Root()
		node.Ino = getIno(e.Info)
		node.Status = local.ScanningStatus
		state.Nodes.Upsert(node)
		return []Command{CmdScan{"."}}
	}
	return []Command{CmdStop{}}
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
	cmds := []Command{}
	parent, err := state.Nodes.ByPath(e.Path)
	if err != nil {
		// TODO handle error
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

// getIno is a small helper function to get the inode number from a fs.FileInfo
// (Linux only).
func getIno(info fs.FileInfo) uint64 {
	return info.Sys().(*syscall.Stat_t).Ino
}
