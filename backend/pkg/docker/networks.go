package docker

import (
	"context"
	"errors"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// ListNetworks retrieves all Docker networks with ID, Name, Driver, and Scope
func ListNetworks() ([]map[string]string, error) {
	// Create a new Docker client
	cli, err := client.NewClientWithOpts(client.WithVersion("1.41"))
	if err != nil {
		return nil, err
	}

	// List all networks
	networks, err := cli.NetworkList(context.Background(), types.NetworkListOptions{})
	if err != nil {
		return nil, err
	}

	// Format response to include ID, Name, Driver, and Scope
	networkDetails := []map[string]string{}
	for _, network := range networks {
		networkDetails = append(networkDetails, map[string]string{
			"id":     network.ID,
			"name":   network.Name,
			"driver": network.Driver,
			"scope":  network.Scope,
		})
	}

	return networkDetails, nil
}

func InspectNetwork(networkID string) (map[string]interface{}, error) {
	// Create a new Docker client
	cli, err := client.NewClientWithOpts(client.WithVersion("1.41"))
	if err != nil {
		return nil, err
	}

	// Inspect the specified network
	networkResource, err := cli.NetworkInspect(context.Background(), networkID, types.NetworkInspectOptions{})
	if err != nil {
		return nil, errors.New("network not found")
	}

	// Prepare the network details in response format
	networkDetails := map[string]interface{}{
		"id":         networkResource.ID,
		"name":       networkResource.Name,
		"driver":     networkResource.Driver,
		"scope":      networkResource.Scope,
		"containers": networkResource.Containers, // Shows attached containers
	}

	return networkDetails, nil
}

type NetworkContainer struct {
	ContainerID string `json:"container_id"`
	Name        string `json:"name"`
}

func ListContainersInNetwork(networkID string) ([]NetworkContainer, error) {
	// Create a new Docker client
	cli, err := client.NewClientWithOpts(client.WithVersion("1.41"))
	if err != nil {
		return nil, err
	}

	// Inspect the network to retrieve details of attached containers
	networkResource, err := cli.NetworkInspect(context.Background(), networkID, types.NetworkInspectOptions{})
	if err != nil {
		return nil, err
	}

	// Collect attached containers
	containers := []NetworkContainer{}
	for containerID, container := range networkResource.Containers {
		containers = append(containers, NetworkContainer{
			ContainerID: containerID,
			Name:        container.Name,
		})
	}

	// Return an error message if no containers are attached
	if len(containers) == 0 {
		return nil, fmt.Errorf("No containers are attached to the network: %s", networkID)
	}

	return containers, nil
}

func RemoveNetwork(networkID string) (string, error) {
    // Create a new Docker client
    cli, err := client.NewClientWithOpts(client.WithVersion("1.41"))
    if err != nil {
        return "", err
    }

    // Try to remove the network
    if err := cli.NetworkRemove(context.Background(), networkID); err != nil {
        // Handle network not found error
        if client.IsErrNotFound(err) {
            return "", errors.New("Network not found. Invalid Network ID.")
        }
        return "", err
    }

    return "Network removed successfully", nil
}