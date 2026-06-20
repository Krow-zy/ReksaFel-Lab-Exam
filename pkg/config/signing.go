package config

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
)

// GenerateHMAC computes the HMAC-SHA256 signature of a message using a key.
// Standard library crypto packages are used for secure signature computation.
func GenerateHMAC(message []byte, key []byte) []byte {
	h := hmac.New(sha256.New, key)
	_, _ = h.Write(message)
	return h.Sum(nil)
}

// VerifyHMAC checks if the HMAC-SHA256 signature matches the message.
// Uses subtle.ConstantTimeCompare to prevent timing attacks.
func VerifyHMAC(message []byte, signature []byte, key []byte) bool {
	expectedSignature := GenerateHMAC(message, key)
	return subtle.ConstantTimeCompare(signature, expectedSignature) == 1
}
