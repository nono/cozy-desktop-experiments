package client

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/nono/cozy-desktop-experiments/ng/state/remote"
)

type Fake struct {
	Address   string
	SyncCount int
	ByID      map[remote.ID]*remote.Doc
	Feed      []Change

	GenerateID  func() remote.ID
	GenerateRev func(generation int) remote.Rev
	GenerateSeq func(generation int) remote.Seq
}

type Change struct {
	Seq int
	*remote.Doc
	Skip bool
}

func NewFake(address string) remote.Client {
	generateID := func() remote.ID {
		return newUUID()
	}
	generateRev := func(generation int) remote.Rev {
		return newRev(generation)
	}
	generateSeq := func(generation int) remote.Seq {
		// TODO improve it
		return remote.Seq(fmt.Sprintf("%d-seq", generation))
	}

	root := &remote.Doc{
		ID:   remote.RootID,
		Rev:  generateRev(1),
		Type: remote.Directory,
	}
	trash := &remote.Doc{
		ID:    remote.TrashID,
		Rev:   generateRev(1),
		Type:  remote.Directory,
		Name:  remote.TrashName,
		DirID: root.ID,
	}
	byID := map[remote.ID]*remote.Doc{
		root.ID:  root,
		trash.ID: trash,
	}
	fake := &Fake{
		Address:     address,
		SyncCount:   0,
		ByID:        byID,
		Feed:        []Change{},
		GenerateID:  generateID,
		GenerateRev: generateRev,
		GenerateSeq: generateSeq,
	}
	for _, doc := range []*remote.Doc{root, trash} {
		fake.addToChangesFeed(doc)
	}
	// TODO add some design docs to the changes feed

	return fake
}

// TODO Add limit
func (f *Fake) Changes(seq *remote.Seq) (*remote.ChangesResponse, error) {
	since := 0
	if seq != nil {
		since = seq.ExtractGeneration()
	}
	lastSeq := 0
	docs := []*remote.Doc{}
	for _, c := range f.Feed {
		if c.Seq < since {
			continue
		}
		docs = append(docs, c.Doc)
		lastSeq = c.Seq
	}
	return &remote.ChangesResponse{Docs: docs, Seq: f.GenerateSeq(lastSeq), Pending: 0}, nil
}

func (f *Fake) CreateDir(parentID remote.ID, name string) (*remote.Doc, error) {
	// TODO find a way to simulate latency
	if name == "" {
		return nil, errors.New("CreateDir: name is missing")
	}
	if strings.Contains("name", "/") {
		return nil, errors.New("CreateDir: name is invalid")
	}
	if _, ok := f.ByID[parentID]; !ok {
		return nil, errors.New("CreateDir: parent does not exist")
	}

	dir := &remote.Doc{
		ID:    f.GenerateID(),
		Rev:   f.GenerateRev(1),
		Type:  remote.Directory,
		Name:  name,
		DirID: parentID,
	}
	f.ByID[dir.ID] = dir
	f.addToChangesFeed(dir)
	return dir, nil
}

func (f *Fake) Trash(doc *remote.Doc) (*remote.Doc, error) {
	if doc.ID == remote.RootID || doc.ID == remote.TrashID {
		return nil, errors.New("Trash: invalid ID (root or trash)")
	}
	was, ok := f.ByID[doc.ID]
	if !ok {
		return nil, errors.New("Trash: doc not found")
	}
	if was.Rev != doc.Rev {
		return nil, errors.New("Trash: invalid revision")
	}
	was.DirID = remote.TrashID
	was.Rev = f.GenerateRev(extractGeneration(was.Rev) + 1)
	f.addToChangesFeed(was)
	return was, nil
}

func (f *Fake) Refresh() error {
	return nil
}

func (f *Fake) Synchronized() error {
	f.SyncCount++
	return nil
}

func (f *Fake) CheckInvariants() error {
	root, ok := f.ByID[remote.RootID]
	if !ok {
		return errors.New("root is missing")
	}
	if root.Type != remote.Directory {
		return errors.New("root is not a directory")
	}
	trash, ok := f.ByID[remote.TrashID]
	if !ok {
		return errors.New("trash is missing")
	}
	if trash.Type != remote.Directory {
		return errors.New("trash is not a directory")
	}
	if trash.Name != remote.TrashName {
		return errors.New("trash has not the expected name")
	}
	if trash.DirID != root.ID {
		return errors.New("trash has not the expected DirID")
	}

	max := len(f.ByID) + 1
	seen := map[string]*remote.Doc{} // "DirID/Name" -> doc
	for _, doc := range f.ByID {
		if doc.ID == remote.RootID {
			continue
		}
		if err := f.checkCanMoveUpToRoot(doc, max); err != nil {
			return err
		}

		key := fmt.Sprintf("%s/%s", doc.DirID, doc.Name)
		if other, ok := seen[key]; ok {
			return fmt.Errorf("%#v and %#v has same path", doc, other)
		} else {
			seen[key] = doc
		}
	}

	return nil
}

// checkCanMoveUpToRoot ensures that the document is reachable by finding its
// parent, and the parent of its parent, etc. until the root is found. It
// ensures that there is no loop like A is the parent of B and B the parent of
// A.
func (f *Fake) checkCanMoveUpToRoot(doc *remote.Doc, remaining int) error {
	parent, ok := f.ByID[doc.DirID]
	if !ok {
		return fmt.Errorf("%#v parent is missing", doc)
	}
	if parent.Type != remote.Directory {
		return fmt.Errorf("%#v is expected to be a directory", parent)
	}

	if parent.ID == remote.RootID {
		return nil
	} else if remaining == 0 {
		return errors.New("there is a loop")
	}
	return f.checkCanMoveUpToRoot(parent, remaining-1)
}

func (f *Fake) addToChangesFeed(doc *remote.Doc) {
	for i, change := range f.Feed {
		if change.Doc.ID == doc.ID {
			f.Feed[i].Skip = true
		}
	}
	change := Change{
		Seq: len(f.Feed),
		Doc: doc,
	}
	f.Feed = append(f.Feed, change)
}

func newUUID() remote.ID {
	guid := uuid.Must(uuid.NewV4())
	id := fmt.Sprintf("%s", guid)
	id = strings.Replace(id, "-", "", -1)
	return remote.ID(id)
}

func newRev(generation int) remote.Rev {
	rev := fmt.Sprintf("%d-rev", generation) // TODO improve it
	return remote.Rev(rev)
}

func extractGeneration(rev remote.Rev) int {
	parts := strings.Split(string(rev), "-")
	n, err := strconv.Atoi(parts[0])
	if err != nil {
		panic(fmt.Errorf("cannot extract generation from rev %s", rev))
	}
	return n
}
