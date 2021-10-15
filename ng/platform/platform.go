// Package platform provides a way to plug a local.FS implementation and a
// remote.Client that will be manipulated from the state by commands, and will
// send information via an events channel.
package platform

import (
	"github.com/nono/cozy-desktop-experiments/ng/state"
	"github.com/nono/cozy-desktop-experiments/ng/state/local"
	"github.com/nono/cozy-desktop-experiments/ng/state/remote"
)

// Platform is just a way to plug several components to make them work
// together.
type Platform struct {
	fs     local.FS
	client remote.Client
	events chan state.Event
}

// New returns a new platform with the given local.FS and remote.Client.
func New(fs local.FS, client remote.Client) *Platform {
	return &Platform{
		fs:     fs,
		client: client,
		events: make(chan state.Event),
	}
}

// FS returns the local.FS of the platform.
func (p *Platform) FS() local.FS {
	return p.fs
}

// Client returns the remote.Client of the platform.
func (p *Platform) Client() remote.Client {
	return p.client
}

// Notify sends an event that will be used to update the state.
func (p *Platform) Notify(event state.Event) {
	go func() {
		p.events <- event
	}()
}

// NextEvent fetches the next event to update the state.
func (p *Platform) NextEvent() state.Event {
	return <-p.events
}

// Exec runs the given command in a new goroutine.
func (p *Platform) Exec(cmd state.Command) {
	go func() {
		// TODO recover error
		cmd.Exec(p)
	}()
}
