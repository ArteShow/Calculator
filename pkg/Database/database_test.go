package database

import (
	"database/sql"
	"path/filepath"
	"testing"
)

func setupTestDB(t *testing.T) (*sql.DB, string) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	db, err := CreateDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create DB: %v", err)
	}

	err = CreateTable(db, "users", map[string]string{
		"id":       "INTEGER PRIMARY KEY",
		"username": "TEXT UNIQUE",
		"password": "TEXT",
	})
	if err != nil {
		t.Fatalf("Failed to create users table: %v", err)
	}

	err = CreateTable(db, "calculations", map[string]string{
		"id":     "INTEGER PRIMARY KEY",
		"userId": "INTEGER",
	})
	if err != nil {
		t.Fatalf("Failed to create calculations table: %v", err)
	}

	return db, dbPath
}

func TestInsertAndGetUser(t *testing.T) {
	db, dbPath := setupTestDB(t)
	defer db.Close()

	user := map[string]interface{}{
		"id":       1,
		"username": "testuser",
		"password": "secret",
	}
	err := InsertData(db, "users", user)
	if err != nil {
		t.Fatalf("InsertData failed: %v", err)
	}

	gotUser, err := GetUserByUsername(dbPath, "testuser")
	if err != nil {
		t.Fatalf("GetUserByUsername failed: %v", err)
	}
	if gotUser.Username != "testuser" || gotUser.Password != "secret" {
		t.Fatalf("Unexpected user data: %+v", gotUser)
	}
}

func TestDeleteData(t *testing.T) {
	db, _ := setupTestDB(t)
	defer db.Close()

	_ = InsertData(db, "users", map[string]interface{}{
		"id":       2,
		"username": "deleteMe",
		"password": "pw",
	})
	err := DeleteData(db, "users", "username = 'deleteMe'")
	if err != nil {
		t.Fatalf("DeleteData failed: %v", err)
	}

	var count int
	_ = db.QueryRow("SELECT COUNT(*) FROM users WHERE username = 'deleteMe'").Scan(&count)
	if count != 0 {
		t.Fatalf("Expected 0 users, got %d", count)
	}
}

func TestUpdateData(t *testing.T) {
	db, _ := setupTestDB(t)
	defer db.Close()

	_ = InsertData(db, "users", map[string]interface{}{
		"id":       3,
		"username": "updateMe",
		"password": "old",
	})
	err := UpdateData(db, "users", map[string]interface{}{
		"password": "newpass",
	}, "username = 'updateMe'")
	if err != nil {
		t.Fatalf("UpdateData failed: %v", err)
	}

	var pw string
	_ = db.QueryRow("SELECT password FROM users WHERE username = 'updateMe'").Scan(&pw)
	if pw != "newpass" {
		t.Fatalf("Expected updated password, got %s", pw)
	}
}

func TestGetUserID(t *testing.T) {
	db, dbPath := setupTestDB(t)
	defer db.Close()

	_ = InsertData(db, "users", map[string]interface{}{
		"id":       5,
		"username": "getid",
		"password": "pw",
	})
	id, err := GetUserID(dbPath, "getid")
	if err != nil {
		t.Fatalf("GetUserID failed: %v", err)
	}
	if id != "5" {
		t.Fatalf("Expected user ID 5, got %s", id)
	}
}

func TestGetMaxId(t *testing.T) {
	db, dbPath := setupTestDB(t)
	defer db.Close()

	_ = InsertData(db, "users", map[string]interface{}{
		"id":       7,
		"username": "maxid",
		"password": "pw",
	})
	maxID, err := GetMaxId(dbPath, "users")
	if err != nil {
		t.Fatalf("GetMaxId failed: %v", err)
	}
	if maxID != 7 {
		t.Fatalf("Expected max ID 7, got %d", maxID)
	}
}

func TestGetMaxExpressionIdByUserId(t *testing.T) {
	db, _ := setupTestDB(t)
	defer db.Close()

	_ = InsertData(db, "calculations", map[string]interface{}{
		"id":     10,
		"userId": 99,
	})
	maxID, err := GetMaxExpressionIdByUserId(db, 99)
	if err != nil {
		t.Fatalf("GetMaxExpressionIdByUserId failed: %v", err)
	}
	if maxID != 10 {
		t.Fatalf("Expected max expression ID 10, got %d", maxID)
	}
}
