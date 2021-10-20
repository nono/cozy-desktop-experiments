package state

import (
	"github.com/nono/cozy-desktop-experiments/state/local"
	"github.com/nono/cozy-desktop-experiments/state/remote"
	"github.com/nono/cozy-desktop-experiments/state/types"
)

// State is the global state of the client. It is updated with the information
// from the notified events. And the decisions to which commands to send is
// based on this state.
type State struct {
	Local  *local.State
	Remote *remote.State
	Clock  types.Clock
}

// Sync is the event loop to update the state and send commands, via the
// platform.
func Sync(platform Platform) error {
	state := &State{
		Local:  local.NewState(),
		Remote: remote.NewState(),
	}
	cmds := EventStart{}.Update(state)
	for {
		for _, cmd := range cmds {
			if _, ok := cmd.(CmdStop); ok {
				return state.Local.CheckEventualConsistency()
			}
			platform.Exec(cmd)
		}
		state.Clock++
		event := platform.NextEvent()
		cmds = event.Update(state)
		if err := state.Local.CheckInvariants(); err != nil {
			return err
		}
	}
}
