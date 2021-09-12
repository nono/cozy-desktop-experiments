package sync

import "os"

type Operation interface {
	Go(platform *Platform)
}

type OpStop struct{}

func (o OpStop) Go(platform *Platform) {
	os.Exit(0)
}

type OpStat struct {
	Path string
}

func (o OpStat) Go(platform *Platform) {
	go func() {
		info, err := platform.Local.Stat(o.Path)
		platform.Events <- EventStatDone{Op: o, Info: info, Error: err}
	}()
}

type OpScan struct {
	Path string
}

func (o OpScan) Go(platform *Platform) {
	go func() {
		entries, err := platform.Local.ReadDir(o.Path)
		platform.Events <- EventScanDone{Op: o, Entries: entries, Error: err}
	}()
}
