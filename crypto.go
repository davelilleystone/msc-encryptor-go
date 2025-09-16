package main

import (
	"fmt"
	"os"
)

// func deriveKey(password []byte, salt []byte) []byte {

// }

// func encryptBytes(data []byte, password []byte) ([]byte, error) {

// }

// func decryptBytes(data []byte, password []byte) ([]byte, error) {

// }

// func decryptFile(src string, dest string, password []byte) error {

// }

func encryptFile(src string, dest string, password []byte) error {
	data, err := os.ReadFile(src)

	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}

	// call encrypt bytes here

	err = os.WriteFile(dest, data, 0644)

	if err != nil {
		return fmt.Errorf("error writing file: %v", err)
	}

	return err
}
