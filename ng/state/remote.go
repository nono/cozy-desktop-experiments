package state

import (
	"fmt"

	"github.com/nono/cozy-desktop-experiments/ng/state/remote"
)

type OpRefreshToken struct {
}

func (op OpRefreshToken) Go(platform Platform) {
	go func() {
		err := platform.Client().Refresh()
		platform.Notify(EventRefreshDone{Error: err})
	}()
}

type EventRefreshDone struct {
	Error error
}

func (e EventRefreshDone) Update(state *State) []Operation {
	return []Operation{}
}

type OpChanges struct {
	Seq *remote.Seq
}

func (o OpChanges) Go(platform Platform) {
	go func() {
		res, err := platform.Client().Changes(o.Seq)
		if err == nil {
			platform.Notify(EventChangesDone{Docs: res.Docs, Seq: &res.Seq, Pending: res.Pending})
		} else {
			platform.Notify(EventChangesDone{Error: err})
		}
	}()
}

type EventChangesDone struct {
	Docs    []*remote.Doc
	Seq     *remote.Seq
	Pending int
	Error   error
}

func (e EventChangesDone) Update(state *State) []Operation {
	fmt.Printf("Update %#v\n", e) // TODO
	return []Operation{}
}
