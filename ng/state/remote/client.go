package remote

import (
	"strconv"
	"strings"
)

type Client interface {
	// Changes is used to request the changes feed since a sequence number.
	// TODO add a limit option
	// TODO add an option to skip deleted docs
	Changes(seq *Seq) (*ChangesResponse, error)

	// CreateDir will create a directory on the Cozy.
	CreateDir(parentID ID, name string) (*Doc, error)

	// Trash will put a file or directory inside the trash on the Cozy.
	Trash(doc *Doc) (*Doc, error)

	// Refresh can be used to refresh the OAuth access token.
	Refresh() error

	// Synchronized can be called to inform the Cozy that the client is now
	// synchronized. The data of last synchronization is shown in
	// cozy-settings.
	Synchronized() error
}

type ChangesResponse struct {
	Docs    []*Doc
	Seq     Seq
	Pending int
}

type Seq string

func (s Seq) ExtractGeneration() int {
	parts := strings.SplitN(string(s), "-", 2)
	n, err := strconv.Atoi(parts[0])
	if err != nil {
		panic(err)
	}
	return n
}
