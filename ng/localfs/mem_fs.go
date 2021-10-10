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
		info: &memFileInfo{
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
	info *memFileInfo
}

func (f *memFile) Name() string               { return f.info.name }
func (f *memFile) IsDir() bool                { return false }
func (f *memFile) Type() fs.FileMode          { return f.info.mode }
func (f *memFile) Info() (fs.FileInfo, error) { return f.info, nil }

type memFileHandler struct {
	f *memFile
	// TODO closed bool
	// TODO pos int
}

func (fh *memFileHandler) Stat() (fs.FileInfo, error) { return fh.f.Info() }
func (fh *memFileHandler) Read(b []byte) (int, error) { return 0, errors.New("Not implemeted") }
func (fh *memFileHandler) Close() error               { return nil }

type memDir struct {
	info     *memFileInfo
	path     string
	children []fs.DirEntry
}

func (d *memDir) Name() string               { return d.info.name }
func (d *memDir) IsDir() bool                { return true }
func (d *memDir) Type() fs.FileMode          { return d.info.mode.Type() }
func (d *memDir) Info() (fs.FileInfo, error) { return d.info, nil }

type memDirHandler struct {
	d   *memDir
	pos int
}

func (dh *memDirHandler) Stat() (fs.FileInfo, error) { return dh.d.Info() }
func (dh *memDirHandler) Read(b []byte) (int, error) {
	return 0, &fs.PathError{Op: "read", Path: dh.d.path, Err: fs.ErrInvalid}
}
func (dh *memDirHandler) Close() error { return nil }
func (dh *memDirHandler) ReadDir(count int) ([]fs.DirEntry, error) {
	if dh.pos >= len(dh.d.children) {
		if count <= 0 {
			return []fs.DirEntry{}, nil
		}
		return nil, io.EOF
	}

	from := dh.pos
	to := dh.pos + count
	if count <= 0 || to > len(dh.d.children) {
		to = len(dh.d.children)
	}
	dh.pos += to - from
	return dh.d.children[from:to], nil
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
	name = strings.TrimSuffix(name, "./")
	dir, ok := mem.ByPath[name]
	if !ok {
		return nil, &os.PathError{Op: "open", Path: name, Err: os.ErrInvalid}
	}
	return &memDirHandler{d: dir, pos: 0}, nil
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
	handler, err := mem.Open(name)
	if err != nil {
		return nil, err
	}
	if dh, ok := handler.(*memDirHandler); ok {
		return dh.ReadDir(-1)
	}
	return nil, &os.PathError{Op: "readdir", Path: name, Err: os.ErrInvalid}
}

func (mem memFS) Mkdir(name string) error {
	name = filepath.Clean(name)
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
		info: &memFileInfo{
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

// TODO add a CheckInvariants method

var _ fs.FS = memFS{}
var _ fs.DirEntry = &memFile{}
var _ fs.DirEntry = &memDir{}
var _ fs.File = &memFileHandler{}
var _ fs.ReadDirFile = &memDirHandler{}
var _ fs.FileInfo = &memFileInfo{}
