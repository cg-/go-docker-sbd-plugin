package main

import (
	"github.com/docker/go-plugins-helpers/volume"
)

type volumeName struct {
	name        string
	connections int
}

type fsDriver struct {
	name       string
}

func newFsDriver(name string) fsDriver {
	d := fsDriver{
		name:    name,
	}

	return d
}

func (d fsDriver) Create(r volume.Request) volume.Response {
	return volume.Response{Err: "Not implemented."}
}

func (d fsDriver) Remove(r volume.Request) volume.Response {
	return volume.Response{Err: "Not implemented."}
}

func (d fsDriver) Path(r volume.Request) volume.Response {
	return volume.Response{Err: "Not implemented."}
}

func (d fsDriver) Mount(r volume.Request) volume.Response {
	return volume.Response{Err: "Not implemented."}
}

func (d fsDriver) Unmount(r volume.Request) volume.Response {
	return volume.Response{Err: "Not implemented."}
}

func (d fsDriver) Get(r volume.Request) volume.Response {
	return volume.Response{Err: "Not implemented."}
}

func (d fsDriver) List(r volume.Request) volume.Response {
	return volume.Response{Err: "Not implemented."}
}

func (d *fsDriver) mountpoint(name string) string {
	return "test"
}

func (d *fsDriver) mountVolume(name, destination string) error {
	return nil
}

func (d *fsDriver) unmountVolume(target string) error {
	return nil
}
