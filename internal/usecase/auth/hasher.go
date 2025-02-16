package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type BcryptPasswordHasher struct{}

func NewBcryptPasswordHasher() *BcryptPasswordHasher {
	return &BcryptPasswordHasher{}
}

func (h *BcryptPasswordHasher) GenerateFromPassword(password []byte, cost int) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(password, cost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password using bcrypt: %w", err)
	}

	return hashedPassword, nil
}

func (h *BcryptPasswordHasher) CompareHashAndPassword(hashedPassword, password []byte) error {
	if err := bcrypt.CompareHashAndPassword(hashedPassword, password); err != nil {
		return fmt.Errorf("failed to compare hash and password using bcrypt: %w", err)
	}

	return nil
}
