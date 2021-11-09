package client

import (
	"errors"

	"github.com/nono/cozy-desktop-experiments/state/remote"
)

// NewReadOnly returns a mock of a cozy-stack client that can be used in tests,
// and it panics if any write operation is called.
func NewReadOnly(address string) *ReadOnly {
	fake := NewFake(address)
	return &ReadOnly{Fake: fake}
}

type ReadOnly struct {
	Fake *Fake
}

// Changes is required by the remote.Client interface.
func (ro *ReadOnly) Changes(seq *remote.Seq, limit int, skipTrashed bool) (*remote.ChangesResponse, error) {
	return ro.Fake.Changes(seq, limit, skipTrashed)
}

// CreateDir is required by the remote.Client interface.
func (ro *ReadOnly) CreateDir(parentID remote.ID, name string) (*remote.Doc, error) {
	panic(errors.New("CreateDir has been called on ReadOnly client"))
}

// Trash is required by the remote.Client interface.
func (ro *ReadOnly) Trash(doc *remote.Doc) (*remote.Doc, error) {
	panic(errors.New("Trash has been called on ReadOnly client"))
}

// EmptyTrash is required by the remote.Client interface.
func (ro *ReadOnly) EmptyTrash() error {
	panic(errors.New("EmptyTrash has been called on ReadOnly client"))
}

// Refresh is required by the remote.Client interface.
func (ro *ReadOnly) Refresh() error {
	return ro.Fake.Refresh()
}

// Synchronized is required by the remote.Client interface.
func (ro *ReadOnly) Synchronized() error {
	return ro.Fake.Synchronized()
}

var _ remote.Client = &ReadOnly{}
