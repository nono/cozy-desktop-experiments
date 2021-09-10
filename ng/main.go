package main

import (
	"fmt"
	"os"

	"github.com/nono/cozy-desktop-experiments/ng/sync"
)

func main() {
	fmt.Println("Start")
	localDir := "."
	local := os.DirFS(localDir)
	platform := sync.NewPlatform(local)
	sync.Start(platform)
}
