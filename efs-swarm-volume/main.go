package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"driver"

	"github.com/docker/go-plugins-helpers/volume"
)

// ensureMount
func mountWithEfs(efsId string, mountpoint string) error {

	cmd := exec.Command("mount")

	cmd.Args = append(cmd.Args, "-t", "efs")
	cmd.Args = append(cmd.Args, "-o", "tls")
	cmd.Args = append(cmd.Args, fmt.Sprintf("%s:/", efsId))
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

	efsId := os.Getenv("EFS_ID")

	err := mountWithEfs(efsId, "/mnt/volumes")
	if err != nil {
		log.Print(err)
	}

	d := driver.NewDriver("/mnt/volumes")

	h := volume.NewHandler(d)

	fmt.Printf("EFS Volume Plugin initializing on %s\n", efsId)
	err = h.ServeUnix("/run/docker/plugins/efsvolumes.sock", 0)

	log.Print(err)

	return
}
