package main

import (
	"fmt"
	"os"

	"github.com/nono/cozy-desktop-experiments/ng/localfs"
	"github.com/nono/cozy-desktop-experiments/ng/platform"
	"github.com/nono/cozy-desktop-experiments/ng/state"
)

func main() {
	fmt.Println("Start")
	localDir := "."
	localFS := localfs.DirFS(localDir)
	// localFS := localfs.MemFS()
	platform := platform.New(localFS)
	if err := state.Sync(platform); err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Done.\n")
}
