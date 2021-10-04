package remote

type Client interface {
	Changes(seq *Seq) (*ChangesResponse, error)
	Refresh() error
	Synchronized() error
}

type ChangesResponse struct {
	Docs    []*Doc
	Seq     Seq
	Pending int
}

type Seq string
