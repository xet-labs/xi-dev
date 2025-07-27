package lib

import (
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthLib struct {
	jwtSecret []byte

	once sync.Once
	rw   sync.RWMutex
}

var Auth = &AuthLib{
	jwtSecret: []byte("supersecretkey"),
}

func (a *AuthLib) HashPass(pw string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(hash), err
}

func (a *AuthLib) CheckPass(hash, pw string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw))
}

func (a *AuthLib) GenToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.jwtSecret)
}

func (a *AuthLib) ParseToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		return a.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return token.Claims.(jwt.MapClaims), nil
}

func (a *AuthLib) Hash(pw string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(pw), 14)
    return string(bytes), err
}
