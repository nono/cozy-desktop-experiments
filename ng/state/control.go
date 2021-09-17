package state

import (
	"os"
)

type EventStart struct{}

func (e EventStart) Update(state *State) []Operation {
	return []Operation{OpStat{"."}}
}

type OpStop struct{}

func (o OpStop) Go(platform Platform) {
	os.Exit(0)
}
