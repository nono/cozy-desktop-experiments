package client

import (
	"errors"
	"fmt"
	"hash/crc32"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/nono/cozy-desktop-experiments/ng/state/remote"
)

// https://github.com/cozy/cozy-stack/blob/master/model/vfs/vfs.go#L24
const forbiddenFilenameChars = "/\x00\n\r"

// Fake can be used to simulate a cozy-stack client (and the stack its-self)
// for tests.
type Fake struct {
	Address   string
	SyncCount int
	ByID      map[remote.ID]*remote.Doc
	Feed      []Change

	// Those functions can be overloaded for some tests where we want to
	// control the values.
	GenerateID   func() remote.ID
	GenerateRev  func(id remote.ID, generation int) remote.Rev
	GenerateSeq  func(generation int) remote.Seq
	ConflictName func(id remote.ID, name string) string
}

// Change describes an entry in the changes feed of a fake stack/client.
type Change struct {
	Seq int
	*remote.Doc
	Skip bool
}

// NewFake creates a fake client that can be used for tests. It doesn't make
// any HTTP request, it just simulate them via an in-memory mock.
func NewFake(address string) remote.Client {

	fake := &Fake{
		Address:      address,
		SyncCount:    0,
		ByID:         map[remote.ID]*remote.Doc{},
		Feed:         []Change{},
		GenerateID:   newUUID,
		GenerateRev:  newRev,
		GenerateSeq:  newSeq,
		ConflictName: conflictName,
	}

	return fake
}

// AddInitialDocs will create the tree for a new instance (root, trash, etc.).
func (f *Fake) AddInitialDocs() {
	root := &remote.Doc{
		ID:   remote.RootID,
		Rev:  f.GenerateRev(remote.RootID, 1),
		Type: remote.Directory,
	}
	trash := &remote.Doc{
		ID:    remote.TrashID,
		Rev:   f.GenerateRev(remote.TrashID, 1),
		Type:  remote.Directory,
		Name:  remote.TrashName,
		DirID: root.ID,
	}
	// TODO add directories like Photos
	// TODO add some design docs to the changes feed
	for _, doc := range []*remote.Doc{root, trash} {
		f.ByID[doc.ID] = doc
		f.addToChangesFeed(doc)
	}
}

// Changes is required by the remote.Client interface.
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

// CreateDir is required by the remote.Client interface.
func (f *Fake) CreateDir(parentID remote.ID, name string) (*remote.Doc, error) {
	// TODO find a way to simulate latency
	if name == "" {
		return nil, errors.New("CreateDir: name is missing")
	}
	if strings.ContainsAny(name, forbiddenFilenameChars) {
		return nil, errors.New("CreateDir: name is invalid")
	}
	if _, ok := f.ByID[parentID]; !ok {
		return nil, errors.New("CreateDir: parent does not exist")
	}

	id := f.GenerateID()
	dir := &remote.Doc{
		ID:    id,
		Rev:   f.GenerateRev(id, 1),
		Type:  remote.Directory,
		Name:  name,
		DirID: parentID,
	}
	f.ByID[dir.ID] = dir
	f.addToChangesFeed(dir)
	return dir, nil
}

// Trash is required by the remote.Client interface.
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
	if f.isInTrash(was) {
		return nil, errors.New("Trash: already in the trash")
	}
	for _, other := range f.ByID {
		if other.DirID == remote.TrashID && other.Name == was.Name {
			was.Name = f.ConflictName(was.ID, was.Name)
		}
	}
	was.DirID = remote.TrashID
	was.Rev = f.GenerateRev(was.ID, extractGeneration(was.Rev)+1)
	f.addToChangesFeed(was)
	return was, nil
}

// Refresh is required by the remote.Client interface.
func (f *Fake) Refresh() error {
	return nil
}

// Synchronized is required by the remote.Client interface.
func (f *Fake) Synchronized() error {
	f.SyncCount++
	return nil
}

// CheckInvariants checks that we don't have inconsistencies in the fake
// client. It can be used as a way to detect some bugs in the Fake code.
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

// addToChangesFeed adds an entry for the given document in the changes feed.
// It masks previous entries for the same document, as CouchDB does.
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

func (f *Fake) isInTrash(doc *remote.Doc) bool {
	switch doc.DirID {
	case remote.RootID:
		return false
	case remote.TrashID:
		return true
	}
	parent, ok := f.ByID[doc.DirID]
	if !ok {
		panic(errors.New("parent not found"))
	}
	return f.isInTrash(parent)
}

// newUUID returns a compact UUID, similar to those used by CouchDB.
func newUUID() remote.ID {
	guid := uuid.Must(uuid.NewV4())
	id := fmt.Sprintf("%s", guid)
	id = strings.Replace(id, "-", "", -1)
	return remote.ID(id)
}

// newRev takes a generation number and returns a new revision for it.
func newRev(id remote.ID, generation int) remote.Rev {
	hashable := fmt.Sprintf("%d-%s", generation, id)
	hash := crc32.ChecksumIEEE([]byte(hashable))
	rev := fmt.Sprintf("%d-%0x", generation, hash)
	return remote.Rev(rev)
}

// newSeq takes a generation number and returns a new sequence for it.
func newSeq(generation int) remote.Seq {
	seq := fmt.Sprintf("%d-seq", generation) // TODO improve it
	return remote.Seq(seq)
}

// https://github.com/cozy/cozy-stack/blob/master/model/vfs/rand_suffix.go
var conflictRand = uint32(time.Now().UnixNano() + int64(os.Getpid()))

func conflictName(id remote.ID, name string) string {
	conflictRand = conflictRand*1664525 + 1013904223
	suffix := strconv.Itoa(int(1e9 + conflictRand%1e9))[1:]
	// https://github.com/cozy/cozy-stack/blob/master/model/vfs/vfs.go#L46
	return fmt.Sprintf("%s (__cozy__: %s)", name, suffix)
}

// extractGeneration returns the generation number of a revision. It is the
// first part of the revision, before the "-".
func extractGeneration(rev remote.Rev) int {
	parts := strings.Split(string(rev), "-")
	n, err := strconv.Atoi(parts[0])
	if err != nil {
		panic(fmt.Errorf("cannot extract generation from rev %s", rev))
	}
	return n
}
