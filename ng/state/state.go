package state

import (
	"github.com/nono/cozy-desktop-experiments/ng/state/local"
	"github.com/nono/cozy-desktop-experiments/ng/state/remote"
	"github.com/nono/cozy-desktop-experiments/ng/state/types"
)

type State struct {
	Local  *local.State
	Remote *remote.State
	Clock  types.Clock
}

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
