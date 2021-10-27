package remote

import (
	"strings"

	"github.com/nono/cozy-desktop-experiments/state/types"
)

// Docs is keeping the information about the files on the Cozy.
type Docs struct {
	ByID            map[ID]*Doc
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
		ByID: make(map[ID]*Doc),
		Seq:  nil,
	}
}

// Upsert will add or update the given doc in the docs.
func (docs *Docs) Upsert(doc *Doc) {
	docs.ByID[doc.ID] = doc
}

// MarkAsDeleted is used to mark a document as deleted.
func (docs *Docs) MarkAsDeleted(id ID) {
	delete(docs.ByID, id)
}
