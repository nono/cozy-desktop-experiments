package platform

import (
	"github.com/nono/cozy-desktop-experiments/ng/state"
	"github.com/nono/cozy-desktop-experiments/ng/state/local"
	"github.com/nono/cozy-desktop-experiments/ng/state/remote"
)

type Platform struct {
	fs     local.FS
	client remote.Client
	events chan state.Event
}

func New(fs local.FS, client remote.Client) *Platform {
	return &Platform{
		fs:     fs,
		client: client,
		events: make(chan state.Event),
	}
}

func (p *Platform) FS() local.FS {
	return p.fs
}

func (p *Platform) Client() remote.Client {
	return p.client
}

func (p *Platform) Notify(event state.Event) {
	p.events <- event
}

func (p *Platform) NextEvent() state.Event {
	return <-p.events
}

func (p *Platform) Exec(cmd state.Command) {
	go func() {
		// TODO recover error
		cmd.Exec(p)
	}()
}
