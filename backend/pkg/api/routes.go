package api

import (
	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/containers", ListContainersHandler).Methods("GET")
	router.HandleFunc("/containers/all", ListAllContainersHandler).Methods("GET")
	router.HandleFunc("/containers/start", StartContainerHandler).Methods("POST")
	router.HandleFunc("/containers/stop", StopContainerHandler).Methods("POST")
	router.HandleFunc("/containers/remove", RemoveContainerHandler).Methods("DELETE")
	router.HandleFunc("/containers/logs", GetContainerLogsHandler).Methods("POST")
	router.HandleFunc("/containers/stats", GetContainerStatsHandler).Methods("POST")
	router.HandleFunc("/containers/inspect", InspectContainerHandler).Methods("POST")
	router.HandleFunc("/containers/remove/all", RemoveAllContainersHandler).Methods("DELETE")

	router.HandleFunc("/images", ListImagesHandler).Methods("GET")
	router.HandleFunc("/images/dangling", ListDanglingImagesHandler).Methods("GET")
	router.HandleFunc("/images/remove", RemoveImageHandler).Methods("DELETE")							
	router.HandleFunc("/images/remove/all", RemoveAllImagesHandler).Methods("DELETE")
	router.HandleFunc("/images/dangling/remove/all", RemoveAllDanglingImagesHandler).Methods("DELETE")
	router.HandleFunc("/images/inspect", InspectImageHandler).Methods("POST")
	router.HandleFunc("/images/pull", PullImageHandler).Methods("POST")


	router.HandleFunc("/volumes", ListVolumesHandler).Methods("GET")
	router.HandleFunc("/volumes/inspect", InspectVolumeHandler).Methods("POST")
	router.HandleFunc("/volumes/containers", ListContainersAttachedToVolumeHandler).Methods("POST")
	router.HandleFunc("/volumes/remove", RemoveVolumeHandler).Methods("DELETE")


	router.HandleFunc("/networks", ListNetworksHandler).Methods("GET")
	router.HandleFunc("/networks/inspect", InspectNetworkHandler).Methods("POST")
	router.HandleFunc("/networks/containers", ListContainersInNetworkHandler).Methods("POST")
	router.HandleFunc("/networks/remove", RemoveNetworkHandler).Methods("DELETE")
	
	return router
}
