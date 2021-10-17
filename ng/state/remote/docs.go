package remote

import "github.com/nono/cozy-desktop-experiments/ng/state/types"

// State is keeping the information about the files on the Cozy.
type State struct {
	ByID        map[ID]*Doc
	Seq         *Seq
	Refreshing  bool
	RefreshedAt types.Clock
}

// Doc describes a CouchDB document for the io.cozy.files doctype on the Cozy.
type Doc struct {
	ID    ID
	Rev   Rev
	Type  string
	Name  string
	DirID ID
}

// ID is used for identifying the CouchDB documents.
type ID string

// Rev is used by CouchDB to avoid conflicts.
type Rev string

// NewState creates a new state.
func NewState() *State {
	return &State{
		ByID: make(map[ID]*Doc),
		Seq:  nil,
	}
}
