package setup

import (
	"database/sql"
	"os"
	"testing"

	database "github.com/ArteShow/Calculator/pkg/database"
	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) (*sql.DB, string) {
	// Create a temporary database for testing
	dbPath := "./test.db"
	db, err := database.OpenDatabase(dbPath)
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	return db, dbPath
}

func TestCreateTable(t *testing.T) {
	db, _ := setupTestDB(t)
	defer db.Close()

	// Define columns for a test table
	columns := map[string]string{
		"id":       "INTEGER PRIMARY KEY AUTOINCREMENT",
		"username": "TEXT NOT NULL UNIQUE",
		"password": "TEXT NOT NULL",
	}

	err := CreateTable(db, "test_table", columns)
	assert.NoError(t, err)

	// Check if the table was created
	var exists int
	err = db.QueryRow("SELECT 1 FROM sqlite_master WHERE type='table' AND name='test_table'").Scan(&exists)
	assert.NoError(t, err)
	assert.Equal(t, 1, exists)
}

func TestInsertIfNotExists(t *testing.T) {
	db, _ := setupTestDB(t)
	defer db.Close()

	// Create the 'users' table
	columns := map[string]string{
		"id":       "INTEGER PRIMARY KEY AUTOINCREMENT",
		"username": "TEXT NOT NULL UNIQUE",
		"password": "TEXT NOT NULL",
	}
	err := CreateTable(db, "users", columns)
	assert.NoError(t, err)

	// Insert user data
	userData := map[string]interface{}{
		"username": "testuser",
		"password": "password123",
	}

	err = InsertIfNotExists(db, "users", userData, "username")
	assert.NoError(t, err)

	// Try inserting the same data again
	err = InsertIfNotExists(db, "users", userData, "username")
	assert.NoError(t, err)

	// Check that the data is inserted once
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", "testuser").Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestSetup(t *testing.T) {
	db, dbPath := setupTestDB(t)
	defer db.Close()

	// Run setup
	Setup()

	// Check that the JWT key was inserted
	var key string
	err := db.QueryRow("SELECT key FROM jwt WHERE id = 1").Scan(&key)
	assert.NoError(t, err)
	assert.NotEmpty(t, key)

	// Clean up by removing the test database file
	err = os.Remove(dbPath)
	assert.NoError(t, err)
}
