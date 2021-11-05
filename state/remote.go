package state

import (
	"fmt"
	"path/filepath"

	"github.com/nono/cozy-desktop-experiments/state/common"
	"github.com/nono/cozy-desktop-experiments/state/remote"
	"github.com/nono/cozy-desktop-experiments/state/types"
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
	state.Docs.Refreshing = false
	state.Docs.RefreshedAt = state.Clock
	state.Docs.FetchingChanges = true
	return []Command{
		CmdChanges{Limit: nbChangesPerPage, Seq: state.Docs.Seq, SkipTrashed: true},
	}
}

// CmdChanges is a command to fetch the changes feed of the Cozy.
type CmdChanges struct {
	Limit       int
	Seq         *remote.Seq
	SkipTrashed bool
}

// Exec is required by Command interface.
func (cmd CmdChanges) Exec(platform Platform) {
	res, err := platform.Client().Changes(cmd.Seq, cmd.Limit, cmd.SkipTrashed)
	if err == nil {
		platform.Notify(EventChangesDone{Docs: res.Docs, Seq: &res.Seq, Pending: res.Pending})
	} else {
		platform.Notify(EventChangesDone{Error: err})
	}
}

// EventChangesDone is used to notify of the result of the changes feed.
type EventChangesDone struct {
	Docs    []*remote.ChangedDoc
	Seq     *remote.Seq
	Pending int
	Error   error
}

// Update is required by Event interface.
func (e EventChangesDone) Update(state *State) []Command {
	// TODO handle error
	state.Docs.Seq = e.Seq
	for _, change := range e.Docs {
		if change.Deleted {
			state.Docs.MarkAsDeleted(change.Doc.ID)
		} else {
			state.Docs.Upsert(change.Doc)
		}
	}
	if e.Pending > 0 {
		return []Command{
			CmdChanges{Limit: nbChangesPerPage, Seq: state.Docs.Seq, SkipTrashed: true},
		}
	}
	state.Docs.FetchingChanges = false
	return state.findNextCommand()
}

// CmdSynchronized is a command to let the Cozy know that the client has reach
// a stable point of synchronization.
type CmdSynchronized struct {
}

// Exec is required by Command interface.
func (cmd CmdSynchronized) Exec(platform Platform) {
	err := platform.Client().Synchronized()
	platform.Notify(EventSynchronized{Error: err})
}

// EventSynchronized is notified when the Cozy has been informed of the
// synchronization, or the call has failed.
type EventSynchronized struct {
	Error error
}

// Update is required by Event interface.
func (e EventSynchronized) Update(state *State) []Command {
	// TODO handle error
	// TODO continuous synchonization
	state.Nodes.PrintTree()
	return []Command{CmdStop{}}
}

const nbChangesPerPage = 10_000

// CmdCreateDir is a command for creating a directory on the Cozy.
type CmdCreateDir struct {
	ParentID remote.ID
	Name     string
}

// Exec is required by Command interface.
func (cmd CmdCreateDir) Exec(platform Platform) {
	doc, err := platform.Client().CreateDir(cmd.ParentID, cmd.Name)
	platform.Notify(EventCreateDirDone{Doc: doc, Error: err})
}

// EventCreateDirDone is notified when a directory has been by the desktop
// client on the Cozy.
type EventCreateDirDone struct {
	Doc   *remote.Doc
	Error error
}

// Update is required by Event interface.
func (e EventCreateDirDone) Update(state *State) []Command {
	state.Docs.Upsert(e.Doc)

	parentLink, ok := state.Links.ByRemoteID[e.Doc.DirID]
	if !ok {
		fmt.Println("TODO") // TODO handle error
	}
	parentNode, ok := state.Nodes.ByID[parentLink.LocalID]
	if !ok {
		fmt.Println("TODO") // TODO handle error
	}
	// FIXME if the node has been renamed/moved while the directory has been
	// created on the Cozy, it won't be found
	parentPath := state.Nodes.Path(parentNode)
	node, err := state.Nodes.ByPath(filepath.Join(parentPath, e.Doc.Name))
	if err != nil {
		fmt.Println("TODO") // TODO handle error
	}
	link := &common.Link{
		LocalID:  node.ID,
		RemoteID: e.Doc.ID,
		ParentID: parentLink.ID,
		Name:     e.Doc.Name,
		Type:     types.DirType,
	}
	state.Links.Add(link)
	return state.findNextCommand()
}
