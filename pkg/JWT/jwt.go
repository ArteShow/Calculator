package jwt

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	config "github.com/ArteShow/Calculator/pkg/Config"
	database "github.com/ArteShow/Calculator/pkg/Database"
	"github.com/golang-jwt/jwt/v5"
	_ "modernc.org/sqlite"
)

func CreateJWT(userID int, role string, key string) (string, error) {
	now := time.Now()
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": userID,
        "iat":     now.Unix(),
		"nbf": now.Add(time.Second * 1).Unix(),
		"exp": now.Add(time.Minute * 15).Unix(),
		"role":  role,
    })

    return token.SignedString([]byte(key))
}

func InsertJWTKey(db *sql.DB, jwtKey string) error {
	query := "INSERT INTO jwt (key) VALUES (?)"
	_, err := db.Exec(query, jwtKey)
	return err
}

func ParseJWT(tokenString string, secret []byte) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
}

func GenerateJWTKey(length int) (string, error) {
	// Create a byte slice to hold the random key
	key := make([]byte, length)

	// Read random bytes from the crypto/rand package
	_, err := rand.Read(key)
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT key: %w", err)
	}

	// Encode the random bytes into a base64 string
	return base64.URLEncoding.EncodeToString(key), nil
}

func GetJWTKey() string {
	query := "SELECT key FROM jwt WHERE id = 1;"
	db, err := database.OpenDatabase(config.GetDatabasePath())
	if err != nil {
		panic(err)
	}
	defer db.Close()
	var key string
	err = db.QueryRow(query).Scan(&key)
	if err != nil {
		panic(err)
	}
	return key
}

func ExtractUserIDFromToken(tokenString string, secret []byte) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return secret, nil
	})
	if err != nil || !token.Valid {
		return "", fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("cannot parse claims")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return "", fmt.Errorf("user_id not found in token")
	}

	return fmt.Sprintf("%.0f", userID), nil
}

func ValidateJWT(tokenString string, key string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(key), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}