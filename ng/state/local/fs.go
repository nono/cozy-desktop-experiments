package local

import "io/fs"

type FS interface {
	fs.StatFS
	fs.ReadDirFS
	Mkdir(name string) error
}
