package localfs

import (
	"io/fs"
	"os"

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

var _ fs.FS = dirFS(".")
