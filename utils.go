package ouser

import (
	"crypto/sha512"
	"encoding/hex"
)

func sha(s string) string {
	k := sha512.Sum512([]byte(s))
	return hex.EncodeToString(k[:])[58:98]
}
