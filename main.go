package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"syscall"

	"golang.org/x/term"
)

func exitWithError(msg string, a ...any) {
	fmt.Printf("Error: "+msg+"\n", a...)
	flag.Usage()
	os.Exit(1)
}

func main() {
	// define flag
	action := flag.String("action", "", "Select action to perform, either 'encrypt' or 'decrypt'")
	src := flag.String("src", "", "Select source file")
	dest := flag.String("dest", "", "Select destination file")

	// custom usage output
	flag.Usage = func() {
		fmt.Println("Example: encryptor -action encrypt -src sample.txt -dest sample.enc")
		flag.PrintDefaults()
	}

	// parse CLI args
	flag.Parse()

	// validate action flag
	if *action == "" {
		exitWithError("-action flag is required.")
	} else if *action != "encrypt" && *action != "decrypt" {
		exitWithError("invalid -action value %q", *action)
	}

	// validate src flag
	if *src == "" {
		exitWithError("-src flag is required.")
	}

	// validate dest flag
	if *dest == "" {
		exitWithError("-dest flag is required.")
	}
	// success
	fmt.Println("Action:", *action)
	fmt.Println("Source file:", *src)
	fmt.Println("Destination file:", *dest)

	// get password
	fmt.Printf("\nPlease enter password to %s %s:\n", *action, *src)
	password, err := term.ReadPassword(int(syscall.Stdin))

	if err != nil {
		log.Fatalf("Failed to read password: %v", err)
	}

	passwordStr := string(password)

	if *action == "encrypt" {
		fmt.Println("call encrypt routine")
	} else if *action == "decrypt" {
		fmt.Println("call decrypt routine")
	} else {
		fmt.Println("something has gone wrong")
	}

}
