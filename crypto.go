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

// returns securely generated random bytes of the given length.
func randomBytes(length int) []byte {
	secret := make([]byte, length)
	rand.Read(secret)
	return secret
}

// hash user password into 256bit key suitable for AES encryption
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
	// append salt and nonce to ciphertext / tag
	saltPlusNonce := append(salt, nonce...)
	// append ciphertext to salt and nonce, we now have [salt][nonce][ciphertext + tag]
	encryptedData := append(saltPlusNonce, ciphertext...)

	return encryptedData, nil
}

func decryptBytes(data []byte, password []byte) ([]byte, error) {
	salt := data[:16]
	nonce := data[16:28]
	ciphertext := data[28:]
	key := deriveKey(password, salt)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to wrap cipher in GCM: %w", err)
	}
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("decryption failed (wrong password or tampered file): %w", err)
	}

	return plaintext, nil
}

func decryptFile(src string, dest string, password []byte) error {

	data, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}
	decrypted, err := decryptBytes(data, password)
	if err != nil {
		return err
	}
	err = os.WriteFile(dest, decrypted, 0600)
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}
	return nil

}

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
