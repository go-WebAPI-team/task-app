package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashPasswowrd(pw string) string {
	hash := sha256.Sum256([]byte(pw))
	return hex.EncodeToString(hash[:])
}