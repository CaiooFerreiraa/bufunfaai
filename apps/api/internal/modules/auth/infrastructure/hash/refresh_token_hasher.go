package hash

import (
	"crypto/sha256"
	"encoding/hex"

	platformcrypto "github.com/bufunfaai/bufunfaai/apps/api/internal/platform/crypto"
)

type RefreshTokenManager struct{}

func NewRefreshTokenManager() *RefreshTokenManager {
	return &RefreshTokenManager{}
}

func (manager *RefreshTokenManager) GenerateToken() (string, error) {
	return platformcrypto.GenerateOpaqueToken(32)
}

func (manager *RefreshTokenManager) HashToken(rawToken string) string {
	sum := sha256.Sum256([]byte(rawToken))
	return hex.EncodeToString(sum[:])
}
