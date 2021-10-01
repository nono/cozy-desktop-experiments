package localfs

import (
	"io/fs"
	"os"

	"github.com/nono/cozy-desktop-experiments/ng/state/local"
)

func DirFS(dir string) local.FS {
	return dirFS(dir)
}

type dirFS string

func (dir dirFS) Open(name string) (fs.File, error) {
	if !fs.ValidPath(name) {
		return nil, &os.PathError{Op: "open", Path: name, Err: os.ErrInvalid}
	}
	f, err := os.Open(string(dir) + "/" + name)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (dir dirFS) Stat(name string) (fs.FileInfo, error) {
	if !fs.ValidPath(name) {
		return nil, &os.PathError{Op: "stat", Path: name, Err: os.ErrInvalid}
	}
	info, err := os.Stat(string(dir) + "/" + name)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func (dir dirFS) ReadDir(name string) ([]fs.DirEntry, error) {
	if !fs.ValidPath(name) {
		return nil, &os.PathError{Op: "readdir", Path: name, Err: os.ErrInvalid}
	}
	entries, err := os.ReadDir(string(dir) + "/" + name)
	if err != nil {
		return nil, err
	}
	return entries, nil
}

var _ fs.FS = dirFS(".")
