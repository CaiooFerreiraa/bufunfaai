package hash

import platformcrypto "github.com/bufunfaai/bufunfaai/apps/api/internal/platform/crypto"

type PasswordHasher struct {
	hasher *platformcrypto.Argon2idPasswordHasher
}

func NewPasswordHasher(hasher *platformcrypto.Argon2idPasswordHasher) *PasswordHasher {
	return &PasswordHasher{hasher: hasher}
}

func (hasher *PasswordHasher) HashPassword(password string) (string, error) {
	return hasher.hasher.Hash(password)
}

func (hasher *PasswordHasher) ComparePassword(encodedHash string, password string) error {
	return hasher.hasher.Compare(encodedHash, password)
}
