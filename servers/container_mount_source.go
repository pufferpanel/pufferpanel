package servers

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofrs/uuid/v5"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/mount"
	"github.com/moby/moby/client"

	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"github.com/pufferpanel/pufferpanel/v3/utils"
)

var containerMountSource string

func InitContainerMountSource() (err error) {
	containerMountSource = config.DockerRootPath.Value()
	if containerMountSource != "" || os.Getenv("PUFFER_PLATFORM") == "" {
		// either the path was set from env or we're not in a container
		return
	}

	path := filepath.Join(os.TempDir(), "puffer-cid")
	err = os.MkdirAll(path, 0755)
	if err != nil {
		return
	}
	defer os.RemoveAll(path)

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
	utils.Close(file)

	docker, err := client.New(client.FromEnv)
	if err != nil {
		return
	}
	defer utils.Close(docker)
	ctx := context.Background()

	containers, err := docker.ContainerList(ctx, client.ContainerListOptions{})
	if err != nil {
		return
	}

	var found []string
	var self container.Summary
	for _, c := range containers.Items {
		rc, err := docker.CopyFromContainer(ctx, c.ID, client.CopyFromContainerOptions{
			SourcePath: path,
		})
		if err != nil {
			// failed, so either file or container doesn't exist, meaning that's not us
			continue
		}
		// not interested in the contents, just need to know the file existed in the container
		utils.Close(rc.Content)
		found = append(found, c.ID)
		self = c
	}

	if len(found) > 1 {
		logging.Debug.Printf("Multiple containers found that could be us: %s\n", strings.Join(found, ", "))
		return pufferpanel.ErrContainerNotUnique
	}

	if len(found) == 0 {
		return pufferpanel.ErrNoContainerFound
	}

	var dataMount *container.MountPoint = nil
	for _, mount := range self.Mounts {
		mountPath, e := filepath.Abs(mount.Destination)
		if e != nil {
			logging.Debug.Println("Failed normalizing mount destination path, trying without normalizing")
			mountPath = mount.Destination
		}

		dataRoot, e := filepath.Abs(config.DataRootFolder.Value())
		if e != nil {
			logging.Debug.Println("Failed normalizing data root path, trying without normalizing")
			dataRoot = config.DataRootFolder.Value()
		}

		if mountPath == dataRoot {
			dataMount = &mount
			break
		}
	}

	if dataMount == nil {
		return pufferpanel.ErrNoMountFound
	}

	if dataMount.Type != mount.TypeBind && dataMount.Type != mount.TypeVolume {
		logging.Debug.Printf("Unsupported mount type found: %s\n", dataMount.Type)
		return pufferpanel.ErrUnsupportedMountType
	}

	containerMountSource = dataMount.Source
	logging.Debug.Printf("Found own container id %s and host path %s\n", self.ID, dataMount.Source)
	return
}
