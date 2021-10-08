package client

import (
	"errors"
	"fmt"
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
		Type: "directory",
	}
	trash := &remote.Doc{
		ID:    remote.TrashID,
		Rev:   generateRev(1),
		Type:  "directory",
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
		Type:  "directory",
		Name:  name,
		DirID: parentID,
	}
	f.ByID[dir.ID] = dir
	f.addToChangesFeed(dir)
	return dir, nil
}

func (f *Fake) Refresh() error {
	return nil
}

func (f *Fake) Synchronized() error {
	f.SyncCount++
	return nil
}

// TODO add a CheckInvariants method

func (f *Fake) addToChangesFeed(doc *remote.Doc) {
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
