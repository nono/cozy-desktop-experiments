package state

// EventStart is used to starts the client.
type EventStart struct{}

// Update is required by Event interface.
func (e EventStart) Update(state *State) []Command {
	state.Remote.Refreshing = true
	return []Command{
		CmdStat{"."},
		CmdRefreshToken{},
	}
}

// CmdStop is a command for stopping the client.
type CmdStop struct{}

// Exec is required by Command interface.
func (cmd CmdStop) Exec(platform Platform) {
	panic("Unreachable code")
}
