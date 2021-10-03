package remote

import "github.com/nono/cozy-desktop-experiments/ng/state/types"

type State struct {
	ByID        map[ID]*Doc
	Seq         *Seq
	Refreshing  bool
	RefreshedAt types.Clock
}

type Doc struct {
	ID    ID
	DirID ID
	Type  string
	Name  string
}

type ID string

func NewState() *State {
	return &State{
		ByID: make(map[ID]*Doc),
		Seq:  nil,
	}
}
