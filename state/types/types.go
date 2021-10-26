// Package types is just for sharing some basic types between packages without
// introducing a loop of imports.
package types

// Clock is a logical clock for the state.
type Clock uint64

// Type is used to differentiate files to directories.
type Type int

const (
	UnknownType Type = iota
	FileType
	DirType
)
