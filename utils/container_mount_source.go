package utils

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	mountType "github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/gofrs/uuid/v5"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/logging"
)

var NotUniqueError = errors.New("Multiple containers found that could be us")
var NoneFoundError = errors.New("No container found")
var NoMountFound = errors.New("No mount found that matches the data root")
var UnsupportedMount = errors.New("Unsupported mount type found")

var ContainerMountSource string = config.DockerRootPath.Value()

func InitContainerMountSource() (err error) {
	if ContainerMountSource != "" || os.Getenv("PUFFER_PLATFORM") == "" {
		// either the path was set from env or we're not in a container
		return
	}

	path := filepath.Join(os.TempDir(), "puffer-cid")
	err = os.Mkdir(path, 0755)
	if err != nil {
		return
	}
	defer func(path string) {
		_ = os.Remove(path)
	}(path)

	id, err := uuid.NewV4()
	if err != nil {
		return
	}

	path = filepath.Join(path, id.String())
	file, err := os.Create(path)
	if err != nil {
		return
	}

	// we only need the file to exist, we never read or write, so close it right away
	err = file.Close()
	if err != nil {
		return
	}

	docker, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return
	}
	defer func() {
		errC := docker.Close()
		if err == nil {
			err = errC
		}
	}()
	ctx := context.Background()
	docker.NegotiateAPIVersion(ctx)

	containers, err := docker.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		return
	}

	var found []string
	var self types.Container
	for _, c := range containers {
		rc, _, err := docker.CopyFromContainer(ctx, c.ID, path)
		if err != nil {
			// failed, so either file or container doesn't exist, meaning that's not us
			continue
		}
		// not interested in the contents, just need to know the file existed in the container
		rc.Close()
		found = append(found, c.ID)
		self = c
	}

	if len(found) > 1 {
		logging.Debug.Printf("Multiple containers found that could be us: %s\n", strings.Join(found, ", "))
		return NotUniqueError
	}

	if len(found) == 0 {
		return NoneFoundError
	}

	var dataMount *types.MountPoint = nil
	for _, mount := range self.Mounts {
		mountPath, e := filepath.Abs(mount.Destination)
		if e != nil {
			logging.Debug.Println("Failed normalizing mount destination path, trying without normalizing")
			mountPath = mount.Destination
		}

		dataRoot, e := filepath.Abs(config.DataRootFolder.Value())
		if e != nil {
			logging.Debug.Println("Failed normalizing data root path, trying without normalizing")
			mountPath = config.DataRootFolder.Value()
		}

		if mountPath == dataRoot {
			dataMount = &mount
			break
		}
	}

	if dataMount == nil {
		return NoMountFound
	}

	if dataMount.Type != mountType.TypeBind && dataMount.Type != mountType.TypeVolume {
		logging.Debug.Printf("Unsupported mount type found: %s\n", dataMount.Type)
		return UnsupportedMount
	}

	ContainerMountSource = dataMount.Source
	logging.Debug.Printf("Found own container id %s and host path %s\n", self.ID, dataMount.Source)
	return
}
