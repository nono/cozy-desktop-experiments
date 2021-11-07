package main

import (
	"fmt"
	"os"

	"github.com/nono/cozy-desktop-experiments/client"
	"github.com/nono/cozy-desktop-experiments/localfs"
	"github.com/nono/cozy-desktop-experiments/platform"
	"github.com/nono/cozy-desktop-experiments/state"
)

func main() {
	localDir := "."
	localFS := localfs.NewDirFS(localDir)
	// localFS := localfs.NewMemFS()

	// remoteClient := client.New("http://cozy.localhost:8080/")
	remoteClient := client.NewFake("http://cozy.localhost:8080/")
	remoteClient.AddInitialDocs()

	fmt.Println("Start")
	platform := platform.New(localFS, remoteClient)
	if err := state.Sync(platform); err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Done.\n")
}
