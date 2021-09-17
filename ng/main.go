package main

import (
	"fmt"

	"github.com/nono/cozy-desktop-experiments/ng/local"
	"github.com/nono/cozy-desktop-experiments/ng/platform"
	"github.com/nono/cozy-desktop-experiments/ng/state"
)

func main() {
	fmt.Println("Start")
	localDir := "."
	localFS := local.DirFS(localDir)
	// localFS := local.MemFS()
	platform := platform.New(localFS)
	state.Sync(platform)
}
