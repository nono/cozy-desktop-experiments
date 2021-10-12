package state

import (
	"fmt"

	"github.com/nono/cozy-desktop-experiments/ng/state/remote"
)

type CmdRefreshToken struct {
}

func (cmd CmdRefreshToken) Exec(platform Platform) {
	err := platform.Client().Refresh()
	platform.Notify(EventTokenRefreshed{Error: err})
}

type EventTokenRefreshed struct {
	Error error
}

func (e EventTokenRefreshed) Update(state *State) []Command {
	// TODO handle error
	state.Remote.Refreshing = false
	state.Remote.RefreshedAt = state.Clock
	return []Command{
		CmdChanges{state.Remote.Seq},
	}
}

type CmdChanges struct {
	Seq *remote.Seq
}

func (cmd CmdChanges) Exec(platform Platform) {
	res, err := platform.Client().Changes(cmd.Seq)
	if err == nil {
		platform.Notify(EventChangesDone{Docs: res.Docs, Seq: &res.Seq, Pending: res.Pending})
	} else {
		platform.Notify(EventChangesDone{Error: err})
	}
}

type EventChangesDone struct {
	Docs    []*remote.Doc
	Seq     *remote.Seq
	Pending int
	Error   error
}

func (e EventChangesDone) Update(state *State) []Command {
	fmt.Printf("Update %#v\n", e) // TODO
	return []Command{}
}
