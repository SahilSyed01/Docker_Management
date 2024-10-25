package api

import (
	"encoding/json"
	"net/http"

	"Docker_Management/pkg/docker"

	"github.com/docker/docker/api/types"
)

// ContainerResponse represents the structure of the container information
type ContainerResponse struct {
	ID     string `json:"id"`
	Image  string `json:"image"`
	Status string `json:"status"`
}

func ListContainersHandler(w http.ResponseWriter, r *http.Request) {
	containers, err := docker.ListContainers()
	if err != nil {
		http.Error(w, "Failed to list containers", http.StatusInternalServerError)
		return
	}

	// Prepare the response slice
	var response []ContainerResponse
	for _, container := range containers {
		response = append(response, ContainerResponse{
			ID:     container.ID,
			Image:  container.Image,
			Status: container.Status,
		})
	}

	// Set the content type to application/json
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
    w.Header().Set("Access-Control-Allow-Credentials", "true")

	// Marshal the response to JSON
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func ListAllContainersHandler(w http.ResponseWriter, r *http.Request) {
	containers, err := docker.ListAllContainers()
	if err != nil {
		http.Error(w, "Failed to list containers", http.StatusInternalServerError)
		return
	}

	// Prepare the response slice
	var response []ContainerResponse
	for _, container := range containers {
		response = append(response, ContainerResponse{
			ID:     container.ID,
			Image:  container.Image,
			Status: container.Status,
		})
	}

	// Set the content type to application/json
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
    w.Header().Set("Access-Control-Allow-Credentials", "true")

	// Marshal the response to JSON
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

type RequestBody struct {
	ID string `json:"id"` // JSON field to hold the container ID
}

// StartContainerHandler handles the HTTP request to start a container
func StartContainerHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody RequestBody

	// Decode the JSON request body
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
    
	// Call the StartContainer function with the provided container ID
	message, err := docker.StartContainer(requestBody.ID)
	if err != nil {
		http.Error(w, "Failed to start container: "+err.Error(), http.StatusInternalServerError)
		return
	}
    w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
    w.Header().Set("Access-Control-Allow-Credentials", "true")
	// Respond with the message
	w.Write([]byte(message))
}

func StopContainerHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody RequestBody

	// Decode the JSON request body
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Call the StopContainer function with the provided container ID
	message, err := docker.StopContainer(requestBody.ID)
	if err != nil {
		http.Error(w, "Failed to stop container: "+err.Error(), http.StatusInternalServerError)
		return
	}
    w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
    w.Header().Set("Access-Control-Allow-Credentials", "true")
	// Respond with the message
	w.Write([]byte(message))
}

func RemoveContainerHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody RequestBody

	// Decode the JSON request body
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Call the RemoveContainer function with the provided container ID
	message, err := docker.RemoveContainer(requestBody.ID)
	if err != nil {
		http.Error(w, "Failed to remove container: "+err.Error(), http.StatusInternalServerError)
		return
	}
    w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
    w.Header().Set("Access-Control-Allow-Credentials", "true")
	// Respond with the message
	w.Write([]byte(message))
}

type LogResponse struct {
	ID   string `json:"id"`
	Logs string `json:"logs"`
}
func RemoveAllContainersHandler(w http.ResponseWriter, r *http.Request) {
    // Call the RemoveAllContainers function
    results, err := docker.RemoveAllContainers()
    if err != nil {
        http.Error(w, "Failed to remove containers: "+err.Error(), http.StatusInternalServerError)
        return
    }
    
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
    w.Header().Set("Access-Control-Allow-Credentials", "true")
    // Set the response content type to JSON
    w.Header().Set("Content-Type", "application/json")

    // Return the result messages in JSON format
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": "Containers Deleted",
        "details": results,
    })
}
// GetContainerLogsHandler handles the HTTP request to get logs for a container
func GetContainerLogsHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody RequestBody

	// Decode the JSON request body
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Call the GetContainerLogs function with the provided container ID
	logs, err := docker.GetContainerLogs(requestBody.ID)
	if err != nil {
		http.Error(w, "Failed to retrieve logs: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare the log response
	logResponse := LogResponse{
		ID:   requestBody.ID,
		Logs: logs,
	}
    
	// Set the response content type to JSON
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
    w.Header().Set("Access-Control-Allow-Credentials", "true")
	json.NewEncoder(w).Encode(logResponse)
}

type StatsResponse struct {
	ID     string  `json:"id"`
	CPU    float64 `json:"cpu_usage_percent"`
	Memory float64 `json:"memory_usage_percent"`
}

// GetContainerStatsHandler handles the HTTP request to get stats for a container
func GetContainerStatsHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody RequestBody

	// Decode the JSON request body
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Call the GetContainerStats function with the provided container ID
	stats, err := docker.GetContainerStats(requestBody.ID)
	if err != nil {
		http.Error(w, "Failed to retrieve stats: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare the stats response
	statsResponse := StatsResponse{
		ID:     stats.ID,
		CPU:    stats.CPU,
		Memory: stats.Memory,
	}

	// Set the response content type to JSON
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
    w.Header().Set("Access-Control-Allow-Credentials", "true")
	json.NewEncoder(w).Encode(statsResponse)
}

type InspectResponse struct {
	ContainerInfo types.ContainerJSON `json:"container_info"` // Use types.ContainerJSON directly
}

// InspectContainerHandler handles the HTTP request to inspect a container
func InspectContainerHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody RequestBody

	// Decode the JSON request body
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Call the InspectContainer function with the provided container ID
	containerInfo, err := docker.InspectContainer(requestBody.ID)
	if err != nil {
		http.Error(w, "Failed to inspect container: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare the response
	inspectResponse := InspectResponse{
		ContainerInfo: containerInfo,
	}

	// Set the response content type to JSON
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
    w.Header().Set("Access-Control-Allow-Credentials", "true")
	json.NewEncoder(w).Encode(inspectResponse)
}

