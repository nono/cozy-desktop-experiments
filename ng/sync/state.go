package sync

type State struct {
}

func NewState() *State {
	return &State{}
}

func (s *State) Update(event Event) (State, []Operation) {
	return s, nil
}
