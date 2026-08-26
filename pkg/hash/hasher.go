package hash

import (
	"errors"
	"fmt"

	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	"golang.org/x/crypto/bcrypt"
)

type BcryptHasher struct{}

func NewBcryptHasher() *BcryptHasher {
	return &BcryptHasher{}
}

func (b *BcryptHasher) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", exception.ErrInternal)
	}

	return string(hash), nil
}

func (b *BcryptHasher) Compare(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return exception.ErrInvalidPassword
		}
		return fmt.Errorf("password comparison failed: %w", exception.ErrInternal)
	}

	return nil
}
