package sync

import "os"

type Operation interface {
	Go(platform *Platform)
}

type OpStop struct{}

func (o OpStop) Go(platform *Platform) {
	os.Exit(0)
}
