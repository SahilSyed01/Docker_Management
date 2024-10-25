package main

import (
	"Docker_Management/pkg/api"
	"Docker_Management/pkg/config"
	"log"
	"net/http"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Connect to MongoDB using the loaded MongoURI
	// db.ConnectDB(config.AppConfig.MongoURI)

	// Set up routes
	log.Printf("Starting server on :%s", config.AppConfig.ServerPort)
	router := api.SetupRouter()

	// Start the server
	log.Printf("Started Server on :%s",config.AppConfig.ServerPort)
	if err := http.ListenAndServe(":"+config.AppConfig.ServerPort, router); err != nil {
		log.Fatal(err)
	}
}
