package config

import (
	"os"
)

type Config struct {
	MongoURI   string
	ServerPort string
}

var AppConfig Config

func LoadConfig() {
	// Load environment variables from .env file
	// err := godotenv.Load()
	// if err != nil {
	//     log.Println("No .env file found, using default values")
	// }

	AppConfig = Config{
		//MongoURI:   getEnv("MONGO_URI", "mongodb://localhost:27017"),
		ServerPort: getEnv("SERVER_PORT", "8090"),
	}
}

// getEnv retrieves the value of the environment variable or returns a fallback value if not set.
func getEnv(key string, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}
