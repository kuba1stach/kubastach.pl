package config

import "os"

type Config struct {
	CosmosConfig *CosmosConfig
}

type CosmosConfig struct {
	Endpoint      string
	Key           string
	DatabaseName  string
	ContainerName string
}

func LoadConfig() (*Config, error) {
	// For simplicity, using environment variables directly.
	// In a real application, consider using a library like Viper or similar.
	cosmosConfig := &CosmosConfig{
		Endpoint:      getEnv("COSMOS_ENDPOINT", "https://kubastachplproddb.documents.azure.com:443/"),
		Key:           getEnv("COSMOS_KEY", "=="),
		DatabaseName:  getEnv("COSMOS_DATABASE", "kubastachpl"),
		ContainerName: getEnv("COSMOS_CONTAINER", "data"),
	}

	return &Config{
		CosmosConfig: cosmosConfig,
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
