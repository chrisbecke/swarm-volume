package driver

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/docker/go-plugins-helpers/volume"
)

//------------------------------

const showHidden = false

func NewDriver(root string) *simplefsDriver {

	d := &simplefsDriver{
		mounts: map[string]*activeMount{},
		root:   root,
	}
	return d
}

///////////////////////////////////////////////////////////////////////////////

// ActiveMount holds active mounts
type activeMount struct {
	connections int
	//	mountpoint  string
	//	createdAt   time.Time
	ids map[string]int
}

type simplefsDriver struct {
	sync.RWMutex

	root string

	mounts map[string]*activeMount
}

// API volumeDriver.Create
func (d *simplefsDriver) Create(r *volume.CreateRequest) error {
	d.Lock()
	defer d.Unlock()

	err := os.MkdirAll(d.mountpoint(r.Name), 0755)

	return err
}

// volumeDriver.List
func (d *simplefsDriver) List() (*volume.ListResponse, error) {
	d.Lock()
	defer d.Unlock()

	files, err := ioutil.ReadDir(d.root)

	if err != nil {
		return &volume.ListResponse{}, err
	}

	var vols []*volume.Volume
	for _, file := range files {
		if file.IsDir() && (showHidden || !strings.HasPrefix(file.Name(), ".")) {

			vols = append(vols, &volume.Volume{Name: file.Name()})
		}
	}

	return &volume.ListResponse{Volumes: vols}, nil
}

// volumeDriver.Get
func (d *simplefsDriver) Get(r *volume.GetRequest) (*volume.GetResponse, error) {
	d.Lock()
	defer d.Unlock()

	s := make(map[string]interface{})

	stat, err := os.Stat(d.mountpoint(r.Name))

	if err != nil {
		return &volume.GetResponse{}, err
	}
	vo := &volume.Volume{
		Name:      stat.Name(),
		CreatedAt: stat.ModTime().Format(time.RFC3339),
		Status:    s,
	}
	return &volume.GetResponse{Volume: vo}, nil
}

// volumeDriver.Remove
func (d *simplefsDriver) Remove(r *volume.RemoveRequest) error {
	d.Lock()
	defer d.Unlock()

	v, ok := d.mounts[r.Name]
	if ok && v.connections != 0 {
		log.Printf("Error: %d Existing local mounts", v.connections)
	}

	err := os.RemoveAll(d.mountpoint(r.Name))

	delete(d.mounts, r.Name)

	return err
}

// Volumedriver.Path
func (d *simplefsDriver) Path(r *volume.PathRequest) (*volume.PathResponse, error) {
	d.Lock()
	defer d.Unlock()

	v, ok := d.mounts[r.Name]
	if !ok || v.connections == 0 {
		err := fmt.Errorf("no mountpoint for volume.")
		log.Printf("Path error. name: %s, err: %v", r.Name, err)
		return &volume.PathResponse{}, err
	}

	return &volume.PathResponse{Mountpoint: d.mountpoint(r.Name)}, nil
}

// VolumeDriver.Mount
func (d *simplefsDriver) Mount(r *volume.MountRequest) (*volume.MountResponse, error) {
	d.Lock()
	defer d.Unlock()

	mountpoint := d.mountpoint(r.Name)

	v, ok := d.mounts[r.Name]
	if !ok {
		v = &activeMount{
			ids: map[string]int{},
		}
		d.mounts[r.Name] = v
	}

	v.ids[r.ID]++
	v.connections++

	log.Printf("Mounted registration: %+v", v)

	return &volume.MountResponse{Mountpoint: mountpoint}, nil
}

// VolumeDriver.Unmount
func (d *simplefsDriver) Unmount(r *volume.UnmountRequest) error {
	log.Printf("Unmount %v", r)

	v, ok := d.mounts[r.Name]
	if !ok {
		err := fmt.Errorf("Volume not found in active Mounts: %s", r.Name)
		log.Printf("Unmount failed: %v", err)
		return err
	}

	if v.connections == 0 {
		err := fmt.Errorf("Mount has no active connections: %s", r.Name)
		log.Printf("Unmount failed: %v", err)
		return err
	}

	i, ok := v.ids[r.ID]
	if !ok {
		err := fmt.Errorf("Mount %s does not know about this client ID: %s", r.Name, r.ID)
		log.Printf("Unmount failed: %v", err)
		return err
	}

	i--
	v.connections--

	if i <= 1 {
		delete(v.ids, r.ID)
	} else {
		v.ids[r.ID] = i
	}

	if len(v.ids) == 0 {
		log.Printf("Unmounting volume %s with %v clients", r.Name, v.connections)

		delete(d.mounts, r.Name)
	}

	return nil
}

func (d *simplefsDriver) Capabilities() *volume.CapabilitiesResponse {
	return &volume.CapabilitiesResponse{Capabilities: volume.Capability{Scope: "global"}}
}

// mountpoint of a docker volume
func (d *simplefsDriver) mountpoint(Name string) string {
	return filepath.Join(d.root, Name)
}
