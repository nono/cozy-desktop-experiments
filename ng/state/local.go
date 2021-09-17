package state

import (
	"fmt"
	"io/fs"
	"syscall"
)

type OpStat struct {
	Path string
}

func (o OpStat) Go(platform Platform) {
	go func() {
		info, err := platform.FS().Stat(o.Path)
		platform.Notify(EventStatDone{Op: o, Info: info, Error: err})
	}()
}

type EventStatDone struct {
	Op    OpStat
	Info  fs.FileInfo
	Error error
}

func (e EventStatDone) Update(state *State) []Operation {
	fmt.Printf("Update %#v\n", e.Info)
	if e.Op.Path == "." && e.Error == nil && e.Info.IsDir() {
		fmt.Printf("inode number = %v\n", e.Info.Sys().(*syscall.Stat_t).Ino)
		return []Operation{OpScan{"."}}
	}
	return []Operation{OpStop{}}
}

type OpScan struct {
	Path string
}

func (o OpScan) Go(platform Platform) {
	go func() {
		entries, err := platform.FS().ReadDir(o.Path)
		platform.Notify(EventScanDone{Op: o, Entries: entries, Error: err})
	}()
}

type EventScanDone struct {
	Op      OpScan
	Entries []fs.DirEntry
	Error   error
}

func (e EventScanDone) Update(state *State) []Operation {
	fmt.Printf("Update\n")
	for _, e := range e.Entries {
		fmt.Printf("* %#v\n", e)
	}
	return []Operation{OpStop{}}
}
