package sync

type Platform struct {
	Events chan Event
	Local  LocalFS
}

func NewPlatform(local LocalFS) *Platform {
	return &Platform{
		Events: make(chan Event),
		Local:  local,
	}
}
