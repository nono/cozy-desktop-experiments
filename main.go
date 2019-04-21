package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"bazil.org/fuse"
	fspkg "bazil.org/fuse/fs"
	"github.com/nono/cozy-fuse/config"
	"github.com/nono/cozy-fuse/local"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: cozy-fuse <MOUNT_POINT>\n")
	flag.PrintDefaults()
}

func main() {
	fuseDebug := flag.Bool("debug", false, "enable verbose FUSE debugging")
	configPath := flag.String("config", "config.json", "use this configuration file")
	flag.Parse()
	if flag.NArg() > 0 {
		usage()
		os.Exit(2)
	}
	cfg, err := config.Load(*configPath)
	if err != nil {
		if err == config.ErrNotExist {
			usage()
		} else {
			fmt.Fprintf(os.Stderr, "Error on config: %s\n", err)
		}
		os.Exit(2)
	}
	if *fuseDebug {
		fuse.Debug = func(msg interface{}) {
			fmt.Printf("fuse debug: %v\n", msg)
		}
	}

	fmt.Printf("config = %#v\n", cfg)

	fmt.Println("cozy-fuse: mounting")
	options := []fuse.MountOption{
		fuse.FSName("cozy-fuse"),
		fuse.Subtype("cozy-fuse"),
		fuse.ReadOnly(),
	}
	c, err := fuse.Mount(cfg.Mount, options...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error on mount: %s\n", err)
		os.Exit(1)
	}
	defer c.Close()
	defer fuse.Unmount(cfg.Mount)

	fs := &local.FS{}
	err = fspkg.Serve(c, fs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error on serve: %s\n", err)
		os.Exit(1)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	select {
	case <-c.Ready:
	case <-signalChan:
	}
	if err := c.MountError; err != nil {
		fmt.Fprintf(os.Stderr, "Error on ready: %s\n", err)
		os.Exit(1)
	}
}
