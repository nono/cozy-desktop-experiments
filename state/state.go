package state

import (
	"fmt"
	"path/filepath"

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
	if state.Nodes.Root().Status != local.StableStatus || state.Docs.FetchingChanges {
		return []Command{}
	}

	root := state.Links.Root()
	if cmd, ok := state.findMetadataCommand(root); ok {
		return []Command{cmd}
	}

	return []Command{
		CmdSynchronized{Clock: state.Clock},
	}
}

func (state *State) findMetadataCommand(link *common.Link) (Command, bool) {
	node := state.Nodes.ByID[link.LocalID]
	doc := state.Docs.ByID[link.RemoteID]

	if node == nil || doc == nil {
		panic(fmt.Errorf("Unexpected state for %#v", link)) // FIXME
	}

	// TODO compare node & doc
	if link.Type != types.DirType {
		return nil, false
	}

	nodes := state.Nodes.Children(node)
	for _, child := range nodes {
		if _, ok := state.Links.ByLocalID[child.ID]; !ok {
			return &CmdCreateDir{
				ParentID:     doc.ID,
				Name:         child.Name,
				LocalID:      child.ID,
				ParentLinkID: link.ID,
			}, true
		}
	}

	docs := state.Docs.Children(doc)
	if len(docs) != len(nodes) {
		for _, child := range docs {
			if _, ok := state.Links.ByRemoteID[child.ID]; !ok {
				parentPath := state.Nodes.Path(node)
				return &CmdMkdir{
					Path:         filepath.Join(parentPath, child.Name),
					RemoteID:     child.ID,
					ParentLinkID: link.ID,
				}, true
			}
		}
	}

	links := state.Links.Children(link)
	for _, link := range links {
		if cmd, ok := state.findMetadataCommand(link); ok {
			return cmd, ok
		}
	}
	return nil, false
}
