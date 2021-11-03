package local

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/nono/cozy-desktop-experiments/state/types"
)

// Nodes is keeping the information about the files in the synchronized
// directory of the local file system.
type Nodes struct {
	ByID       map[ID]*Node
	ByParentID map[ID]map[ID]*Node // parentID -> map of children
	ByIno      map[uint64]*Node
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
	Type     types.Type
	Status   Status
	// TODO executable bit, {c,m,birth}time
}

// ID is a synthetic number for identifying a node.
type ID uint64

// RootID is the identifier for the root (ie, the synchronized directory).
const RootID ID = 1

// nextID is the next available identifier to assign to a new node.
var nextID ID = RootID + 1 // 0 is for unknown

// Status describes what operations are in progress on a node to improve the
// knowledge we have of it. For a directory, we can scan it to find what are
// its children.
type Status int

const (
	// InitialStatus is the default status of a new node
	InitialStatus Status = iota
	// ScanningStatus means that we are looking for the direct children of a directory
	ScanningStatus
	// ScannedStatus means that we know the direct children, but not yet the nodes below them
	ScannedStatus
	// StableStatus means that that the information about this node is reliable
	StableStatus
)

// NewNodes creates a new state for managing nodes (data about files or
// directories on the local file system).
func NewNodes() *Nodes {
	nodes := &Nodes{
		ByID:       make(map[ID]*Node),
		ByParentID: make(map[ID]map[ID]*Node),
		ByIno:      make(map[uint64]*Node),
	}
	root := &Node{
		ID:       RootID,
		ParentID: RootID,
		Name:     "",
		Type:     types.DirType,
		Status:   InitialStatus,
	}
	nodes.Upsert(root)
	return nodes
}

// Root returns the node for the root of the synchronized directory.
func (nodes *Nodes) Root() *Node {
	return nodes.ByID[RootID]
}

// ByPath returns the node with the given path (if it exists).
func (nodes *Nodes) ByPath(path string) (*Node, error) {
	parts := strings.Split(path, string(filepath.Separator))
	node := nodes.Root()
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
		parentID := node.ID
		node = nil
		for _, n := range nodes.ByParentID[parentID] {
			if n.Name == part {
				node = n
			}
		}
	}
}

// Upsert will add or update the given node in the nodes.
func (nodes *Nodes) Upsert(n *Node) {
	if n.Ino != 0 {
		if was, ok := nodes.ByIno[n.Ino]; ok {
			n.ID = was.ID
		}
		nodes.ByIno[n.Ino] = n
	}
	if n.ID == 0 {
		n.ID = nextID
		nextID++
	} else {
		nodes.detachParent(n.ID)
	}
	nodes.ByID[n.ID] = n
	if n.ID != RootID {
		nodes.attachParent(n)
	}
}

func (nodes *Nodes) MarkAsScanned(node *Node) {
	node.Status = ScannedStatus
	for _, child := range nodes.Children(node) {
		if child.Type == types.DirType && child.Status != StableStatus {
			return
		}
	}
	node.Status = StableStatus
	if node.ID == RootID {
		return
	}
	parent := nodes.ByID[node.ParentID]
	nodes.MarkAsScanned(parent)
}

// Children returns a map of id -> node for the children of the given
// directory.
func (nodes *Nodes) Children(parent *Node) map[ID]*Node {
	return nodes.ByParentID[parent.ID]
}

// attachParent updates the ByParentID field when a child is added to a
// directory.
func (nodes *Nodes) attachParent(child *Node) {
	children := nodes.ByParentID[child.ParentID]
	if children == nil {
		children = make(map[ID]*Node)
	}
	children[child.ID] = child
	nodes.ByParentID[child.ParentID] = children
}

// attachParent updates the ByParentID field when a child is removed a
// directory.
func (nodes *Nodes) detachParent(childID ID) {
	if was, ok := nodes.ByID[childID]; ok {
		delete(nodes.ByParentID[was.ParentID], was.ID)
	}
}

// CheckInvariants checks a few properties that should always be true. It can
// be used to detect bugs in the nodes.Local implementation.
func (nodes *Nodes) CheckInvariants() error {
	for id, n := range nodes.ByID {
		if n.ID != id {
			return fmt.Errorf("local node %#v should have id %v", n, id)
		}
	}
	for ino, n := range nodes.ByIno {
		if n.Ino != ino {
			return fmt.Errorf("local node %#v should have ino %v", n, ino)
		}
	}
	for _, n := range nodes.ByID {
		if n.ID == RootID {
			continue
		}
		if n.ID == n.ParentID {
			return fmt.Errorf("local node %#v should not be its own parent", n)
		}
		children, ok := nodes.ByParentID[n.ParentID]
		if ok {
			_, ok = children[n.ID]
		}
		if !ok {
			return fmt.Errorf("local node %#v has no parent", n)
		}
	}
	for _, children := range nodes.ByParentID {
		for id, child := range children {
			if n, ok := nodes.ByID[id]; !ok || n != child {
				return fmt.Errorf("invalid indexation ByParentID for %q", id)
			}
		}
	}
	return nil
}

// CheckEventualConsistency checks properties that should be true if a stable
// nodes is reached, ie when no changes are made to the local file system, and
// we wait that the desktop client says it is synchronized. Then, properties
// like all nodes have a type that is file or directory, not unknown, should be
// true.
func (nodes *Nodes) CheckEventualConsistency() error {
	for _, n := range nodes.ByID {
		if n.ID <= 0 {
			return fmt.Errorf("local node %#v should have an id > 0", n)
		}
		if n.Ino <= 0 {
			return fmt.Errorf("local node %#v should have an ino > 0", n)
		}
		if n.ParentID <= 0 {
			return fmt.Errorf("local node %#v should have a parentID > 0", n)
		}
		if _, ok := nodes.ByID[n.ParentID]; !ok {
			return fmt.Errorf("local node %#v should have a parent", n)
		}
		if n.Name == "" && n.ID != RootID {
			return fmt.Errorf("local node %#v should have a name", n)
		}
		if n.Type != types.FileType && n.Type != types.DirType {
			return fmt.Errorf("local node %#v should be a file or directory", n)
		}
	}
	return nodes.checkTree()
}

// checkTree is a naive way to check that we can reach all the nodes starting
// from a root.
func (nodes *Nodes) checkTree() error {
	inTree := map[ID]struct{}{
		RootID: {},
	}
	set := make(map[ID]*Node, len(nodes.ByID))
	for id, n := range nodes.ByID {
		set[id] = n
	}

	for {
		nb := len(set)
		if nb == 0 {
			return nil
		}
		for id, n := range set {
			if _, ok := inTree[n.ParentID]; ok {
				inTree[id] = struct{}{}
				delete(set, id)
			}
		}
		if nb == len(set) {
			return fmt.Errorf("local is not a tree: %v\n", set)
		}
	}
}

// PrintTree can be used for debug.
func (nodes *Nodes) PrintTree() {
	fmt.Printf("---\n")
	nodes.printTree(nodes.Root(), 0)
	fmt.Printf("---\n")
}

func (nodes *Nodes) printTree(node *Node, indent int) {
	for i := 0; i < indent; i++ {
		fmt.Print("  ")
	}
	fmt.Printf("- %#v\n", node)
	for _, n := range nodes.ByID {
		if n.ParentID == node.ID && n != node { // n != node is needed to avoid looping on the root
			nodes.printTree(n, indent+1)
		}
	}
}
