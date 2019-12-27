package local

import (
	"context"
	"fmt"
	"os"

	"bazil.org/fuse"
	fspkg "bazil.org/fuse/fs"
	"github.com/nono/cozy-fuse/config"
)

// Mount will mount the file system and will serve in a goroutine the fuse requests.
func Mount(ctx context.Context, cfg *config.Config) error {
	options := []fuse.MountOption{
		fuse.FSName("cozy-fuse"),
		fuse.Subtype("cozy-fuse"),
		fuse.ReadOnly(),
	}
	c, err := fuse.Mount(cfg.Mount, options...)
	if err != nil {
		return err
	}
	// XXX c.Close() blocks if fuse.Unmount has not been called before
	defer c.Close()
	defer fuse.Unmount(cfg.Mount)
	ch := make(chan error)
	go serve(c, ch)
	select {
	case err = <-ch:
		return err
	case <-ctx.Done():
	}
	return nil
}

func serve(c *fuse.Conn, ch chan error) {
	fs := &FS{}
	err := fspkg.Serve(c, fs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error on serve: %s\n", err)
		ch <- err
	}
	<-c.Ready
	if err := c.MountError; err != nil {
		fmt.Fprintf(os.Stderr, "Error on ready: %s\n", err)
		ch <- err
	}
}
