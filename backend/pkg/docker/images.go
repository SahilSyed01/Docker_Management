package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

func ListImages() ([]types.ImageSummary, error) {
	// Create a new Docker client
	cli, err := client.NewClientWithOpts(client.WithVersion("1.41"))
	if err != nil {
		return nil, err
	}

	// List images
	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		return nil, err
	}

	return images, nil
}

func ListDanglingImages() ([]types.ImageSummary, error) {
	// Create a new Docker client
	cli, err := client.NewClientWithOpts(client.WithVersion("1.41"))
	if err != nil {
		return nil, err
	}

	// Get all images
	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		return nil, err
	}

	// Filter for dangling images
	var danglingImages []types.ImageSummary
	for _, img := range images {
		if len(img.RepoTags) == 0 { // No tags means it's a dangling image
			danglingImages = append(danglingImages, img)
		}
	}

	return danglingImages, nil
}

func RemoveImage(imageID string) (string, error) {
	// Create a new Docker client
	cli, err := client.NewClientWithOpts(client.WithVersion("1.41"))
	if err != nil {
		return "", err
	}

	// Remove the image
	_, err = cli.ImageRemove(context.Background(), imageID, types.ImageRemoveOptions{Force: true})
	if err != nil {
		return "", err
	}

	return "Image removed successfully", nil
}

func RemoveAllImages() ([]string, error) {
	// Create a new Docker client
	cli, err := client.NewClientWithOpts(client.WithVersion("1.41"))
	if err != nil {
		return nil, err
	}

	// Get the list of all images
	images, err := cli.ImageList(context.Background(), types.ImageListOptions{All: true})
	if err != nil {
		return nil, err
	}

	// Slice to hold the result messages
	var results []string

	// Iterate over each image and try to remove it
	for _, image := range images {
		_, err := cli.ImageRemove(context.Background(), image.ID, types.ImageRemoveOptions{Force: true})
		if err != nil {
			// Append error message if the image is being used
			results = append(results, fmt.Sprintf("Cannot remove, image is being used: %s", image.ID))
		} else {
			// Append success message
			results = append(results, fmt.Sprintf("Image deleted: %s", image.ID))
		}
	}

	return results, nil
}

func RemoveAllDanglingImages() ([]string, error) {
	// Create a new Docker client
	cli, err := client.NewClientWithOpts(client.WithVersion("1.41"))
	if err != nil {
		return nil, err
	}

	// Set up filter to list only dangling images
	imageFilter := filters.NewArgs()
	imageFilter.Add("dangling", "true")

	// Get the list of dangling images
	images, err := cli.ImageList(context.Background(), types.ImageListOptions{Filters: imageFilter})
	if err != nil {
		return nil, err
	}

	// Check if there are no dangling images
	if len(images) == 0 {
		return []string{"No dangling images found"}, nil
	}

	// Slice to hold the result messages
	var results []string

	// Iterate over each image and try to remove it
	for _, image := range images {
		// Attempt to remove the image
		_, err := cli.ImageRemove(context.Background(), image.ID, types.ImageRemoveOptions{Force: true})
		if err != nil {
			results = append(results, fmt.Sprintf("Failed to remove image: %s", image.ID))
		} else {
			results = append(results, fmt.Sprintf("Image deleted: %s", image.ID))
		}
	}

	return results, nil
}
func InspectImage(imageID string) (types.ImageInspect, error) {
    // Create a new Docker client
    cli, err := client.NewClientWithOpts(client.WithVersion("1.41"))
    if err != nil {
        return types.ImageInspect{}, err
    }

    // Inspect the image
    imageInspect, _, err := cli.ImageInspectWithRaw(context.Background(), imageID)
    if err != nil {
        return types.ImageInspect{}, err
    }

    return imageInspect, nil
}