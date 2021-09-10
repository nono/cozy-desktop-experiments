package sync

import "io/fs"

type Platform struct {
	Events chan Event
	Local  fs.FS
}

func NewPlatform(local fs.FS) *Platform {
	return &Platform{
		Events: make(chan Event),
		Local:  local,
	}
}
