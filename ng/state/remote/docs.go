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
	Rev   Rev
	Type  string
	Name  string
	DirID ID
}

type ID string
type Rev string

func NewState() *State {
	return &State{
		ByID: make(map[ID]*Doc),
		Seq:  nil,
	}
}
