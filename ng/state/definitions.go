package state

import (
	"io/fs"
)

type Platform interface {
	FS() FS
	Notify(event Event)
	NextEvent() Event
}

type FS interface {
	fs.StatFS
	fs.ReadDirFS
}

type Event interface {
	Update(state *State) []Operation
}

type Operation interface {
	Go(platform Platform)
}
