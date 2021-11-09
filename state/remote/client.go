package remote

import (
	"strconv"
	"strings"
)

// Client is an interface for a client that can make requests to the
// cozy-stack to manipulate files.
type Client interface {
	// Changes is used to request the changes feed since a sequence number.
	Changes(seq *Seq, limit int, skipTrashed bool) (*ChangesResponse, error)

	// CreateDir will create a directory on the Cozy.
	CreateDir(parentID ID, name string) (*Doc, error)

	// Trash will put a file or directory inside the trash on the Cozy.
	Trash(doc *Doc) (*Doc, error)

	// EmptyTrash will delete the files and directories in the trash. Useful
	// for tests.
	EmptyTrash() error

	// Refresh can be used to refresh the OAuth access token.
	Refresh() error

	// Synchronized can be called to inform the Cozy that the client is now
	// synchronized. The data of last synchronization is shown in
	// cozy-settings.
	Synchronized() error
}

// ChangesResponse describes the successful response to a call to the changes
// feed.
type ChangesResponse struct {
	Docs    []*ChangedDoc
	Seq     Seq
	Pending int
}

// ChangedDoc is a result item from the changes feed.
type ChangedDoc struct {
	*Doc
	Deleted bool
}

// Seq is the short for sequence. It is a way to keep a position on the changes
// feed for the next calls.
type Seq string

// ExtractGeneration returns the first part of a sequence. The generation is
// the number before the "-".
func (s Seq) ExtractGeneration() int {
	parts := strings.SplitN(string(s), "-", 2)
	n, err := strconv.Atoi(parts[0])
	if err != nil {
		panic(err)
	}
	return n
}
