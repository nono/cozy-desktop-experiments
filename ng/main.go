package main

import (
	"fmt"

	"github.com/nono/cozy-desktop-experiments/ng/sync"
)

func main() {
	fmt.Println("Start")
	localDir := "."
	local := sync.DirFS(localDir)
	platform := sync.NewPlatform(local)
	sync.Start(platform)
}
