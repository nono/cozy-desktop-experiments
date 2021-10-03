package state

type EventStart struct{}

func (e EventStart) Update(state *State) []Operation {
	state.Remote.Refreshing = true
	return []Operation{
		OpStat{"."},
		OpRefreshToken{},
	}
}

type OpStop struct{}

func (o OpStop) Go(platform Platform) {
	panic("Unreachable code")
}
