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
	Exec(cmd Command)
}

type Event interface {
	Update(state *State) []Command
}

type Command interface {
	Exec(platform Platform)
}
