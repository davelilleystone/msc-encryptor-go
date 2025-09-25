// Benchmarking harness for Go AES-GCM encryptor utility.
//
// Generated with the assistance of ChatGPT (GPT-5).
// All generated code has been reviewed and adapted by the author.
//
// This program benchmarks encryption and decryption times on test files
// (1MB, 10MB, 100MB, 500MB) using the existing encryptFile / decryptFile
// functions. Results are printed to console and written to benchmarks/results.csv.

package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var testFiles = []string{
	"benchmarks/test_1mb.bin",
	"benchmarks/test_10mb.bin",
	"benchmarks/test_100mb.bin",
	"benchmarks/test_500mb.bin",
}

const runsPerFile = 5
const password = "benchmark-password"
const resultsCSV = "benchmarks/results.csv"

// Wrappers to avoid name collisions with your actual implementation
func runEncrypt(src, dest, password string) error {
	return encryptFile(src, dest, []byte(password)) // call your real encryptFile
}

func runDecrypt(src, dest, password string) error {
	return decryptFile(src, dest, []byte(password)) // call your real decryptFile
}

func benchmarkFile(path string, runs int) (encryptAvg, decryptAvg float64, size int64, err error) {
	src := filepath.Clean(path)
	enc := src + ".enc"
	dec := src + ".dec"

	var totalEnc, totalDec float64

	// Get file size
	info, err := os.Stat(src)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("could not stat %s: %w", src, err)
	}
	size = info.Size()

	for i := 0; i < runs; i++ {
		// Encrypt
		start := time.Now()
		if err := runEncrypt(src, enc, password); err != nil {
			return 0, 0, 0, fmt.Errorf("encrypt failed: %w", err)
		}
		elapsed := time.Since(start).Seconds()
		totalEnc += elapsed

		// Decrypt
		start = time.Now()
		if err := runDecrypt(enc, dec, password); err != nil {
			return 0, 0, 0, fmt.Errorf("decrypt failed: %w", err)
		}
		elapsed = time.Since(start).Seconds()
		totalDec += elapsed

		// Cleanup decrypted file
		_ = os.Remove(dec)
	}

	// Cleanup encrypted file
	_ = os.Remove(enc)

	encryptAvg = totalEnc / float64(runs)
	decryptAvg = totalDec / float64(runs)

	return encryptAvg, decryptAvg, size, nil
}

func main() {
	results := [][]string{{"file", "size_bytes", "encrypt_avg", "decrypt_avg"}}

	fmt.Println("Running Go benchmarks...\n")
	for _, file := range testFiles {
		encAvg, decAvg, size, err := benchmarkFile(file, runsPerFile)
		if err != nil {
			fmt.Printf("Error benchmarking %s: %v\n", file, err)
			continue
		}

		fmt.Printf("%s: %.1f MB | Encrypt %.4fs | Decrypt %.4fs\n",
			filepath.Base(file), float64(size)/1e6, encAvg, decAvg)

		results = append(results, []string{
			filepath.Base(file),
			fmt.Sprintf("%d", size),
			fmt.Sprintf("%f", encAvg),
			fmt.Sprintf("%f", decAvg),
		})
	}

	// Write to CSV
	f, err := os.Create(resultsCSV)
	if err != nil {
		fmt.Printf("Could not create %s: %v\n", resultsCSV, err)
		return
	}
	defer f.Close()

	w := csv.NewWriter(f)
	if err := w.WriteAll(results); err != nil {
		fmt.Printf("Could not write results: %v\n", err)
	}
	fmt.Printf("\nBenchmark results saved to %s\n", resultsCSV)
}
