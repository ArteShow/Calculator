package jwt

import (
	"errors"
	"fmt"
	"time"

	config "github.com/ArteShow/Calculator/pkg/Config"
	database "github.com/ArteShow/Calculator/pkg/Database"
	"github.com/golang-jwt/jwt/v5"
)

func CreateJWT(userID int, role string, key string) (string, error) {
	now := time.Now()
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": userID,
        "iat":     now.Unix(),
		"nbf": now.Add(time.Minute * 1).Unix(),
		"exp": now.Add(time.Minute * 15).Unix(),
		"role":  role,
    })

    return token.SignedString([]byte(key))
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

func ValidateJWT(tokenString string, key string) (*jwt.Token, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is correct
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return key, nil
	})

	if err != nil {
		return nil, err
	}

	// Check if token is valid
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}