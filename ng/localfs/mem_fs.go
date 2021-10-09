package localfs

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/nono/cozy-desktop-experiments/ng/state/local"
)

func MemFS() local.FS {
	baseDir := &memDir{
		info: memFileInfo{
			name:    ".",
			size:    4096,
			mode:    fs.ModeDir | 0755,
			modTime: time.Now(),
			sys:     &syscall.Stat_t{Ino: 1},
		},
		path: ".",
	}
	return memFS{
		ByPath: map[string]*memDir{".": baseDir},
	}
}

type memFS struct {
	ByPath map[string]*memDir
}

type memFile struct {
	info memFileInfo
}

func (f memFile) Stat() (fs.FileInfo, error) { return f.info, nil }
func (f memFile) Read(b []byte) (int, error) { return 0, errors.New("Not implemeted") }
func (f memFile) Close() error               { return nil }

type memDir struct {
	info     memFileInfo
	path     string
	children []fs.DirEntry
}

func (d memDir) Name() string               { return d.info.name }
func (d memDir) IsDir() bool                { return true }
func (d memDir) Type() fs.FileMode          { return d.info.mode }
func (d memDir) Info() (fs.FileInfo, error) { return d.info, nil }
func (d memDir) Stat() (fs.FileInfo, error) { return d.info, nil }
func (d memDir) Read(b []byte) (int, error) {
	return 0, &fs.PathError{Op: "read", Path: d.path, Err: fs.ErrInvalid}
}
func (d memDir) Close() error { return nil }
func (d memDir) ReadDir(count int) ([]fs.DirEntry, error) {
	if count <= 0 || len(d.children) < count {
		return d.children, nil
	}
	return nil, io.EOF
}

type memFileInfo struct {
	name    string
	size    int64
	mode    fs.FileMode
	modTime time.Time
	sys     *syscall.Stat_t
}

func (info memFileInfo) Name() string       { return info.name }
func (info memFileInfo) Size() int64        { return info.size }
func (info memFileInfo) Mode() fs.FileMode  { return info.mode }
func (info memFileInfo) ModTime() time.Time { return info.modTime }
func (info memFileInfo) IsDir() bool        { return info.mode.IsDir() }
func (info memFileInfo) Sys() interface{}   { return info.sys }

func (mem memFS) Open(name string) (fs.File, error) {
	dir, ok := mem.ByPath[name]
	if !ok {
		return nil, &os.PathError{Op: "open", Path: name, Err: os.ErrInvalid}
	}
	return dir, nil
}

func (mem memFS) Stat(name string) (fs.FileInfo, error) {
	if !fs.ValidPath(name) {
		return nil, &os.PathError{Op: "stat", Path: name, Err: os.ErrInvalid}
	}
	f, err := mem.Open(name)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	return f.Stat()
}

func (mem memFS) ReadDir(name string) ([]fs.DirEntry, error) {
	if !fs.ValidPath(name) {
		return nil, &os.PathError{Op: "readdir", Path: name, Err: os.ErrInvalid}
	}
	dir, err := mem.Open(name)
	if err != nil {
		return nil, err
	}
	if d, ok := dir.(memDir); ok {
		return d.ReadDir(-1)
	}
	return nil, &os.PathError{Op: "readdir", Path: name, Err: os.ErrInvalid}
}

func (mem memFS) Mkdir(name string) error {
	if !fs.ValidPath(name) || name == "." {
		return &os.PathError{Op: "mkdir", Path: name, Err: os.ErrInvalid}
	}
	if _, ok := mem.ByPath[name]; ok {
		return &os.PathError{Op: "mkdir", Path: name, Err: os.ErrInvalid}
	}
	parentname, dirname := filepath.Split(name)
	if parentname == "" {
		parentname = "."
	} else {
		parentname = strings.TrimSuffix(parentname, Separator)
	}
	parent, ok := mem.ByPath[parentname]
	if !ok {
		return &os.PathError{Op: "mkdir", Path: name, Err: os.ErrInvalid}
	}
	dir := &memDir{
		info: memFileInfo{
			name:    dirname,
			size:    4096,
			mode:    fs.ModeDir | 0755,
			modTime: time.Now(),
			sys:     &syscall.Stat_t{Ino: mem.NextIno()},
		},
		path: name,
	}
	mem.ByPath[name] = dir
	parent.children = append(parent.children, dir)
	return nil
}

func (mem memFS) NextIno() uint64 {
	return uint64(len(mem.ByPath) + 1) // TODO
}

var _ fs.FS = memFS{}
var _ fs.File = memFile{}
var _ fs.ReadDirFile = memDir{}
