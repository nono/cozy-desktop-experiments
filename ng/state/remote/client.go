package remote

import (
	"strconv"
	"strings"
)

type Client interface {
	Changes(seq *Seq) (*ChangesResponse, error)
	CreateDir(parentID ID, name string) (*Doc, error)
	Trash(doc *Doc) (*Doc, error)
	Refresh() error
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
