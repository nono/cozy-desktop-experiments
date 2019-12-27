package remote

import (
	"fmt"

	"github.com/nono/cozy-fuse/config"
)

// Do is just an example of using the client.
func Do(cfg *config.Config) {
	cli := cfg.Client()
	d, err := cli.GetDirByPath("/")
	if err != nil {
		fmt.Printf("GetDirByPath error: %s\n", err)
	} else {
		fmt.Printf("GetDirByPath: %#v\n", d)
	}
}
