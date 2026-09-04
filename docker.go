package pufferpanel

import (
	"context"
	"io"
	"slices"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"github.com/pufferpanel/pufferpanel/v3/utils"
)

var dockerClient *client.Client

func GetDockerClient() (*client.Client, error) {
	var err error = nil
	if dockerClient == nil {
		dockerClient, err = client.NewClientWithOpts(client.FromEnv)
		ctx := context.Background()
		dockerClient.NegotiateAPIVersion(ctx)
	}
	return dockerClient, err
}

func DoesContainerExist(id string, ctx context.Context) (bool, error) {
	client, err := GetDockerClient()
	if err != nil {
		return false, err
	}

	opts := container.ListOptions{
		Filters: filters.NewArgs(),
	}

	opts.All = true
	opts.Filters.Add("name", id)

	existingContainers, err := client.ContainerList(ctx, opts)
	if err != nil {
		return false, err
	}

	for _, v := range existingContainers {
		if slices.Contains(v.Names, "/"+id) {
			return true, nil
		}
	}

	return false, nil
}

func PullDockerImage(environment *Environment, ctx context.Context, imageName string, force bool) error {
	client, err := GetDockerClient()
	if err != nil {
		return err
	}

	if !force {
		exists := false

		parts := strings.SplitN(imageName, ":", 2)
		if len(parts) != 2 {
			imageName = imageName + ":latest"
		}

		opts := image.ListOptions{
			All:     true,
			Filters: filters.NewArgs(),
		}
		opts.Filters.Add("reference", imageName)
		images, err := client.ImageList(ctx, opts)

		if err != nil {
			return err
		}

		for _, v := range images {
			for _, z := range v.RepoTags {
				if z == imageName {
					exists = true
					break
				}
			}
			if exists {
				break
			}
		}

		environment.Log(logging.Debug, "Does image %v exist? %v", imageName, exists)

		if exists {
			return nil
		}
	}

	op := image.PullOptions{}

	environment.Log(logging.Debug, "Downloading image %v", imageName)
	environment.DisplayToConsole(true, "Downloading image for container, please wait\n")

	r, err := client.ImagePull(ctx, imageName, op)
	defer utils.Close(r)
	if err != nil {
		return err
	}

	w := &ImageWriter{Parent: environment.ConsoleTracker}
	_, err = io.Copy(w, r)

	if err != nil {
		return err
	}

	environment.Log(logging.Debug, "Downloaded image %v", imageName)
	environment.DisplayToConsole(true, "Downloaded image for container\n")
	return err
}
