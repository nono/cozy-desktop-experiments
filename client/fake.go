package client

import (
	"errors"
	"fmt"
	"hash/crc32"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/nono/cozy-desktop-experiments/state/remote"
	"github.com/nono/cozy-desktop-experiments/state/types"
)

// https://github.com/cozy/cozy-stack/blob/master/model/vfs/vfs.go#L24
const forbiddenFilenameChars = "/\x00\n\r"

// Fake can be used to simulate a cozy-stack client (and the stack its-self)
// for tests.
//
// TODO find a way to simulate latency
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
	*remote.ChangedDoc
	Skip bool
}

// NewFake creates a fake client that can be used for tests. It doesn't make
// any HTTP request, it just simulate them via an in-memory mock.
func NewFake(address string) *Fake {
	return &Fake{
		Address:      address,
		SyncCount:    0,
		ByID:         map[remote.ID]*remote.Doc{},
		Feed:         []Change{},
		GenerateID:   newUUID,
		GenerateRev:  newRev,
		GenerateSeq:  newSeq,
		ConflictName: conflictName,
	}
}

// AddInitialDocs will create the tree for a new instance (root, trash, etc.).
func (f *Fake) AddInitialDocs(changed ...*remote.ChangedDoc) {
	var docs []*remote.Doc
	if len(changed) == 0 {
		root := &remote.Doc{
			ID:   remote.RootID,
			Rev:  f.GenerateRev(remote.RootID, 1),
			Type: types.DirType,
		}
		trash := &remote.Doc{
			ID:    remote.TrashID,
			Rev:   f.GenerateRev(remote.TrashID, 1),
			Type:  types.DirType,
			Name:  remote.TrashName,
			DirID: root.ID,
		}
		docs = []*remote.Doc{root, trash}
	} else {
		for _, doc := range changed {
			docs = append(docs, doc.Doc)
		}
	}
	for _, doc := range docs {
		f.ByID[doc.ID] = doc
		f.addToChangesFeed(doc)
	}
}

// MatchSequence will create dumb entries in the changes feed until it reaches
// the given sequence generation number. It is used to compensate for the lack
// of design docs.
func (f *Fake) MatchSequence(seq remote.Seq) {
	gen := seq.ExtractGeneration()
	for {
		last := f.Feed[len(f.Feed)-1]
		if last.Seq >= gen {
			return
		}
		f.Feed = append(f.Feed, Change{
			Seq:        len(f.Feed) + 1,
			ChangedDoc: &remote.ChangedDoc{Doc: &remote.Doc{}},
			Skip:       true,
		})
	}
}

// Changes is required by the remote.Client interface.
func (f *Fake) Changes(seq *remote.Seq, limit int, skipTrashed bool) (*remote.ChangesResponse, error) {
	since := 0
	if seq != nil {
		since = seq.ExtractGeneration()
	}
	lastSeq := since
	docs := []*remote.ChangedDoc{}
	for _, c := range f.Feed {
		if c.Seq <= since {
			continue
		}
		lastSeq = c.Seq
		if c.Skip {
			continue
		}
		if skipTrashed && (c.Deleted || f.isInTrash(c.ChangedDoc.Doc)) {
			continue
		}
		docs = append(docs, c.ChangedDoc)
	}
	pending := 0
	if len(docs) > limit {
		pending = len(docs) - limit
		docs = docs[:limit]
	}
	return &remote.ChangesResponse{Docs: docs, Seq: f.GenerateSeq(lastSeq), Pending: pending}, nil
}

// CreateDir is required by the remote.Client interface.
func (f *Fake) CreateDir(parentID remote.ID, name string) (*remote.Doc, error) {
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
		Type:  types.DirType,
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
	f.updateTreeRev(was)
	was.DirID = remote.TrashID
	was.Rev = f.GenerateRev(was.ID, extractGeneration(was.Rev)+1)
	f.addToChangesFeed(was)
	return was, nil
}

// EmptyTrash is required by remote.Client interface.
func (f *Fake) EmptyTrash() error {
	for id, doc := range f.ByID {
		if !f.isInTrash(doc) {
			continue
		}
		f.addDeletedToChangesFeed(id)
		delete(f.ByID, id)
	}
	return nil
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

// DocsByID returns a map of id -> doc (for testing purpose).
func (f *Fake) DocsByID() map[remote.ID]*remote.Doc {
	return f.ByID
}

// CheckInvariants checks that we don't have inconsistencies in the fake
// client. It can be used as a way to detect some bugs in the Fake code.
func (f *Fake) CheckInvariants() error {
	root, ok := f.ByID[remote.RootID]
	if !ok {
		return errors.New("root is missing")
	}
	if root.Type != types.DirType {
		return errors.New("root is not a directory")
	}
	trash, ok := f.ByID[remote.TrashID]
	if !ok {
		return errors.New("trash is missing")
	}
	if trash.Type != types.DirType {
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
	if parent.Type != types.DirType {
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
		Seq:        len(f.Feed) + 1,
		ChangedDoc: &remote.ChangedDoc{Doc: doc, Deleted: false},
	}
	f.Feed = append(f.Feed, change)
}

// addDeletedToChangesFeed adds an entry to the changes feed for a deleted
// document. It masks previous entries for the same document, as CouchDB does.
func (f *Fake) addDeletedToChangesFeed(id remote.ID) {
	for i, change := range f.Feed {
		if change.Doc.ID == id {
			f.Feed[i].Skip = true
		}
	}
	doc := &remote.Doc{ID: id}
	change := Change{
		Seq:        len(f.Feed) + 1,
		ChangedDoc: &remote.ChangedDoc{Doc: doc, Deleted: true},
	}
	f.Feed = append(f.Feed, change)
}

// updateTreeRev is used to update the revisions of all the descendants of the
// given directory.
func (f *Fake) updateTreeRev(dir *remote.Doc) {
	var children []*remote.Doc
	for id, doc := range f.ByID {
		if doc.DirID == dir.ID {
			children = append(children, f.ByID[id])
		}
	}
	sort.Slice(children, func(i, j int) bool {
		return children[i].Name < children[j].Name
	})
	for _, child := range children {
		dir.Rev = f.GenerateRev(dir.ID, extractGeneration(dir.Rev)+1)
		f.addToChangesFeed(child)
		f.updateTreeRev(child)
	}
}

// isInTrash returns true if the doc is a descendant of the trash.
func (f *Fake) isInTrash(doc *remote.Doc) bool {
	switch doc.ID {
	case remote.RootID:
		return false
	case remote.TrashID:
		return true
	}
	parent, ok := f.ByID[doc.DirID]
	if !ok {
		panic(fmt.Errorf("parent not found for %#v", doc))
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

var _ remote.Client = &Fake{}
