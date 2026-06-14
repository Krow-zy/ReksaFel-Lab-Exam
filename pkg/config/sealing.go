package config

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

// ConfigSealer manages local configuration encryption and decryption.
type ConfigSealer struct {
	key []byte
}

// NewConfigSealer instantiates a sealer with a 32-byte key for AES-256.
func NewConfigSealer(key []byte) (*ConfigSealer, error) {
	if len(key) != 32 {
		return nil, errors.New("key must be exactly 32 bytes for AES-256-GCM")
	}
	return &ConfigSealer{key: key}, nil
}

// Seal encrypts raw configuration bytes using AES-GCM with a random nonce.
func (cs *ConfigSealer) Seal(plainText []byte) ([]byte, error) {
	block, err := aes.NewCipher(cs.key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Nonce size for GCM is typically 12 bytes
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// Seal appends the ciphertext to the nonce so it can be unpacked later
	cipherText := aesGCM.Seal(nonce, nonce, plainText, nil)
	return cipherText, nil
}

// Unseal decrypts sealed configuration bytes using the appended nonce.
func (cs *ConfigSealer) Unseal(cipherText []byte) ([]byte, error) {
	block, err := aes.NewCipher(cs.key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()
	if len(cipherText) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, actualCipherText := cipherText[:nonceSize], cipherText[nonceSize:]
	plainText, err := aesGCM.Open(nil, nonce, actualCipherText, nil)
	if err != nil {
		return nil, err
	}

	return plainText, nil
}
