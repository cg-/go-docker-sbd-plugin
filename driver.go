package main

import (
	"os"
	//"fmt"
	"github.com/docker/go-plugins-helpers/volume"
	"github.com/cg-/go-nbd"
)

/**
	*	A struct representing a driver for the shared block devices. This is
	* how the main plugin program will create instances of NBDs to connect
	* with the actual block device.
	*/
type fsDriver struct {
	mounts     string
	device		 string
	nbds 			 map[string]*nbd.NbdConnector
}

/**
	*	Constructor.
	* mounts: where the mountpoints should be placed
	* device: the block device to be shared
	*/
func newFsDriver(mounts, device string) fsDriver {
	d := fsDriver{
		mounts: 	 mounts,
		device:		 device,
		nbds:		   make(map[string]*nbd.NbdConnector),
	}
	// create the mount directory if it doesn't exist
	dir := mounts
  _, err := os.Stat(dir)
  if os.IsNotExist(err) {
		os.Mkdir(dir, 0777)
  }
	return d
}

/**
	*	Creates a new volume.
	*
	* This will create the mountpoint on the system.
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
	return volume.Response{Mountpoint: d.mounts + "/" + "0"}
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
	dir := d.mounts + "/" + r.Name
  _, err := os.Stat(dir)
  if os.IsNotExist(err) {
		os.Mkdir(dir, 0777)
  }else{
		return volume.Response{Err: "This mountpoint already exists. Please manully remove it if you're sure it's not being used. Otherwise choose a new name."}
	}

	nbd, err := nbd.CreateNbdConnector(d.device, dir)
	if err != nil {
		return volume.Response{Err: err.Error()}
	}
	d.nbds[r.Name] = nbd
	d.nbds[r.Name].Mount()
	return volume.Response{Mountpoint: dir}
}

/**
	*	Unmounts the volume.
	*/
func (d fsDriver) Unmount(r volume.Request) volume.Response {
	d.nbds[r.Name].Unmount()
	os.RemoveAll(d.mounts + "/" + r.Name)
	return volume.Response{}
}

/**
	* Remounts every volume except the one specified
	*/
func (d fsDriver) RemountAllBut(name string) {
	for k, _ := range d.nbds {
		if k == name {
			continue
		}
		d.nbds[k].Remount()
	}
}

/**
	*	Returns all the info about the volume...
	*/
func (d fsDriver) Get(r volume.Request) volume.Response {
	return volume.Response{Volume: &volume.Volume{Name: r.Name, Mountpoint: d.mounts + "/" + r.Name}}
}

/**
	*	Returns all the info about all volumes. Since we only care about one, this
	* is identical to Get.
	*/
func (d fsDriver) List(r volume.Request) volume.Response {
	volStack := make([]*volume.Volume, 0)
	for k, _ := range d.nbds {
		volStack = append(volStack, &volume.Volume{Name: k, Mountpoint: d.mounts + "/" + k})
	}

	return volume.Response{Volumes: volStack}
}
