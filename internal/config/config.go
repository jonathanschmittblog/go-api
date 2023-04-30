package config

import (
	"encoding/json"
	"os"
)

type ServiceConfiguration struct {
	Server   ServerConfiguration   `json:"server"`
	Database DatabaseConfiguration `json:"database"`
}

type ServerConfiguration struct {
	Port     int    `json:"port"`
	PageSize int    `json:"page_size"`
	Version  string `json:"version"`
}

type DatabaseConfiguration struct {
	DatabaseMinutesIdle  int    `json:"minutes_idle"`
	DatabaseMaxIdleConns int    `json:"max_idle_conns"`
	DatabaseMaxOpenConns int    `json:"max_open_conns"`
	ConnectionString     string `json:"connection_string"`
}

func GetConfig() (ServiceConfiguration, error) {
	config := ServiceConfiguration{}
	configValue, err := os.ReadFile("./config/config.json")
	if err != nil {
		return config, err
	}

	secretValue, err := os.ReadFile("./config/secret.json")
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(configValue, &config)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(secretValue, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}
