// Package state is where all the logic happens. There is a global state that
// is updated when events are received, and in response it gives command to
// execute to the platforms. Those commands can be to change the local
// filesystem, to change the remote Cozy, or for controls (exit for example).
package state

import (
	"github.com/nono/cozy-desktop-experiments/ng/state/local"
	"github.com/nono/cozy-desktop-experiments/ng/state/remote"
)

// Platform is used as a way to plug a local.FS and remote.Client. They can
// execute commands and send events via the platform.
type Platform interface {
	FS() local.FS
	Client() remote.Client
	Notify(event Event)
	NextEvent() Event
	Exec(cmd Command)
}

// Event is a way to give information to update the state when a change happens
// on the local filesystem or the remote Cozy.
type Event interface {
	Update(state *State) []Command
}

// Command is used to ask the local filesystem or the remote Cozy to execute
// things.
type Command interface {
	Exec(platform Platform)
}
