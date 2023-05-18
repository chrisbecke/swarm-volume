package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"driver"

	"github.com/docker/go-plugins-helpers/volume"
)

// ensureMount sudo mount -t nfs4 -o nfsvers=4.1,rsize=1048576,wsize=1048576,hard,timeo=600,retrans=2,noresvport fs-000ec294979105cf3.efs.eu-west-1.amazonaws.com:/ efs
func mountWithNfs(mountType string, mountOptions string, mountDevice string, mountpoint string) error {

	cmd := exec.Command("mount")

	cmd.Args = append(cmd.Args, "-t", mountType)
	cmd.Args = append(cmd.Args, "-o", mountOptions)
	cmd.Args = append(cmd.Args, mountDevice)
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

	nfsType := os.Getenv("NFS_TYPE")
	nfsOptions := os.Getenv("NFS_OPTIONS")
	nfsDevice := os.Getenv("NFS_DEVICE")

	err := mountWithNfs(nfsType, nfsOptions, nfsDevice, "/mnt/volumes")
	if err != nil {
		log.Print(err)
	}
	defer unmount("/mnt/volumes")

	d := driver.NewDriver("/mnt/volumes")

	h := volume.NewHandler(d)

	fmt.Printf("NFS Volume Plugin initializing on %s\n", nfsDevice)
	err = h.ServeUnix("/run/docker/plugins/nfsvol.sock", 0)

	log.Print(err)

	return
}
