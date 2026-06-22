package telemetry

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"testing"
)

func TestVerifyFileIntegrity(t *testing.T) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "integrity_test_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	content := []byte("zerogap-exam-integrity-validation-payload-2026")
	if _, err := tmpFile.Write(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Compute expected hash
	h := sha256.New()
	h.Write(content)
	expectedHash := hex.EncodeToString(h.Sum(nil))

	// Verify integrity
	isValid, err := VerifyFileIntegrity(tmpFile.Name(), expectedHash)
	if err != nil {
		t.Fatalf("VerifyFileIntegrity failed: %v", err)
	}
	if !isValid {
		t.Error("Expected file integrity to be valid")
	}

	// Verify corruption detection
	isValid, err = VerifyFileIntegrity(tmpFile.Name(), "incorrecthash1234567890abcdef")
	if err != nil {
		t.Fatalf("VerifyFileIntegrity failed: %v", err)
	}
	if isValid {
		t.Error("Expected file integrity to be invalid for corrupted hash")
	}
}
