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
	ops := EventStart{}.Update(state)
	for {
		for _, op := range ops {
			if _, ok := op.(OpStop); ok {
				return state.Local.CheckEventualConsistency()
			}
			op.Go(platform)
		}
		state.Clock++
		event := platform.NextEvent()
		ops = event.Update(state)
		if err := state.Local.CheckInvariants(); err != nil {
			return err
		}
	}
}
