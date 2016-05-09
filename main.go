package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/docker/go-plugins-helpers/volume"
)

const (
	suffix			= "_sbd"
)

var (
	defaultDir  = filepath.Join(volume.DefaultDockerRootDirectory, suffix)
	root        = flag.String("root", defaultDir, "Shared block device volumes root directory")
)

func main() {
	var Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	Usage()

	// Open the block device
	device, err := os.OpenFile("/tmp/test", os.O_RDWR, os.FileMode(0666))

	if(err != nil){
		fmt.Fprintf(os.Stdout, "Trouble opening the block device: %s", err)
		os.Exit(1)
	}
	defer device.Close()

	// Create a driver instance for the block device
	d := newFsDriver("test", device)
	h := volume.NewHandler(d)
	fmt.Println(h.ServeUnix("root", "sbd"))

	os.Exit(0)

	//d := newFsDriver("test")
	//h := volume.NewHandler(d)
	//fmt.Println(h.ServeUnix("root", "glusterfs"))
}
