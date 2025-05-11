package jwt

import (
	"database/sql"
	"testing"

	_ "github.com/ArteShow/Calculator/pkg/database"
	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) (*sql.DB, string) {
	db, dbPath := setupTestDB(t)
	return db, dbPath
}

func TestCreateJWT(t *testing.T) {
	userID := 1
	role := "admin"
	key := "secretkey"

	token, err := CreateJWT(userID, role, key)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestGenerateJWTKey(t *testing.T) {
	key, err := GenerateJWTKey(32)
	assert.NoError(t, err)
	assert.Len(t, key, 44) // Base64 URL encoding will make it 44 chars for 32 byte input
}

func TestInsertJWTKey(t *testing.T) {
	db, _ := setupTestDB(t)
	defer db.Close()

	key := "jwt-secret-key"
	err := InsertJWTKey(db, key)
	assert.NoError(t, err)

	var insertedKey string
	err = db.QueryRow("SELECT key FROM jwt WHERE id = 1").Scan(&insertedKey)
	assert.NoError(t, err)
	assert.Equal(t, key, insertedKey)
}

func TestGetJWTKey(t *testing.T) {
	db, _ := setupTestDB(t)
	defer db.Close()

	key := "jwt-secret-key"
	err := InsertJWTKey(db, key)
	assert.NoError(t, err)

	retrievedKey := GetJWTKey()
	assert.Equal(t, key, retrievedKey)
}

func TestParseJWT(t *testing.T) {
	userID := 1
	role := "admin"
	key := "secretkey"

	token, err := CreateJWT(userID, role, key)
	assert.NoError(t, err)

	parsedToken, err := ParseJWT(token, []byte(key))
	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)
}

func TestValidateJWT(t *testing.T) {
	userID := 1
	role := "admin"
	key := "secretkey"

	token, err := CreateJWT(userID, role, key)
	assert.NoError(t, err)

	validToken, err := ValidateJWT(token, key)
	assert.NoError(t, err)
	assert.True(t, validToken.Valid)

	invalidToken := "invalidToken"
	_, err = ValidateJWT(invalidToken, key)
	assert.Error(t, err)
}

func TestExtractUserIDFromToken(t *testing.T) {
	userID := 1
	role := "admin"
	key := "secretkey"

	token, err := CreateJWT(userID, role, key)
	assert.NoError(t, err)

	extractedUserID, err := ExtractUserIDFromToken(token, []byte(key))
	assert.NoError(t, err)
	assert.Equal(t, "1", extractedUserID)
}
