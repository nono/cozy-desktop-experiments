package remote

const (
	// RootID is the id of the CouchDB document that represents the root on the
	// Cozy.
	RootID = "io.cozy.files.root-dir"
	// TrashID is the id of the CouchDB document that represents the trash on
	// the Cozy.
	TrashID = "io.cozy.files.trash-dir"
	// TrashName is the name of the trash folder on the Cozy.
	TrashName = ".cozy_trash"
	// Directory is the type of directories for the CouchDB documents on the
	// Cozy.
	Directory = "directory"
	// File is the type for, well, files for the CouchDB documents on the Cozy.
	File = "file"
)
