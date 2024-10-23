package docker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

// ListContainers retrieves the list of Docker containers
func ListContainers() ([]types.Container, error) {
	// Create a new Docker client
	cli, err := client.NewClientWithOpts(client.WithVersion("1.41"))
	if err != nil {
		return nil, err
	}

	// Set up the options for listing only running containers
	options := types.ContainerListOptions{
		All:     false, // Only list running containers
		Filters: filters.NewArgs(),
	}

	// Add filter for running containers
	options.Filters.Add("status", "running")

	// List containers
	containers, err := cli.ContainerList(context.Background(), options)
	if err != nil {
		return nil, err
	}

	return containers, err
}

func ListAllContainers() ([]types.Container, error) {
	cli, err := client.NewClientWithOpts(client.WithVersion("1.41"))
	if err != nil {
		return nil, err
	}

	options := types.ContainerListOptions{
		All: true, // List all containers
	}

	containers, err := cli.ContainerList(context.Background(), options)
	if err != nil {
		return nil, err
	}

	return containers, err
}

// StartContainer starts a Docker container if it is not already running
func StartContainer(containerID string) (string, error) {
	// Create a new Docker client
	cli, err := client.NewClientWithOpts(client.WithVersion("1.41"))
	if err != nil {
		return "", err
	}

	// Check the current status of the container
	containerJSON, err := cli.ContainerInspect(context.Background(), containerID)
	if err != nil {
		return "", err
	}

	// Check if the container is already running
	if containerJSON.State.Running {
		return "Container is already running", nil
	}

	// Start the container
	if err := cli.ContainerStart(context.Background(), containerID, types.ContainerStartOptions{}); err != nil {
		return "", errors.New("failed to start the container: " + err.Error())
	}

	// Check the container's state after starting it
	containerJSON, err = cli.ContainerInspect(context.Background(), containerID)
	if err != nil {
		return "", err
	}

	// Check if the container is running
	if !containerJSON.State.Running {
		// If the container is not running, check the exit code
		if containerJSON.State.ExitCode != 0 {
			return "Failed to start the container, it exited with code " + string(containerJSON.State.ExitCode), nil
		}
		return "Failed to start the container, it exited immediately", nil
	}

	return "Container started successfully", nil
}

func StopContainer(containerID string) (string, error) {
	// Create a new Docker client
	cli, err := client.NewClientWithOpts(client.WithVersion("1.41"))
	if err != nil {
		return "", err
	}

	// Check the current status of the container
	containerJSON, err := cli.ContainerInspect(context.Background(), containerID)
	if err != nil {
		return "", err
	}

	// Check if the container is not running
	if !containerJSON.State.Running {
		return "Specified container is not running", nil
	}

	// Stop the container
	if err := cli.ContainerStop(context.Background(), containerID, nil); err != nil {
		return "", errors.New("failed to stop the container: " + err.Error())
	}

	return "Container stopped successfully", nil
}

func RemoveContainer(containerID string) (string, error) {
	// Create a new Docker client
	cli, err := client.NewClientWithOpts(client.WithVersion("1.41"))
	if err != nil {
		return "", err
	}

	// Check if the container exists
	_, err = cli.ContainerInspect(context.Background(), containerID)
	if err != nil {
		if client.IsErrNotFound(err) {
			return "", errors.New("Invalid Container ID")
		}
		return "", err
	}

	// Remove the container
	if err := cli.ContainerRemove(context.Background(), containerID, types.ContainerRemoveOptions{Force: true}); err != nil {
		return "", err
	}

	return "Container removed successfully", nil
}

func RemoveAllContainers() ([]string, error) {
	// Create a new Docker client
	cli, err := client.NewClientWithOpts(client.WithVersion("1.41"))
	if err != nil {
		return nil, err
	}

	// Get the list of all containers (including stopped containers)
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		return nil, err
	}

	// Slice to hold the result messages
	var results []string

	// Iterate over each container and try to remove it
	for _, container := range containers {
		// Check if the container is running
		containerJSON, err := cli.ContainerInspect(context.Background(), container.ID)
		if err != nil {
			results = append(results, fmt.Sprintf("Failed to inspect container: %s", container.ID))
			continue
		}

		// If the container is running, skip deletion and add a message
		if containerJSON.State.Running {
			results = append(results, fmt.Sprintf("Cannot delete the container, it's in running state: %s", container.ID))
		} else {
			// Attempt to remove the container
			err := cli.ContainerRemove(context.Background(), container.ID, types.ContainerRemoveOptions{Force: true})
			if err != nil {
				results = append(results, fmt.Sprintf("Failed to remove container: %s", container.ID))
			} else {
				results = append(results, fmt.Sprintf("Container deleted: %s", container.ID))
			}
		}
	}

	return results, nil
}

func GetContainerLogs(containerID string) (string, error) {
	// Create a new Docker client
	cli, err := client.NewClientWithOpts(client.WithVersion("1.41"))
	if err != nil {
		return "", err
	}

	// Check if the container exists
	_, err = cli.ContainerInspect(context.Background(), containerID)
	if err != nil {
		if client.IsErrNotFound(err) {
			return "", errors.New("Invalid Container ID")
		}
		return "", err
	}

	// Set up options for retrieving logs
	options := types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     false,
		Tail:       "all", // Retrieve all logs
	}

	// Get the logs
	logReader, err := cli.ContainerLogs(context.Background(), containerID, options)
	if err != nil {
		return "", err
	}
	defer logReader.Close()

	// Read the logs
	var logsBuilder strings.Builder
	if _, err := stdcopy.StdCopy(&logsBuilder, &logsBuilder, logReader); err != nil {
		return "", err
	}

	return logsBuilder.String(), nil
}

type ContainerStats struct {
	ID     string  `json:"id"`
	CPU    float64 `json:"cpu_usage_percent"`
	Memory float64 `json:"memory_usage_percent"`
}

// GetContainerStats retrieves statistics for a Docker container by its ID
func GetContainerStats(containerID string) (ContainerStats, error) {
	// Create a new Docker client
	cli, err := client.NewClientWithOpts(client.WithVersion("1.41"))
	if err != nil {
		return ContainerStats{}, err
	}

	// Check if the container exists
	_, err = cli.ContainerInspect(context.Background(), containerID)
	if err != nil {
		if client.IsErrNotFound(err) {
			return ContainerStats{}, errors.New("Invalid Container ID")
		}
		return ContainerStats{}, err
	}

	// Get the container stats
	stats, err := cli.ContainerStats(context.Background(), containerID, false)
	if err != nil {
		return ContainerStats{}, err
	}
	defer stats.Body.Close()

	// Decode the stats
	var stat types.Stats
	if err := json.NewDecoder(stats.Body).Decode(&stat); err != nil {
		return ContainerStats{}, err
	}

	// Calculate CPU usage percentage
	cpuUsage := calculateCPUPercentage(stat)

	// Calculate memory usage percentage
	memoryUsage := calculateMemoryUsagePercentage(stat)

	// Create a ContainerStats instance to return
	containerStats := ContainerStats{
		ID:     containerID,
		CPU:    cpuUsage,
		Memory: memoryUsage,
	}

	return containerStats, nil
}

// Helper function to calculate CPU usage percentage
func calculateCPUPercentage(stats types.Stats) float64 {
	// Calculate CPU usage based on the provided stats
	// Note: This is a basic example; implement the actual calculation based on your requirements
	return float64(stats.CPUStats.CPUUsage.TotalUsage) / float64(stats.CPUStats.SystemUsage) * 100
}

// Helper function to calculate memory usage percentage
func calculateMemoryUsagePercentage(stats types.Stats) float64 {
	// Calculate memory usage based on the provided stats
	return float64(stats.MemoryStats.Usage) / float64(stats.MemoryStats.Limit) * 100
}

func InspectContainer(containerID string) (types.ContainerJSON, error) {
	// Create a new Docker client
	cli, err := client.NewClientWithOpts(client.WithVersion("1.41"))
	if err != nil {
		return types.ContainerJSON{}, err
	}

	// Inspect the container
	containerJSON, err := cli.ContainerInspect(context.Background(), containerID)
	if err != nil {
		if client.IsErrNotFound(err) {
			return types.ContainerJSON{}, errors.New("Invalid Container ID")
		}
		return types.ContainerJSON{}, err
	}

	return containerJSON, nil
}
