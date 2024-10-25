package docker

import (
	"context"
	"errors"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters" // Importing filters package
	"github.com/docker/docker/client"
)

// ListVolumes retrieves all Docker volumes on the system
func ListVolumes() ([]*types.Volume, error) {
	// Create a new Docker client
	cli, err := client.NewClientWithOpts(client.WithVersion("1.41"))
	if err != nil {
		return nil, err
	}

	// List volumes (no filters applied here)
	volumeList, err := cli.VolumeList(context.Background(), filters.Args{})
	if err != nil {
		return nil, err
	}

	// Return the volume list, which contains pointers to Volume objects
	return volumeList.Volumes, nil
}

func InspectVolume(volumeName string) (*types.Volume, error) {
	// Create a new Docker client
	cli, err := client.NewClientWithOpts(client.WithVersion("1.41"))
	if err != nil {
		return nil, err
	}

	// Inspect the volume using its name
	volume, err := cli.VolumeInspect(context.Background(), volumeName)
	if err != nil {
		return nil, fmt.Errorf("error inspecting volume: %v", err)
	}

	// Return the volume details
	return &volume, nil
}

func ListContainersAttachedToVolume(volumeName string) ([]string, error) {
	// Create a Docker client
	cli, err := client.NewClientWithOpts(client.WithVersion("1.41"))
	if err != nil {
		return nil, err
	}

	// List all containers (including stopped ones)
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		return nil, err
	}

	// Iterate through containers and check if they are using the specified volume
	var containerIDs []string
	for _, container := range containers {
		for _, mount := range container.Mounts {
			if mount.Name == volumeName {
				containerIDs = append(containerIDs, container.ID)
			}
		}
	}

	if len(containerIDs) == 0 {
		return nil, fmt.Errorf("no containers attached to volume: %s", volumeName)
	}

	return containerIDs, nil
}

func RemoveVolume(volumeName string) (string, error) {
	// Create a new Docker client
	cli, err := client.NewClientWithOpts(client.WithVersion("1.41"))
	if err != nil {
		return "", err
	}

	// Check if the volume exists before attempting to remove it
	volume, err := cli.VolumeInspect(context.Background(), volumeName)
	if err != nil {
		if client.IsErrNotFound(err) {
			return "", errors.New("Invalid Volume Name: Volume not found")
		}
		return "", err // Return other errors
	}

	// Attempt to remove the volume
	err = cli.VolumeRemove(context.Background(), volume.Name, true) // true = force remove
	if err != nil {
		return "", err // Return any error encountered
	}

	return "Volume removed successfully", nil
}