package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// JwtKey is the key used to sign the JWT tokens (can be set in environment variables for better security)
var JwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

// Claims represents the structure of the JWT token's claims
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// GenerateJWT generates a new JWT token with the given username and an expiration time of 24 hours
func GenerateJWT(username string) (string, error) {
	// Create the claims, which includes the username and expiration time
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // 24 hours expiration
		},
	}

	// Create a new JWT token with the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the JWT key
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}

	return tokenString, nil
}

// ValidateJWT validates the provided JWT token and returns the claims if valid
func ValidateJWT(tokenString string) (*Claims, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Return the JWT key to verify the token's signature
		return JwtKey, nil
	})

	// Check for parsing errors
	if err != nil {
		return nil, fmt.Errorf("invalid or expired token: %v", err)
	}

	// Extract the claims from the token
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
