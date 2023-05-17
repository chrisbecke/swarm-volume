package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"driver"

	"github.com/docker/go-plugins-helpers/volume"
)

//------------------------------

// config.json settings
const socketAddress = "/run/docker/plugins/glusterfs.sock"
const propagatedMount = "/mnt/volumes"

// ensureMount
func mountWithGlusterfs(hosts []string, volume string, mountpoint string) error {

	cmd := exec.Command("glusterfs")

	for _, server := range hosts {
		cmd.Args = append(cmd.Args, "--volfile-server", server)
	}

	cmd.Args = append(cmd.Args, "--volfile-id", volume)
	cmd.Args = append(cmd.Args, mountpoint)

	fmt.Printf("Executing %v\n", cmd)

	_, err := cmd.CombinedOutput()

	return err
}

func unmount(mountpoint string) error {

	cmd := exec.Command("umount", mountpoint)
	_, err := cmd.CombinedOutput()

	return err
}

func main() {

	gfsvol := os.Getenv("GFS_VOLUME")
	gfsservers := strings.Split(os.Getenv("GFS_SERVERS"), ",")

	err := mountWithGlusterfs(gfsservers, gfsvol, propagatedMount)

	if err != nil {
		log.Print(err)
		return
	}
	defer unmount(propagatedMount)

	d := driver.NewDriver(propagatedMount)

	h := volume.NewHandler(d)

	fmt.Printf("GlusterFS Volume Plugin %s connecting to %v\n", gfsvol, gfsservers)
	err = h.ServeUnix(socketAddress, 0)

	log.Print(err)

	return
}
