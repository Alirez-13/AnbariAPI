// # SINGLE REASON: Preserve legacy password utility functions.
package utils

import "AnbariAPI/internal/auth/infrastructure"

func HashPassword(password string) (string, error) {
	return infrastructure.NewBcryptPasswordHasher().HashPassword(password)
}

func CheckPasswordHash(password, hash string) bool {
	return infrastructure.NewBcryptPasswordHasher().CheckPasswordHash(password, hash)
}
