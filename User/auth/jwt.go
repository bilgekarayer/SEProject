package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

type Claims struct {
	UserID    int    `json:"uid"`
	Username  string `json:"uname"`
	Role      string `json:"role"`
	FirstName string `json:"fname"`
	LastName  string `json:"lname"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID int, username, role, firstName, lastName string) (string, error) {
	claims := Claims{
		UserID:    userID,
		Username:  username,
		Role:      role,
		FirstName: firstName,
		LastName:  lastName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString([]byte(os.Getenv("JWT_SECRET")))
}
