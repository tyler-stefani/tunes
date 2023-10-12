package auth

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	ACCESS_TOKEN_DURATION  = time.Minute * 15
	REFRESH_TOKEN_DURATION = time.Hour * 24
)

func HashPassword(unhashed string) (string, error) {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(unhashed), 15)
	return string(hashed), nil
}

func ValidatePassword(unhashed, hashed string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(unhashed))
	if err != nil {
		return false, err
	}
	return true, nil
}

func CreateAccessToken() (string, error) {
	secret := os.Getenv("SECRET_TOKEN")
	if secret == "" {
		return "", fmt.Errorf("could not create token: missing secret")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss": "tunes-server",
			"iat": time.Now().Unix(),
			"exp": time.Now().Add(ACCESS_TOKEN_DURATION).Unix(),
		})
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("could not sign token: %w", err)
	}
	return signed, nil
}

func CreateRefreshToken() (token string, expiration time.Time) {
	return uuid.New().String(), time.Now().Add(REFRESH_TOKEN_DURATION)
}

func ValidateToken(accessToken string) (bool, error) {
	t, err := parseToken(accessToken)

	return t.Valid, err
}

func ShouldRefresh(accessToken string) bool {
	t, err := parseToken(accessToken)

	if t.Valid {
		return false
	} else if errors.Is(err, jwt.ErrTokenExpired) {
		return true
	} else {
		return false
	}
}

func parseToken(accessToken string) (t *jwt.Token, err error) {
	return jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing algorithm: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET_TOKEN")), nil
	})
}
