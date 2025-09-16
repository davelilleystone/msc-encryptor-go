package main

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"os"

	"golang.org/x/crypto/pbkdf2"
)

func deriveKey(password []byte, salt []byte) []byte {
	iterations := 1200000
	keyLength := 32 // 256 bits
	key := pbkdf2.Key(password, salt, iterations, keyLength, sha256.New)
	return key
}

func encryptBytes(data []byte, password []byte) ([]byte, error) {
	salt := make([]byte, 16)
	rand.Read(salt)
	key := deriveKey(password, salt)
	fmt.Println(key)

	return data, nil
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

	encryptBytes(data, password)

	err = os.WriteFile(dest, data, 0600)

	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	return nil
}
