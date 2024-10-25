package api

import (
	"Docker_Management/pkg/docker" // replace with actual path to docker package
	"encoding/json"
	"net/http"
)

// ListVolumesHandler is an HTTP handler to list all Docker volumes
func ListVolumesHandler(w http.ResponseWriter, r *http.Request) {
	volumes, err := docker.ListVolumes()
	if err != nil {
		http.Error(w, "Failed to list volumes", http.StatusInternalServerError)
		return
	}

	if len(volumes) == 0 {
		// If no volumes are present, return an appropriate message
		response := map[string]string{"message": "No Volumes are present"}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Format volumes as a JSON response
	response := []map[string]string{}
	for _, volume := range volumes {
		response = append(response, map[string]string{
			"Name":       volume.Name,
			"Driver":     volume.Driver,
			"Mountpoint": volume.Mountpoint,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
    w.Header().Set("Access-Control-Allow-Credentials", "true")
	json.NewEncoder(w).Encode(response)
}

func InspectVolumeHandler(w http.ResponseWriter, r *http.Request) {
	// Decode the request body
	var reqBody RequestBodyVolume
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Call the function to inspect the volume
	volume, err := docker.InspectVolume(reqBody.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
    w.Header().Set("Access-Control-Allow-Credentials", "true")
	if err := json.NewEncoder(w).Encode(volume); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

type RequestBodyVolume struct {
	Name string `json:"name"` // Ensure the field is defined as "Name" (capitalized to export it)
}

type ResponseVolumeContainers struct {
	Message      string   `json:"message"`
	ContainerIDs []string `json:"container_ids,omitempty"`
}

// ListContainersAttachedToVolumeHandler handles requests to list containers attached to a volume
func ListContainersAttachedToVolumeHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body to get the volume name
	var reqBody RequestBodyVolume
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Call the function to list containers attached to the specified volume
	containerIDs, err := docker.ListContainersAttachedToVolume(reqBody.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create the response struct
	response := ResponseVolumeContainers{
		Message:      "Containers attached to the volume:",
		ContainerIDs: containerIDs,
	}

	// Return the response in JSON format
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
    w.Header().Set("Access-Control-Allow-Credentials", "true")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func RemoveVolumeHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody RequestBodyVolume
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	message, err := docker.RemoveVolume(reqBody.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := map[string]string{"message": message}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
    w.Header().Set("Access-Control-Allow-Credentials", "true")
	json.NewEncoder(w).Encode(response)
}