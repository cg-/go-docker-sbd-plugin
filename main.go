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

	os.Exit(0)

	//d := newFsDriver("test")
	//h := volume.NewHandler(d)
	//fmt.Println(h.ServeUnix("root", "glusterfs"))
}
