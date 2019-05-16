package main

import (
	"os"
)

// Config Shared configuration format
type Config struct {
	Port            int    `json:"port"`
	Host            string `json:"host"`
	MongoDbURI      string `json:"mongo_db_uri"`
	MongoDbDatabase string `json:"mongo_db_database"`
}

// ReadConfigFromEnv Open a configuration at the given path.
func ReadConfigFromEnv() *Config {
	var c = Config{
		Port:            SafeStringToInt(os.Getenv("PORT"), 8080),
		Host:            os.Getenv("HOST"),
		MongoDbDatabase: os.Getenv("MONGO_DB_DATABASE"),
		MongoDbURI:      os.Getenv("MONGO_DB_URI"),
	}
	return &c
}
