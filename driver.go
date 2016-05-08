package main

import (
	"fmt"
	"os"
	"time"
	"syscall"

	"github.com/docker/go-plugins-helpers/volume"
)

/**
	*	A struct representing a driver for the shared block devices. This is
	* how the main plugin program will create instances of NBDs to connect
	* with the actual block device.
	*/
type fsDriver struct {
	name       string
	mountedAt	 string
	device		 *os.File
	nbds 			 map[int]*NBD
}

/**
	*	Constructor.
	* name: a name for the device
	* device: the block device to be shared
	*/
func newFsDriver(name string, device *os.File) fsDriver {
	d := fsDriver{
		name:    	 name,
		mountedAt: "",
		device:		 device,
		nbds:		   make(map[int]*NBD),
	}

	return d
}

/**
	*	Creates a new volume.
	*
	* Since we are using a volume that already exists, we'll just have this
	* return affirmatively.
	*/
func (d fsDriver) Create(r volume.Request) volume.Response {
	return volume.Response{}
}

/**
	*	Deletes a volume.
	*
	* Again, since we're using a volume that already exists, we don't want
	* the plugin to have this functionality, so we'll just have it return
	* affirmatively.
	*/
func (d fsDriver) Remove(r volume.Request) volume.Response {
	return volume.Response{}
}

/**
	*	Returns the path to the mountpoint on the host machine.
	*/
func (d fsDriver) Path(r volume.Request) volume.Response {
	// make sure it's mounted first...
	if(d.mountedAt != ""){
		return volume.Response{Err: "Not mounted."}
	}

	return volume.Response{Mountpoint: d.mountedAt}
}

/**
	*	Mounts the shared block device onto the machine
	*
	*	Mount needs to handle a lot of the heavy lifting. It's main functions are:
	*		- Create a new NBD device linked to the actual block device
	*		- Create a new mount point on the host machine's OS
	*		-	Mount the new NBD to the new mount point
	* 	- Let Docker know where to access the volume
	*/
func (d fsDriver) Mount(r volume.Request) volume.Response {
	stat, _ := d.device.Stat()
	dev := Create(d.device, stat.Size(), make(chan bool))

	// fire up a thread to handle the new nbd
	go dev.Connect()

	// wait for the nbd to start, and figure out what nbd it's mounted at...
	id := 0
	for {
		time.Sleep(1)
		id = dev.GetIdent()
		if(id != 0){
			d.nbds[id] = dev
			break
		}
	}

	// create a new directory to mount the nbd if it doesn't exist
	devpath := fmt.Sprintf("/dev/nbd%d", id)
	path := fmt.Sprintf("/tmp/dev_%d", id)
	_, err := os.Stat(path)
	if(os.IsNotExist(err)){
		os.Mkdir(path, 0777)
	}

	// try to mount it
	err = syscall.Mount(devpath, path, "ext3", 0xC0ED, "rw")
	if(err != nil){
		return volume.Response{Err: err.Error()}
	}

	// set the mountedAt value
	d.mountedAt = path

	// send it to docker
	return volume.Response{Mountpoint: d.mountedAt}
}

/**
	*	Unmounts the volume. This is going to be tricky since it depends on what
	* container is making the call.
	*/
func (d fsDriver) Unmount(r volume.Request) volume.Response {
	return volume.Response{Err: "Not implemented."}
}

/**
	*	Returns all the info about the volume...
	*/
func (d fsDriver) Get(r volume.Request) volume.Response {
	return volume.Response{Volume: &volume.Volume{Name: d.name, Mountpoint: d.mountedAt}}
}

/**
	*	Returns all the info about all volumes. Since we only care about one, this
	* is identical to Get.
	*/
func (d fsDriver) List(r volume.Request) volume.Response {
	return d.Get(r)
}

/**
	*
	*/
func (d fsDriver) Dump() string {
	var toReturn string
	toReturn = fmt.Sprintf("name: %s", d.name)
	return toReturn
}
