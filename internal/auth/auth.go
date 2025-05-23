package auth

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/xgodev/boost/wrapper/log"
)

// User represents an authenticated user
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

// AuthService handles authentication operations
type AuthService struct {
	// Add dependencies here
	jwtSecret []byte
}

// NewAuthService creates a new authentication service
func NewAuthService(jwtSecret string) *AuthService {
	return &AuthService{
		jwtSecret: []byte(jwtSecret),
	}
}

// Login authenticates a user and returns a JWT token
func (s *AuthService) Login(ctx context.Context, username, password string) (string, error) {
	log.Infof("Authenticating user: %s", username)
	
	// Mock implementation for now
	// In a real implementation, this would validate credentials against a database
	if username == "admin" && password == "admin" {
		// Create a new token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":       "1",
			"username": username,
			"email":    "admin@example.com",
			"role":     "admin",
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})
		
		// Sign and get the complete encoded token as a string
		tokenString, err := token.SignedString(s.jwtSecret)
		if err != nil {
			return "", err
		}
		
		return tokenString, nil
	}
	
	return "", errors.New("invalid credentials")
}

// ValidateToken validates a JWT token and returns the user
func (s *AuthService) ValidateToken(ctx context.Context, tokenString string) (*User, error) {
	log.Info("Validating token")
	
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		
		return s.jwtSecret, nil
	})
	
	if err != nil {
		return nil, err
	}
	
	// Validate the token
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		user := &User{
			ID:       claims["id"].(string),
			Username: claims["username"].(string),
			Email:    claims["email"].(string),
			Role:     claims["role"].(string),
		}
		
		return user, nil
	}
	
	return nil, errors.New("invalid token")
}
