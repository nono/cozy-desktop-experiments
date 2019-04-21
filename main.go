package main

import (
	"flag"
	"fmt"
	"os"

	"bazil.org/fuse"
	fspkg "bazil.org/fuse/fs"
	"github.com/nono/cozy-fuse/local"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: cozy-fuse <MOUNT_POINT>\n")
	flag.PrintDefaults()
}

func main() {
	fuseDebug := flag.Bool("debug", false, "enable verbose FUSE debugging")
	flag.Parse()
	if flag.NArg() != 1 {
		usage()
		os.Exit(2)
	}
	if *fuseDebug {
		fuse.Debug = func(msg interface{}) {
			fmt.Printf("fuse debug: %v\n", msg)
		}
	}
	mntPoint := flag.Arg(0)

	fmt.Println("cozy-fuse: mounting")
	options := []fuse.MountOption{
		fuse.FSName("cozy-fuse"),
		fuse.Subtype("cozy-fuse"),
		fuse.ReadOnly(),
	}
	c, err := fuse.Mount(mntPoint, options...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error on mount: %s\n", err)
		os.Exit(1)
	}
	defer c.Close()
	defer fuse.Unmount(mntPoint)

	fs := &local.FS{}
	err = fspkg.Serve(c, fs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error on serve: %s\n", err)
		os.Exit(1)
	}

	<-c.Ready
	if err := c.MountError; err != nil {
		fmt.Fprintf(os.Stderr, "Error on ready: %s\n", err)
		os.Exit(1)
	}
}
