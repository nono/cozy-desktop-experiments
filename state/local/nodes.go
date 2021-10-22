package local

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

// State is keeping the information about the files in the synchronized
// directory of the local file system.
type State struct {
	ByID            map[ID]*Node
	ByIno           map[uint64]*Node
	ScansInProgress int
}

// Node is a file or directory on the local file system.
//
// We don't use the inode number as the main identifier, as inode numbers can
// be reused on Linux, and we may want to says that two files (with two inode
// numbers) are the same like when saving a file by using a temporary file.
type Node struct {
	ID       ID
	Ino      uint64 // 0 means unknown
	ParentID ID
	Name     string
	Type     Type
	// TODO executable bit, {c,m,birth}time
}

// ID is a synthetic number for identifying a node.
type ID uint64

// Type is used to differentiate files to directories.
type Type int

const (
	UnknownType Type = iota
	FileType
	DirType
)

var nextID ID = 2 // 0 = unknown, and 1 is reserved for the root
const rootID ID = 1

// NewState creates a new state.
func NewState() *State {
	state := &State{
		ByID:            make(map[ID]*Node),
		ByIno:           make(map[uint64]*Node),
		ScansInProgress: 0,
	}
	root := &Node{
		ID:       rootID,
		ParentID: rootID,
		Name:     "",
		Type:     DirType,
	}
	state.Upsert(root)
	return state
}

// Root returns the node for the root of the synchronized directory.
func (state *State) Root() *Node {
	return state.ByID[rootID]
}

// ByPath returns the node with the given path (if it exists).
func (state *State) ByPath(path string) (*Node, error) {
	parts := strings.Split(path, string(filepath.Separator))
	node := state.Root()
	for {
		if node == nil {
			return nil, errors.New("Not found")
		}
		if len(parts) == 0 {
			return node, nil
		}
		part := parts[0]
		parts = parts[1:]
		if part == "." {
			continue
		}
		// TODO optimize me
		parentID := node.ID
		node = nil
		for _, n := range state.ByID {
			if n.Name == part && n.ParentID == parentID {
				node = n
			}
		}
	}
}

// Upsert will add or update the given node in the state.
func (state *State) Upsert(n *Node) {
	if n.Ino != 0 {
		if was, ok := state.ByIno[n.Ino]; ok {
			n.ID = was.ID
		}
		state.ByIno[n.Ino] = n
	}
	if n.ID == 0 {
		n.ID = nextID
		nextID++
	}
	state.ByID[n.ID] = n
}

// CheckInvariants checks a few properties that should always be true. It can
// be used to detect bugs in the state.Local implementation.
func (state *State) CheckInvariants() error {
	for id, n := range state.ByID {
		if n.ID != id {
			return fmt.Errorf("local node %#v should have id %v", n, id)
		}
	}
	for ino, n := range state.ByIno {
		if n.Ino != ino {
			return fmt.Errorf("local node %#v should have ino %v", n, ino)
		}
	}
	for _, n := range state.ByID {
		if n.ID > 1 && n.ID == n.ParentID {
			return fmt.Errorf("local node %#v should not be its own parent", n)
		}
	}
	return nil
}

// CheckEventualConsistency checks properties that should be true if a stable
// state is reached, ie when no changes are made to the local file system, and
// we wait that the desktop client says it is synchronized. Then, properties
// like all nodes have a type that is file or directory, not unknown, should be
// true.
func (state *State) CheckEventualConsistency() error {
	for _, n := range state.ByID {
		if n.ID <= 0 {
			return fmt.Errorf("local node %#v should have an id > 0", n)
		}
		if n.Ino <= 0 {
			return fmt.Errorf("local node %#v should have an ino > 0", n)
		}
		if n.ParentID <= 0 {
			return fmt.Errorf("local node %#v should have a parentID > 0", n)
		}
		if _, ok := state.ByID[n.ParentID]; !ok {
			return fmt.Errorf("local node %#v should have a parent", n)
		}
		if n.Name == "" && n.ID != rootID {
			return fmt.Errorf("local node %#v should have a name", n)
		}
		if n.Type != FileType && n.Type != DirType {
			return fmt.Errorf("local node %#v should be a file or directory", n)
		}
	}
	return state.checkTree()
}

// checkTree is a naive way to check that we can reach all the nodes starting
// from a root.
func (state *State) checkTree() error {
	inTree := map[ID]struct{}{
		rootID: {},
	}
	nodes := make(map[ID]*Node, len(state.ByID))
	for id, n := range state.ByID {
		nodes[id] = n
	}

	for {
		nb := len(nodes)
		if nb == 0 {
			return nil
		}
		for id, n := range nodes {
			if _, ok := inTree[n.ParentID]; ok {
				inTree[id] = struct{}{}
				delete(nodes, id)
			}
		}
		if nb == len(nodes) {
			return fmt.Errorf("local is not a tree: %v\n", nodes)
		}
	}
}

// PrintTree can be used for debug.
func (state *State) PrintTree() {
	fmt.Printf("---\n")
	state.printTree(state.Root(), 0)
	fmt.Printf("---\n")
}

func (state *State) printTree(node *Node, indent int) {
	for i := 0; i < indent; i++ {
		fmt.Print("  ")
	}
	fmt.Printf("- %#v\n", node)
	for _, n := range state.ByID {
		if n.ParentID == node.ID && n != node { // n != node is needed to avoid looping on the root
			state.printTree(n, indent+1)
		}
	}
}
