package config

import (
	"bytes"
	"testing"
)

func TestHMACSigningAndVerification(t *testing.T) {
	key := []byte("secret-hmac-key-for-payload-verification")
	message := []byte("this is a secure telemetry payload message representing computer hardware status logs")

	// 1. Generate HMAC
	sig := GenerateHMAC(message, key)
	if len(sig) != 32 { // SHA-256 HMAC produces exactly 32 bytes signature
		t.Errorf("expected HMAC length of 32 bytes, got %d", len(sig))
	}

	// 2. Verification success case
	if !VerifyHMAC(message, sig, key) {
		t.Error("HMAC verification failed for valid signature")
	}

	// 3. Verification failure with modified message
	modifiedMsg := []byte("this is a secure telemetry payload message representing computer hardware status logx")
	if VerifyHMAC(modifiedMsg, sig, key) {
		t.Error("expected HMAC verification to fail for modified message, but it passed")
	}

	// 4. Verification failure with wrong key
	wrongKey := []byte("wrong-hmac-key-for-payload-verification")
	if VerifyHMAC(message, sig, wrongKey) {
		t.Error("expected HMAC verification to fail for incorrect key, but it passed")
	}

	// 5. Verification failure with modified signature
	badSig := bytes.Clone(sig)
	badSig[0] ^= 0xFF
	if VerifyHMAC(message, badSig, key) {
		t.Error("expected HMAC verification to fail for modified signature, but it passed")
	}
}
