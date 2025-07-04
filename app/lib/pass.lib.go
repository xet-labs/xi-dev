package lib

import (
	"sync"

	"golang.org/x/crypto/bcrypt"
)

type PassLib struct {
	once sync.Once
	rw   sync.RWMutex
}

// Global singleton instance
var Pass = &PassLib{}

// Hash password
func (p *PassLib) Hash(pw string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(pw), 14)
    return string(bytes), err
}
