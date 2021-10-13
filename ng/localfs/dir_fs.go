package localfs

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/nono/cozy-desktop-experiments/ng/state/local"
)

const Separator = "/"

func NewDirFS(dir string) local.FS {
	return dirFS(dir)
}

type dirFS string

func (dir dirFS) Open(path string) (fs.File, error) {
	if !fs.ValidPath(path) {
		return nil, &os.PathError{Op: "open", Path: path, Err: os.ErrInvalid}
	}
	f, err := os.Open(string(dir) + Separator + path)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (dir dirFS) Stat(path string) (fs.FileInfo, error) {
	if !fs.ValidPath(path) {
		return nil, &os.PathError{Op: "stat", Path: path, Err: os.ErrInvalid}
	}
	info, err := os.Stat(string(dir) + Separator + path)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func (dir dirFS) ReadDir(path string) ([]fs.DirEntry, error) {
	if !fs.ValidPath(path) {
		return nil, &os.PathError{Op: "readdir", Path: path, Err: os.ErrInvalid}
	}
	entries, err := os.ReadDir(string(dir) + Separator + path)
	if err != nil {
		return nil, err
	}
	return entries, nil
}

func (dir dirFS) Mkdir(path string) error {
	if !fs.ValidPath(path) {
		return &os.PathError{Op: "mkdir", Path: path, Err: os.ErrInvalid}
	}
	return os.Mkdir(string(dir)+Separator+path, 0755)
}

func (dir dirFS) ToMemFS() (*memFS, error) {
	mem := NewMemFS().(*memFS)
	err := dir.addToMemFS(mem, ".")
	return mem, err
}

func (dir dirFS) addToMemFS(mem *memFS, path string) error {
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

var _ fs.FS = dirFS(".")
