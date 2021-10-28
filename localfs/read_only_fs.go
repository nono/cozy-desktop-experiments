package localfs

import (
	"errors"
	"io/fs"

	"github.com/nono/cozy-desktop-experiments/state/local"
)

// NewMemFS returns an in-memory mock of local.FS for tests that panic if any
// write operation is called on it.
func NewReadOnlyFS() local.FS {
	mem := NewMemFS().(*MemFS)
	return &ReadOnlyFS{MemFS: mem}
}

type ReadOnlyFS struct {
	MemFS *MemFS
}

// Open is required by the local.FS interface.
func (ro *ReadOnlyFS) Open(path string) (fs.File, error) {
	return ro.MemFS.Open(path)
}

// Stat is required by the local.FS interface.
func (ro *ReadOnlyFS) Stat(path string) (fs.FileInfo, error) {
	return ro.MemFS.Stat(path)
}

// ReadDir is required by the local.FS interface.
func (ro *ReadOnlyFS) ReadDir(path string) ([]fs.DirEntry, error) {
	return ro.MemFS.ReadDir(path)
}

// Mkdir is required by the local.FS interface.
func (ro *ReadOnlyFS) Mkdir(path string) error {
	panic(errors.New("Mkdir has been called for ReadOnlyFS"))
}

// RemoveAll is required by the local.FS interface.
func (ro *ReadOnlyFS) RemoveAll(path string) error {
	panic(errors.New("RemoveAll has been called for ReadOnlyFS"))
}

var _ fs.FS = &ReadOnlyFS{}
