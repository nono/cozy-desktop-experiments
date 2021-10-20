package local

import "io/fs"

// FS is an interface for making calls to a local file system.
type FS interface {
	fs.StatFS
	fs.ReadDirFS
	Mkdir(path string) error
	RemoveAll(path string) error
}
