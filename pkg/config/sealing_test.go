package config

import (
	"bytes"
	"testing"
)

func TestConfigSealing(t *testing.T) {
	key := []byte("01234567890123456789012345678901") // 32 bytes
	sealer, err := NewConfigSealer(key)
	if err != nil {
		t.Fatalf("Failed to create sealer: %v", err)
	}

	// Test case 1: Successful encrypt and decrypt
	originalData := []byte("my-secret-database-config-credentials")
	sealedData, err := sealer.Seal(originalData)
	if err != nil {
		t.Fatalf("Failed to seal data: %v", err)
	}

	unsealedData, err := sealer.Unseal(sealedData)
	if err != nil {
		t.Fatalf("Failed to unseal data: %v", err)
	}

	if !bytes.Equal(originalData, unsealedData) {
		t.Errorf("Expected decrypted data to match original data, got %s", string(unsealedData))
	}

	// Test case 2: Attempting to decrypt tampered data should fail
	sealedData[len(sealedData)-1] ^= 0xFF // Flip last bit of the auth tag
	_, err = sealer.Unseal(sealedData)
	if err == nil {
		t.Error("Expected decryption of tampered data to fail, but it succeeded")
	}
}
