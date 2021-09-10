package sync

type Event interface {
	Update(state *State) (*State, []Operation)
}

type EventStart struct{}

func (e EventStart) Update(state *State) (*State, []Operation) {
	return state, []Operation{OpStop{}}
}
