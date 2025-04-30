package config

import (
	"encoding/json"
	"os"
)

type DatabaseConfig struct {
	Driver   string `json:"driver"`
	Path    string `json:"path"`
}

func LoadDatabaseConfig() (*DatabaseConfig, error) {
	DatabseConfig := &DatabaseConfig{}
	file, err := os.Open("configs/database.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(DatabseConfig)
	if err != nil {
		return nil, err
	}
	return DatabseConfig, nil
}

func GetDatabasePath() string {
	config, err := LoadDatabaseConfig()
	if err != nil {
		panic(err)
	}
	return config.Path
}

