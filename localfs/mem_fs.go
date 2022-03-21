package localfs

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/nono/cozy-desktop-experiments/state/local"
)

// NewMemFS returns an in-memory mock of local.FS for tests.
// TODO use https://github.com/hack-pad/hackpadfs ?
func NewMemFS() *MemFS {
	baseDir := &memDir{
		info: &memFileInfo{
			name:    ".",
			size:    4096,
			mode:    fs.ModeDir | 0750,
			modTime: time.Now(),
			sys:     &syscall.Stat_t{Ino: 1},
		},
		path: ".",
	}
	return &MemFS{
		ByPath: map[string]*memDir{".": baseDir},
	}
}

type MemFS struct {
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

// Open is required by the local.FS interface.
func (mem *MemFS) Open(path string) (fs.File, error) {
	path = strings.TrimSuffix(path, "./")
	dir, ok := mem.ByPath[path]
	if !ok {
		return nil, &os.PathError{Op: "open", Path: path, Err: os.ErrInvalid}
	}
	return &memDirHandler{d: dir, pos: 0}, nil
}

// Stat is required by the local.FS interface.
func (mem *MemFS) Stat(path string) (fs.FileInfo, error) {
	if !validPath(path) {
		return nil, &os.PathError{Op: "stat", Path: path, Err: os.ErrInvalid}
	}
	f, err := mem.Open(path)
	defer func() { _ = f.Close() }()
	if err != nil {
		return nil, err
	}
	return f.Stat()
}

// ReadDir is required by the local.FS interface.
func (mem *MemFS) ReadDir(path string) ([]fs.DirEntry, error) {
	if !validPath(path) {
		return nil, &os.PathError{Op: "readDir", Path: path, Err: os.ErrInvalid}
	}
	handler, err := mem.Open(path)
	if err != nil {
		return nil, err
	}
	if dh, ok := handler.(*memDirHandler); ok {
		return dh.ReadDir(-1)
	}
	return nil, &os.PathError{Op: "readDir", Path: path, Err: os.ErrInvalid}
}

// Mkdir is required by the local.FS interface.
func (mem *MemFS) Mkdir(path string) error {
	if !validPath(path) || path == "." || endWithSlash(path) {
		return &os.PathError{Op: "mkdir", Path: path, Err: os.ErrInvalid}
	}
	path = filepath.Clean(path)
	if _, ok := mem.ByPath[path]; ok {
		return &os.PathError{Op: "mkdir", Path: path, Err: os.ErrInvalid}
	}
	name, parent, ok := mem.nameAndParent(path)
	if !ok {
		return &os.PathError{Op: "mkdir", Path: path, Err: os.ErrInvalid}
	}
	dir := &memDir{
		info: &memFileInfo{
			name:    name,
			size:    4096,
			mode:    fs.ModeDir | 0750,
			modTime: time.Now(),
			sys:     &syscall.Stat_t{Ino: mem.NextIno()},
		},
		path: path,
	}
	mem.ByPath[path] = dir
	parent.children = append(parent.children, dir)
	return nil
}

// RemoveAll is required by the local.FS interface.
func (mem *MemFS) RemoveAll(path string) error {
	if !validPath(path) || path == "." {
		return &os.PathError{Op: "removeAll", Path: path, Err: os.ErrInvalid}
	}
	path = filepath.Clean(path)
	dir, ok := mem.ByPath[path]
	if !ok {
		return nil
	}
	_, parent, ok := mem.nameAndParent(path)
	if !ok {
		return &os.PathError{Op: "removeAll", Path: path, Err: os.ErrInvalid}
	}

	mem.removeDescendants(dir)

	// https://github.com/golang/go/wiki/SliceTricks#filtering-without-allocating
	children := parent.children[:0]
	for _, child := range parent.children {
		if child != dir {
			children = append(children, child)
		}
	}
	for i := len(children); i < len(parent.children); i++ {
		parent.children[i] = nil
	}
	parent.children = children

	delete(mem.ByPath, path)
	return nil
}

// removeDescendants remove every thing inside a directory (files and
// sub-directories).
func (mem *MemFS) removeDescendants(dir *memDir) {
	for _, child := range dir.children {
		if child.IsDir() {
			mem.removeDescendants(child.(*memDir))
		}
		path := filepath.Join(dir.path, child.Name())
		delete(mem.ByPath, path)
	}
	dir.children = nil
}

// nameAndParent takes a path and return the name, and the parent directory.
func (mem *MemFS) nameAndParent(path string) (string, *memDir, bool) {
	parentPath, name := filepath.Split(path)
	if parentPath == "" || parentPath == "/" {
		parentPath = "."
	} else if endWithSlash(parentPath) {
		parentPath = parentPath[:len(parentPath)-1]
	}
	parent, ok := mem.ByPath[parentPath]
	return name, parent, ok
}

// CheckInvariants check some properties of the mock. It can be used to detect
// some bugs in the MemFS implementation.
func (mem *MemFS) CheckInvariants() error {
	if _, ok := mem.ByPath["."]; !ok {
		return errors.New("root is missing")
	}
	for _, dir := range mem.ByPath {
		if _, ok := mem.ByPath[filepath.Dir(dir.path)]; !ok {
			return fmt.Errorf("%#v has no parent", dir)
		}
		if dir.IsDir() != dir.Type().IsDir() {
			return fmt.Errorf("%#v is both a file and a directory", dir)
		}
		for _, child := range dir.children {
			child := child.(*memDir)
			if child.path != filepath.Join(dir.path, child.info.name) {
				return fmt.Errorf("%#v path is incorrect", child)
			}
		}
	}
	return nil
}

// NextIno returns the next free inode number that can be used.
func (mem *MemFS) NextIno() uint64 {
	return uint64(len(mem.ByPath) + 1) // TODO
}

// validPath checks that the path is valid, not something like foo/../bar or
// with a null character.
func validPath(path string) bool {
	return fs.ValidPath(path) && !strings.Contains(path, "\x00")
}

// endWithSlash returns true if the last character of the path is a slash.
func endWithSlash(path string) bool {
	return path[len(path)-1] == filepath.Separator
}

var _ local.FS = &MemFS{}
var _ fs.FS = &MemFS{}
var _ fs.DirEntry = &memFile{}
var _ fs.DirEntry = &memDir{}
var _ fs.File = &memFileHandler{}
var _ fs.ReadDirFile = &memDirHandler{}
var _ fs.FileInfo = &memFileInfo{}
