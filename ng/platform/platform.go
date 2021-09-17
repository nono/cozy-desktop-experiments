package platform

import (
	"github.com/nono/cozy-desktop-experiments/ng/state"
)

type Platform struct {
	fs     state.FS
	events chan state.Event
}

func New(fs state.FS) *Platform {
	return &Platform{
		fs:     fs,
		events: make(chan state.Event),
	}
}

func (p *Platform) FS() state.FS {
	return p.fs
}

func (p *Platform) Notify(event state.Event) {
	p.events <- event
}

func (p *Platform) NextEvent() state.Event {
	return <-p.events
}
