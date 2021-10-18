package main

import (
	"fmt"
	"os"

	"github.com/nono/cozy-desktop-experiments/ng/client"
	"github.com/nono/cozy-desktop-experiments/ng/localfs"
	"github.com/nono/cozy-desktop-experiments/ng/platform"
	"github.com/nono/cozy-desktop-experiments/ng/state"
)

func main() {
	localDir := "."
	localFS := localfs.NewDirFS(localDir)
	// localFS := localfs.NewMemFS()

	remoteClient := client.New("http://cozy.localhost:8080/")
	// remoteClient := client.NewFake("http://cozy.localhost:8080/")

	fmt.Println("Start")
	platform := platform.New(localFS, remoteClient)
	if err := state.Sync(platform); err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Done.\n")
}
