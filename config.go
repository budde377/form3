package main

import (
	"encoding/json"
	"os"

	"github.com/google/logger"
)

// Config Shared configration format
type Config struct {
	Port            int    `json:"port"`
	Host            string `json:"host"`
	MongoDbURI      string `json:"mongo_db_uri"`
	MongoDbDatabase string `json:"mongo_db_database"`
}

var defaultConfig = Config{
	Port: 8080,
	Host: "http://example.com",
}

// OpenConfigFile Open a configuration at the given path.
func OpenConfigFile(path string) *Config {
	file, err := os.Open(path)
	c := defaultConfig
	if err != nil {
		logger.Warning("Failed to open config file: ", err)
	}
	if file != nil {
		defer file.Close()
	} else {
		return &defaultConfig
	}
	if err := json.NewDecoder(file).Decode(&c); err != nil {
		logger.Fatal("Failed to decode configuration file: ", err)
	}
	return &c
}
