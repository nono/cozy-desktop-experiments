package jsonapi

import (
	"encoding/json"
	"io"

	"github.com/nono/cozy-desktop-experiments/ng/state/remote"
)

type Doc struct {
	Data struct {
		ID   remote.ID `json:"id"`
		Meta struct {
			Rev remote.Rev `json:"rev"`
		} `json:"meta"`
		Attributes struct {
			Type  string    `json:"type"`
			Name  string    `json:"name"`
			DirID remote.ID `json:"dir_id"`
		} `json:"attributes"`
	} `json:"data"`
}

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