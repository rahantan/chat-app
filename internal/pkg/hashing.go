package pkg

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashToken(tokenStr string) string {
	tokenHash := sha256.Sum256([]byte(tokenStr))
	return hex.EncodeToString(tokenHash[:])
}
