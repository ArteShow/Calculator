package config

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupTestConfig(t *testing.T) {
	err := os.MkdirAll("configs", os.ModePerm)
	assert.NoError(t, err)

	dbConfig := DatabaseConfig{
		Driver: "sqlite3",
		Path:   "test.db",
	}
	file, err := os.Create("configs/database.json")
	assert.NoError(t, err)
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(dbConfig)
	assert.NoError(t, err)
}

func teardownTestConfig() {
	_ = os.Remove("configs/database.json")
	_ = os.Remove("configs")
}

func TestLoadDatabaseConfig_Success(t *testing.T) {
	setupTestConfig(t)
	defer teardownTestConfig()

	cfg, err := LoadDatabaseConfig()
	assert.NoError(t, err)
	assert.Equal(t, "sqlite3", cfg.Driver)
	assert.Equal(t, "test.db", cfg.Path)
}

func TestLoadDatabaseConfig_FileNotFound(t *testing.T) {
	teardownTestConfig() // ensure file doesn't exist

	cfg, err := LoadDatabaseConfig()
	assert.Nil(t, cfg)
	assert.Error(t, err)
}

func TestGetDatabasePath(t *testing.T) {
	setupTestConfig(t)
	defer teardownTestConfig()

	path := GetDatabasePath()
	assert.Equal(t, "test.db", path)
}
