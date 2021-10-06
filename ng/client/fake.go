package client

import (
	"errors"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/nono/cozy-desktop-experiments/ng/state/remote"
)

type Fake struct {
	Address   string
	SyncCount int
	ByID      map[remote.ID]*remote.Doc
	Feed      []Change
}

type Change struct {
	Seq int
	*remote.Doc
}

func NewFake(address string) remote.Client {
	root := &remote.Doc{
		ID:   remote.RootID,
		Rev:  newRev(1),
		Type: "directory",
	}
	trash := &remote.Doc{
		ID:    remote.TrashID,
		Rev:   newRev(1),
		Type:  "directory",
		Name:  remote.TrashName,
		DirID: root.ID,
	}
	byID := map[remote.ID]*remote.Doc{
		root.ID:  root,
		trash.ID: trash,
	}
	fake := &Fake{
		Address:   address,
		SyncCount: 0,
		ByID:      byID,
		Feed:      []Change{},
	}
	for _, id := range []remote.ID{root.ID, trash.ID} {
		change := Change{
			Seq: len(fake.Feed),
			Doc: fake.ByID[id],
		}
		fake.Feed = append(fake.Feed, change)
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
	lastSeq := "0-initial"
	docs := []*remote.Doc{}
	for _, c := range f.Feed {
		if c.Seq < since {
			continue
		}
		docs = append(docs, c.Doc)
		lastSeq = fmt.Sprintf("%d-seq", c.Seq) // TODO improve it
	}
	return &remote.ChangesResponse{Docs: docs, Seq: remote.Seq(lastSeq), Pending: 0}, nil
}

func (f *Fake) CreateDir(parentID remote.ID, name string) (*remote.Doc, error) {
	return nil, errors.New("Not yet implemented")
}

func (f *Fake) Refresh() error {
	return nil
}

func (f *Fake) Synchronized() error {
	f.SyncCount++
	return nil
}

func newUUID() string {
	id := uuid.Must(uuid.NewV4())
	return fmt.Sprintf("%s", id)
}

func newRev(generation int) remote.Rev {
	rev := fmt.Sprintf("%d-rev", generation) // TODO improve it
	return remote.Rev(rev)
}
