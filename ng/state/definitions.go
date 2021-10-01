package state

import (
	"github.com/nono/cozy-desktop-experiments/ng/state/local"
	"github.com/nono/cozy-desktop-experiments/ng/state/remote"
)

type Platform interface {
	FS() local.FS
	Client() remote.Client
	Notify(event Event)
	NextEvent() Event
}

type Event interface {
	Update(state *State) []Operation
}

type Operation interface {
	Go(platform Platform)
}
