package remote

type State struct {
	ByID map[ID]*Doc
	Seq  *Seq
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
