package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"os"

	"golang.org/x/crypto/pbkdf2"
)

// randomBytes returns securely generated random bytes of the given length.
// Uses crypto/rand and will crash only if the system CSPRNG fails irrecoverably.

func randomBytes(length int) []byte {
	secret := make([]byte, length)
	rand.Read(secret)
	return secret
}

func deriveKey(password []byte, salt []byte) []byte {
	iterations := 1200000
	keyLength := 32 // 256 bits
	key := pbkdf2.Key(password, salt, iterations, keyLength, sha256.New)
	return key
}

func encryptBytes(data []byte, password []byte) ([]byte, error) {
	salt := randomBytes(16)
	nonce := randomBytes(12)
	key := deriveKey(password, salt)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to wrap cipher in GCM: %w", err)
	}
	ciphertext := gcm.Seal(nil, nonce, data, nil)
	return ciphertext, nil
}

// func decryptBytes(data []byte, password []byte) ([]byte, error) {

// }

// func decryptFile(src string, dest string, password []byte) error {

// }

func encryptFile(src string, dest string, password []byte) error {
	data, err := os.ReadFile(src)

	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}
	encrypted, err := encryptBytes(data, password)
	if err != nil {
		return err
	}
	err = os.WriteFile(dest, encrypted, 0600)
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}
	return nil
}
