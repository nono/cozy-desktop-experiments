package main

import (
	"fmt"

	"github.com/nono/cozy-desktop-experiments/ng/sync"
)

func main() {
	fmt.Println("Start")
	sync.Start(".")
}
