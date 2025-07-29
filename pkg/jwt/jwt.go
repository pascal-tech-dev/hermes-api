package jwt

import (
	"errors"
	"fmt"
	"time"

	"hermes-api/internal/model"

	jwtlib "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Claims represents the JWT claims structure
type Claims struct {
	UserID   uuid.UUID `json:"user_id"`
	Email    string    `json:"email"`
	Username string    `json:"username"`
	jwtlib.RegisteredClaims
}

// GenerateToken creates a JWT token for a user
func GenerateToken(user *model.User, secret string, expiration time.Duration) (string, error) {
	claims := Claims{
		UserID:   user.ID,
		Email:    user.Email,
		Username: user.Username,
		RegisteredClaims: jwtlib.RegisteredClaims{
			ExpiresAt: jwtlib.NewNumericDate(time.Now().Add(expiration)),
			IssuedAt:  jwtlib.NewNumericDate(time.Now()),
			NotBefore: jwtlib.NewNumericDate(time.Now()),
			Issuer:    "hermes-api",
			Subject:   user.ID.String(),
		},
	}

	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseToken parses and validates a JWT token
func ParseToken(tokenString, secret string) (*jwtlib.Token, error) {
	token, err := jwtlib.ParseWithClaims(tokenString, &Claims{}, func(token *jwtlib.Token) (any, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwtlib.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}

// ExtractClaims extracts claims from a JWT token
func ExtractClaims(token *jwtlib.Token) (*Claims, error) {
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

// IsTokenExpired checks if a token is expired
func IsTokenExpired(err error) bool {
	return errors.Is(err, jwtlib.ErrTokenExpired)
}

// GetUserIDFromToken extracts user ID from a token string
func GetUserIDFromToken(tokenString, secret string) (uuid.UUID, error) {
	claims, err := ValidateToken(tokenString, secret)
	if err != nil {
		return uuid.UUID{}, err
	}

	return claims.UserID, nil
}

// ValidateToken validates a token and returns user information
func ValidateToken(tokenString, secret string) (*Claims, error) {
	token, err := ParseToken(tokenString, secret)
	if err != nil {
		return nil, err
	}

	claims, err := ExtractClaims(token)
	if err != nil {
		return nil, err
	}

	return claims, nil
}
