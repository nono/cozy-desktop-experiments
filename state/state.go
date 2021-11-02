package state

import (
	"fmt"

	"github.com/nono/cozy-desktop-experiments/state/common"
	"github.com/nono/cozy-desktop-experiments/state/local"
	"github.com/nono/cozy-desktop-experiments/state/remote"
	"github.com/nono/cozy-desktop-experiments/state/types"
)

// State is the global state of the client. It is updated with the information
// from the notified events. And the decisions to which commands to send is
// based on this state.
type State struct {
	Links *common.Links
	Nodes *local.Nodes
	Docs  *remote.Docs
	Clock types.Clock
}

// Sync is the event loop to update the state and send commands, via the
// platform.
func Sync(platform Platform) error {
	state := &State{
		Links: common.NewLinks(),
		Nodes: local.NewNodes(),
		Docs:  remote.NewDocs(),
	}
	cmds := EventStart{}.Update(state)
	for {
		for _, cmd := range cmds {
			if _, ok := cmd.(CmdStop); ok {
				return state.Nodes.CheckEventualConsistency()
			}
			platform.Exec(cmd)
		}
		state.Clock++
		event := platform.NextEvent()
		cmds = event.Update(state)
		if err := state.Nodes.CheckInvariants(); err != nil {
			return err
		}
	}
}

func (state *State) findNextCommand() []Command {
	// TODO We should optimize to start synchronizing sooner.
	// For example, we can start uploading a file when fetching changes is done
	// (even if the local scans are still in progress).
	if state.Nodes.ScansInProgress > 0 || state.Docs.FetchingChanges {
		return []Command{}
	}

	root := state.Links.Root()
	if cmd := state.findMetadataCommand(root); cmd != nil {
		return []Command{*cmd}
	}

	return []Command{
		CmdSynchronized{},
	}
}

func (state *State) findMetadataCommand(link *common.Link) *Command {
	node := state.Nodes.ByID[link.LocalID]
	doc := state.Docs.ByID[link.RemoteID]

	switch {
	case node == nil && doc == nil:
		// TODO delete link
		return nil
	case node == nil:
		// TODO delete the doc, and then the link
		return nil
	case doc == nil:
		// TODO delete the node, and then the link
		return nil
	}

	// TODO compare node & doc
	if link.Type != types.DirType {
		return nil
	}

	nodes := state.Nodes.Children(node)
	links := state.Links.Children(link)
	docs := state.Docs.Children(doc)

	// TODO WIP
	fmt.Printf("nodes  = %v\n", nodes)
	fmt.Printf("links = %v\n", links)
	fmt.Printf("docs = %v\n", docs)

	for _, link := range links {
		if cmd := state.findMetadataCommand(link); cmd != nil {
			return cmd
		}
	}
	return nil
}
