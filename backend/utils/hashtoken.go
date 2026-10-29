package utils

import "crypto/sha256"

// Tokens are stored as SHA-256 hashes so a leaked sessions table
// can't be used to hijack sessions. Hash raw token bytes before
// any insert or lookup.
func HashToken(token []byte) []byte {
	sum := sha256.Sum256(token)
	return sum[:]
}
