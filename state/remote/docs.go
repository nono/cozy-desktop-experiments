package remote

import (
	"strings"

	"github.com/nono/cozy-desktop-experiments/state/types"
)

// Docs is keeping the information about the files on the Cozy.
type Docs struct {
	ByID            map[ID]*Doc
	ByDirID         map[ID]map[ID]*Doc // dirID -> map of children
	Seq             *Seq
	FetchingChanges bool
	Refreshing      bool
	RefreshedAt     types.Clock
}

// Doc describes a CouchDB document for the io.cozy.files doctype on the Cozy.
type Doc struct {
	ID    ID
	Rev   Rev
	Type  types.Type
	Name  string
	DirID ID
}

// ID is used for identifying the CouchDB documents.
type ID string

// IsDesignDoc returns true if the id is reserved for a design document.
func (id ID) IsDesignDoc() bool {
	return strings.HasPrefix(string(id), "_design")
}

// Rev is used by CouchDB to avoid conflicts.
type Rev string

// NewDocs creates a new docs.
func NewDocs() *Docs {
	return &Docs{
		ByID:    make(map[ID]*Doc),
		ByDirID: make(map[ID]map[ID]*Doc),
		Seq:     nil,
	}
}

// Upsert will add or update the given doc in the docs.
func (docs *Docs) Upsert(doc *Doc) {
	docs.detachParent(doc.ID)
	docs.ByID[doc.ID] = doc
	docs.attachParent(doc)
}

// MarkAsDeleted is used to mark a document as deleted.
func (docs *Docs) MarkAsDeleted(id ID) {
	docs.detachParent(id)
	delete(docs.ByID, id)
}

// Children returns a map of ID -> doc for the children of the given doc.
func (docs *Docs) Children(parent *Doc) map[ID]*Doc {
	return docs.ByDirID[parent.ID]
}

// attachParent updates the ByDirID field when a child is added to a
// directory.
func (docs *Docs) attachParent(child *Doc) {
	children := docs.ByDirID[child.DirID]
	if children == nil {
		children = make(map[ID]*Doc)
	}
	children[child.ID] = child
	docs.ByDirID[child.DirID] = children
}

// attachParent updates the ByDirID field when a child is removed a
// directory.
func (docs *Docs) detachParent(childID ID) {
	if was, ok := docs.ByID[childID]; ok {
		delete(docs.ByDirID[was.DirID], was.ID)
	}
}
