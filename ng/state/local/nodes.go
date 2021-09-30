package local

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

type State struct {
	ByID            map[ID]*Node
	ByIno           map[uint64]*Node
	ScansInProgress int
}

type Node struct {
	ID       ID
	Ino      uint64 // 0 means unknown
	ParentID ID
	Name     string
	Kind     Kind
}

type ID uint64
type Kind int

const (
	UnknownKind Kind = iota
	FileKind
	DirKind
)

var nextID ID = 2 // 0 = unknown, and 1 is reserved for the root
const rootID ID = 1

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
		Kind:     DirKind,
	}
	state.Upsert(root)
	return state
}

func (state *State) Root() *Node {
	return state.ByID[rootID]
}

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
	fmt.Printf("Upsert %#v\n", n)
}

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
		if n.Kind != FileKind && n.Kind != DirKind {
			return fmt.Errorf("local node %#v should be a file or directory", n)
		}
	}
	return state.checkTree()
}

func (state *State) checkTree() error {
	inTree := map[ID]struct{}{
		rootID: struct{}{},
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
