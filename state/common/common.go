package common

import (
	"github.com/nono/cozy-desktop-experiments/state/local"
	"github.com/nono/cozy-desktop-experiments/state/remote"
	"github.com/nono/cozy-desktop-experiments/state/types"
)

// Links is keeping the information that links the local file system with the
// remote Cozy.
type Links struct {
	ByID       map[ID]*Link
	ByParentID map[ID]map[ID]*Link // parentID -> map of children
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
	Type     types.Type
}

// ID is a synthetic number for identifying a node.
type ID uint64

// RootID is the identifier for the root.
const RootID ID = 1

// NewLinks creates a new state.
func NewLinks() *Links {
	root := &Link{
		ID:       RootID,
		LocalID:  local.RootID,
		RemoteID: remote.RootID,
		ParentID: RootID,
		Name:     "",
		Type:     types.DirType,
	}
	return &Links{
		ByID:       map[ID]*Link{RootID: root},
		ByParentID: make(map[ID]map[ID]*Link),
	}
}

// Root returns the link between the local synchronized directory and the root
// of the Cozy.
func (links *Links) Root() *Link {
	return links.ByID[RootID]
}

// Children returns a map of ID -> link for the children of the given link.
func (links *Links) Children(parent *Link) map[ID]*Link {
	return links.ByParentID[parent.ID]
}
