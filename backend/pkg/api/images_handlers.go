package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"Docker_Management/pkg/docker"
)

type ImageResponse struct {
	ID      string `json:"id"`
	Name    string `json:"name"` // Repository name
	Tag     string `json:"tag"`  // Image tag
	Created string `json:"created"`
	Size    string `json:"size"` // Size in MB
}

func ListImagesHandler(w http.ResponseWriter, r *http.Request) {
	// Call the ListImages function
	images, err := docker.ListImages()
	if err != nil {
		http.Error(w, "Failed to list images: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare the response
	var response []ImageResponse
	for _, img := range images {
		// Initialize name and tag
		var name, tag string

		// Check if there are repository tags
		if len(img.RepoTags) > 0 {
			fullTag := img.RepoTags[0]           // Use the first tag as the full tag
			parts := strings.Split(fullTag, ":") // Split by colon

			// Assign name and tag based on the split
			name = parts[0] // Repository name
			if len(parts) > 1 {
				tag = parts[1] // Image tag, if available
			} else {
				tag = "latest" // Default to "latest" if no tag is specified
			}
		}

		// Convert size from bytes to MB
		sizeInMB := float64(img.Size) / (1024 * 1024) // Convert bytes to MB
		sizeFormatted := formatSize(sizeInMB)

		// Create the ImageResponse
		response = append(response, ImageResponse{
			ID:      img.ID,
			Name:    name,
			Tag:     tag,
			Created: time.Unix(img.Created, 0).Format(time.RFC3339), // Format the created time
			Size:    sizeFormatted,                                  // Use formatted size
		})
	}

	// Set the response content type to JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Helper function to format size
func formatSize(size float64) string {
	return fmt.Sprintf("%.0f MB", size) // Format the size to 0 decimal places
}

type DanglingImageResponse struct {
	ID      string `json:"id"`
	Created string `json:"created"`
	Size    string `json:"size"` // Size in MB
}

// ListDanglingImagesHandler handles the HTTP request to list all dangling images
func ListDanglingImagesHandler(w http.ResponseWriter, r *http.Request) {
	// Call the ListDanglingImages function
	images, err := docker.ListDanglingImages()
	if err != nil {
		http.Error(w, "Failed to list dangling images: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if len(images) == 0 {
		response := map[string]string{"message": "No Dangling images"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}
	// Prepare the response
	var response []DanglingImageResponse
	for _, img := range images {
		// Convert size from bytes to MB
		sizeInMB := float64(img.Size) / (1024 * 1024) // Convert bytes to MB
		sizeFormatted := fmt.Sprintf("%.0f MB", sizeInMB)

		// Create the DanglingImageResponse
		response = append(response, DanglingImageResponse{
			ID:      img.ID,
			Created: time.Unix(img.Created, 0).Format(time.RFC3339), // Format the created time
			Size:    sizeFormatted,                                  // Use formatted size
		})
	}

	// Set the response content type to JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

type RemoveImageRequest struct {
	ID string `json:"id"` // ID of the image to remove
}

// RemoveImageHandler handles the HTTP request to remove a Docker image
func RemoveImageHandler(w http.ResponseWriter, r *http.Request) {
	var req RemoveImageRequest

	// Decode the request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Call the RemoveImage function
	message, err := docker.RemoveImage(req.ID)
	if err != nil {
		http.Error(w, "Failed to remove image: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the response content type to JSON and return success message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}

func RemoveAllImagesHandler(w http.ResponseWriter, r *http.Request) {
	// Call the RemoveAllImages function
	results, err := docker.RemoveAllImages()
	if err != nil {
		http.Error(w, "Failed to remove images: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the response content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Return the result messages in JSON format
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Images Deleted",
		"details": results,
	})
}

func InspectImageHandler(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Parse the JSON body
	var requestData struct {
		ID string `json:"id"`
	}

	if err := json.Unmarshal(body, &requestData); err != nil {
		http.Error(w, "Failed to parse JSON body", http.StatusBadRequest)
		return
	}

	// Call the InspectImage function
	imageDetails, err := docker.InspectImage(requestData.ID)
	if err != nil {
		http.Error(w, "Failed to inspect image: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the response content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Return the inspected image details in JSON format
	json.NewEncoder(w).Encode(imageDetails)
}

func PullImageHandler(w http.ResponseWriter, r *http.Request) {
    // Read the request body
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Parse the JSON input
    var requestData struct {
        Image string `json:"image"`
    }
    err = json.Unmarshal(body, &requestData)
    if err != nil || requestData.Image == "" {
        http.Error(w, "Invalid input format. Expected {image: \"image_name:version\"}", http.StatusBadRequest)
        return
    }

    // Pull the image
    result, err := docker.PullImage(requestData.Image)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Send the success response
    response := map[string]string{"message": result}
    json.NewEncoder(w).Encode(response)
}