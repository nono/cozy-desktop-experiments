package state

type EventStart struct{}

func (e EventStart) Update(state *State) []Operation {
	return []Operation{
		OpStat{"."},
		OpChanges{state.Remote.Seq},
	}
}

type OpStop struct{}

func (o OpStop) Go(platform Platform) {
	panic("Unreachable code")
}
