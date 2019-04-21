package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"bazil.org/fuse"
	"github.com/nono/cozy-fuse/config"
	"github.com/nono/cozy-fuse/sync"
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

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt)
		<-signalChan
		cancel()
	}()

	fmt.Println("cozy-fuse: mounting")
	if err := sync.Run(ctx, cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
