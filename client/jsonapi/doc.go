package jsonapi

import (
	"encoding/json"
	"io"

	"github.com/nono/cozy-desktop-experiments/state/remote"
)

// Doc describes a JSON-API document.
// See https://jsonapi.org/format/#document-structure
type Doc struct {
	Data Data `json:"data"`
}

// List describes a JSON-API list of documents.
type List struct {
	Data []Data `json:"data"`
}

// Data describes an item inside data.
type Data struct {
	ID   remote.ID `json:"id"`
	Meta struct {
		Rev remote.Rev `json:"rev"`
	} `json:"meta"`
	Attributes struct {
		Type  string    `json:"type"`
		Name  string    `json:"name"`
		DirID remote.ID `json:"dir_id"`
	} `json:"attributes"`
}

// ParseDoc tries to parse a JSON-API document from a reader, and then converts
// it to remote.Doc.
func ParseDoc(r io.Reader) (*remote.Doc, error) {
	var doc Doc
	if err := json.NewDecoder(r).Decode(&doc); err != nil {
		return nil, err
	}
	return &remote.Doc{
		ID:    doc.Data.ID,
		Rev:   doc.Data.Meta.Rev,
		Type:  doc.Data.Attributes.Type,
		Name:  doc.Data.Attributes.Name,
		DirID: doc.Data.Attributes.DirID,
	}, nil
}

// ParseDoc tries to parse a JSON-API list from a reader, and then converts
// it to a slice of remote.Doc.
func ParseList(r io.Reader) ([]*remote.Doc, error) {
	var list List
	if err := json.NewDecoder(r).Decode(&list); err != nil {
		return nil, err
	}
	docs := make([]*remote.Doc, 0, len(list.Data))
	for _, doc := range list.Data {
		docs = append(docs, &remote.Doc{
			ID:    doc.ID,
			Rev:   doc.Meta.Rev,
			Type:  doc.Attributes.Type,
			Name:  doc.Attributes.Name,
			DirID: doc.Attributes.DirID,
		})
	}
	return docs, nil
}
