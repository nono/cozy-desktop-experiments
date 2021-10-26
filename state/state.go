package state

import (
	"github.com/nono/cozy-desktop-experiments/state/common"
	"github.com/nono/cozy-desktop-experiments/state/local"
	"github.com/nono/cozy-desktop-experiments/state/remote"
	"github.com/nono/cozy-desktop-experiments/state/types"
)

// State is the global state of the client. It is updated with the information
// from the notified events. And the decisions to which commands to send is
// based on this state.
type State struct {
	Links *common.Links
	Nodes *local.Nodes
	Docs  *remote.Docs
	Clock types.Clock
}

// Sync is the event loop to update the state and send commands, via the
// platform.
func Sync(platform Platform) error {
	state := &State{
		Links: common.NewLinks(),
		Nodes: local.NewNodes(),
		Docs:  remote.NewDocs(),
	}
	cmds := EventStart{}.Update(state)
	for {
		for _, cmd := range cmds {
			if _, ok := cmd.(CmdStop); ok {
				return state.Nodes.CheckEventualConsistency()
			}
			platform.Exec(cmd)
		}
		state.Clock++
		event := platform.NextEvent()
		cmds = event.Update(state)
		if err := state.Nodes.CheckInvariants(); err != nil {
			return err
		}
	}
}
