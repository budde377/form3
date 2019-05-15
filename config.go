package main

import (
	"encoding/json"
	"os"

	"github.com/google/logger"
)

type Config struct {
	Port int    `json:"port"`
	Host string `json:"host"`
}

var defaultConfig = Config{
	Port: 8080,
	Host: "http://example.com",
}

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

func OpenConfig() *Config {
	path := "./config/config.json"
	return OpenConfigFile(path)
}
