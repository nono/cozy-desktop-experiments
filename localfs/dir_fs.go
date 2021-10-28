// Package localfs provides an implementation of local.FS that works on a given
// directory of the local file system. It also provides an in-memory mocks for
// tests.
package localfs

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/nono/cozy-desktop-experiments/state/local"
)

// NewDirFS returns a local.FS that makes changes to the given Cozy directory
// on the local disk.
func NewDirFS(dir string) local.FS {
	return DirFS(dir)
}

type DirFS string

// Open is required by the local.FS interface.
func (dir DirFS) Open(path string) (fs.File, error) {
	if !fs.ValidPath(path) {
		return nil, &os.PathError{Op: "open", Path: path, Err: os.ErrInvalid}
	}
	abspath := filepath.Join(string(dir), path)
	f, err := os.Open(abspath)
	if err != nil {
		return nil, err
	}
	return f, nil
}

// Stat is required by the local.FS interface.
func (dir DirFS) Stat(path string) (fs.FileInfo, error) {
	if !fs.ValidPath(path) {
		return nil, &os.PathError{Op: "stat", Path: path, Err: os.ErrInvalid}
	}
	abspath := filepath.Join(string(dir), path)
	info, err := os.Stat(abspath)
	if err != nil {
		return nil, err
	}
	return info, nil
}

// ReadDir is required by the local.FS interface.
func (dir DirFS) ReadDir(path string) ([]fs.DirEntry, error) {
	if !fs.ValidPath(path) {
		return nil, &os.PathError{Op: "readDir", Path: path, Err: os.ErrInvalid}
	}
	abspath := filepath.Join(string(dir), path)
	entries, err := os.ReadDir(abspath)
	if err != nil {
		return nil, err
	}
	return entries, nil
}

// Mkdir is required by the local.FS interface.
func (dir DirFS) Mkdir(path string) error {
	if !fs.ValidPath(path) {
		return &os.PathError{Op: "mkdir", Path: path, Err: os.ErrInvalid}
	}
	abspath := filepath.Join(string(dir), path)
	return os.Mkdir(abspath, 0750)
}

// RemoveAll is required by the local.FS interface.
func (dir DirFS) RemoveAll(path string) error {
	if !fs.ValidPath(path) || path == "." {
		return &os.PathError{Op: "mkdir", Path: path, Err: os.ErrInvalid}
	}
	abspath := filepath.Join(string(dir), path)
	return os.RemoveAll(abspath)
}

// ToMemFS will create a MemFS with the same files and directories. It can be
// useful for testing purpose.
func (dir DirFS) ToMemFS() (*MemFS, error) {
	mem := NewMemFS().(*MemFS)
	err := dir.addToMemFS(mem, ".")
	return mem, err
}

func (dir DirFS) addToMemFS(mem *MemFS, path string) error {
	entries, err := dir.ReadDir(path)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			return errors.New("Unexpected entry")
		}
		entryPath := filepath.Join(path, entry.Name())
		if err := mem.Mkdir(entryPath); err != nil {
			return err
		}
		if err := dir.addToMemFS(mem, entryPath); err != nil {
			return err
		}
	}
	return nil
}

var _ fs.FS = DirFS(".")
