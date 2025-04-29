package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func createJWT(userID int, role string, key string) (string, error) {
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

