package ouser

import (
	"crypto/sha512"
	"encoding/hex"
	"strings"
)

func sha(s string) string {
	k := sha512.Sum512([]byte(s))
	return hex.EncodeToString(k[:])[58:98]
}
func getType(s ...string) string {
	return strings.Join(s, ".")
}
