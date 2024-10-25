package api

import (
	"encoding/json"
	"net/http"
	"Docker_Management/pkg/docker" // replace with the actual path to your docker package
)

func ListNetworksHandler(w http.ResponseWriter, r *http.Request) {
	// Get the list of networks
	networks, err := docker.ListNetworks()
	if err != nil {
		http.Error(w, "Failed to retrieve networks", http.StatusInternalServerError)
		return
	}

	// Check if no networks were found
	if len(networks) == 0 {
		http.Error(w, "No networks found", http.StatusNotFound)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(networks)
}

// InspectNetworkHandler handles the request to inspect a Docker network by its ID
func InspectNetworkHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody RequestBody

	// Decode the request body to get the network ID
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid input format", http.StatusBadRequest)
		return
	}

	// Check if the ID field is empty
	if reqBody.ID == "" {
		http.Error(w, "Network ID is required", http.StatusBadRequest)
		return
	}

	// Inspect the network
	networkDetails, err := docker.InspectNetwork(reqBody.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Send response in JSON format
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(networkDetails)
}

type NetworkRequestBody struct {
    NetworkID string `json:"id"`
}

func ListContainersInNetworkHandler(w http.ResponseWriter, r *http.Request) {
    var reqBody NetworkRequestBody

    // Parse the request body
    if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    // Call the Docker function to get containers attached to the network
    containers, err := docker.ListContainersInNetwork(reqBody.NetworkID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    // Encode the response as JSON and send it back to the client
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "containers_attached_to_network": containers,
    })
}

type NetworkRemoveRequestBody struct {
    NetworkID string `json:"id"`
}

func RemoveNetworkHandler(w http.ResponseWriter, r *http.Request) {
    var reqBody NetworkRemoveRequestBody

    // Parse the request body
    if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    // Call the Docker function to remove the network
    message, err := docker.RemoveNetwork(reqBody.NetworkID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    // Send success message in response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "message": message,
    })
}