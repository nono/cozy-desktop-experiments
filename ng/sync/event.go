package sync

import (
	"fmt"
	"io/fs"
)

type Event interface {
	Update(state *State) []Operation
}

type EventStart struct{}

func (e EventStart) Update(state *State) []Operation {
	return []Operation{OpStat{"."}}
}

type EventStatDone struct {
	Op    OpStat
	Info  fs.FileInfo
	Error error
}

func (e EventStatDone) Update(state *State) []Operation {
	fmt.Printf("Update %#v\n", e.Info)
	if e.Op.Path == "." && e.Error == nil && e.Info.IsDir() {
		return []Operation{OpScan{"."}}
	}
	return []Operation{OpStop{}}
}

type EventScanDone struct {
	Op      OpScan
	Entries []fs.DirEntry
	Error   error
}

func (e EventScanDone) Update(state *State) []Operation {
	fmt.Printf("Update %#v\n", e.Entries)
	return []Operation{OpStop{}}
}
