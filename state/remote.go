package state

import (
	"fmt"

	"github.com/nono/cozy-desktop-experiments/state/remote"
)

// CmdRefreshToken is a command to refresh the access token of the OAuth
// client.
type CmdRefreshToken struct {
}

// Exec is required by Command interface.
func (cmd CmdRefreshToken) Exec(platform Platform) {
	err := platform.Client().Refresh()
	platform.Notify(EventTokenRefreshed{Error: err})
}

// EventTokenRefreshed is notified when the access token has been refreshed, or
// when the attempt has failed.
type EventTokenRefreshed struct {
	Error error
}

// Update is required by Event interface.
func (e EventTokenRefreshed) Update(state *State) []Command {
	// TODO handle error
	state.Remote.Refreshing = false
	state.Remote.RefreshedAt = state.Clock
	return []Command{
		CmdChanges{state.Remote.Seq},
	}
}

// CmdChanges is a command to fetch the changes feed of the Cozy.
type CmdChanges struct {
	Seq *remote.Seq
}

// Exec is required by Command interface.
func (cmd CmdChanges) Exec(platform Platform) {
	res, err := platform.Client().Changes(cmd.Seq)
	if err == nil {
		platform.Notify(EventChangesDone{Docs: res.Docs, Seq: &res.Seq, Pending: res.Pending})
	} else {
		platform.Notify(EventChangesDone{Error: err})
	}
}

// EventChangesDone is used to notify of the result of the changes feed.
type EventChangesDone struct {
	Docs    []*remote.Doc
	Seq     *remote.Seq
	Pending int
	Error   error
}

// Update is required by Event interface.
func (e EventChangesDone) Update(state *State) []Command {
	fmt.Printf("Update %#v\n", e) // TODO
	return []Command{}
}
