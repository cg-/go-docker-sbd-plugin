package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/docker/go-plugins-helpers/volume"
)

const (
	suffix = "_sbd"
)

var (
	defaultDir = "/tmp"
	root       = flag.String("root", defaultDir, "Shared block device volumes root directory")
	bd         = flag.String("bd", "", "Shared block device")
)

func main() {
	var Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if *bd == "" {
		Usage()
		os.Exit(1)
	}

	// Create a driver instance for the block device
	d := newFsDriver(*root, *bd)
	h := volume.NewHandler(d)
	fmt.Println(h.ServeUnix("root", "sbd"))

	time.Sleep(10 * time.Second)
	fmt.Println("remounting")
	d.RemountAllBut("none")

	os.Exit(0)
}
