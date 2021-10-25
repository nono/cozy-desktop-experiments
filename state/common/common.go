package common

import (
	"github.com/nono/cozy-desktop-experiments/state/local"
	"github.com/nono/cozy-desktop-experiments/state/remote"
)

// State is keeping the information that links the local file system with the
// remote Cozy.
type State struct {
	ByID map[ID]*Link
}

// Link is the last known state of a file or directory that was common to a
// local node and a remote doc.
//
// Note: I'm using identifiers, and not pointers, for the LocalID and RemoteID
// as I think it will be easier when we will have persistence.
type Link struct {
	ID       ID
	LocalID  local.ID
	RemoteID remote.ID
	ParentID ID
	Name     string
	Type     Type
}

// ID is a synthetic number for identifying a node.
type ID uint64

// Type is used to differentiate files to directories.
type Type int

const (
	FileType Type = iota + 1
	DirType
)

// NewState creates a new state.
func NewState() *State {
	return &State{
		ByID: make(map[ID]*Link),
	}
}
