package state

type EventStart struct{}

func (e EventStart) Update(state *State) []Command {
	state.Remote.Refreshing = true
	return []Command{
		CmdStat{"."},
		CmdRefreshToken{},
	}
}

type CmdStop struct{}

func (cmd CmdStop) Exec(platform Platform) {
	panic("Unreachable code")
}
