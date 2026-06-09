// # SINGLE REASON: Hash auth passwords with bcrypt.
package infrastructure

import "golang.org/x/crypto/bcrypt"

type BcryptPasswordHasher struct{}

func NewBcryptPasswordHasher() BcryptPasswordHasher {
	return BcryptPasswordHasher{}
}

func (BcryptPasswordHasher) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (BcryptPasswordHasher) CheckPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
