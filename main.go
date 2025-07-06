package main

import (
	"fmt"
	"tiny-jvm/jvm" // Import our jvm package
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: tiny-jvm <class_file_path>")
		os.Exit(1)
	}

	filepath := os.Args[1]

	isValid, err := jvm.VerifyClassMagic(filepath)
	if err != nil {
		fmt.Printf("Error verifying magic number for '%s': %v\n", filepath, err)
		os.Exit(1)
	}

	if isValid {
		fmt.Printf("Successfully verified magic number for '%s'. It is a valid Java class file.\n", filepath)
	} else {
		fmt.Printf("Magic number for '%s' is invalid. It is NOT a valid Java class file.\n", filepath)
	}

	// You can also read the raw magic number
	magic, err := jvm.ReadClassMagic(filepath)
	if err != nil {
		fmt.Printf("Error reading magic number: %v\n", err)
	} else {
		fmt.Printf("Read raw magic number: 0x%X\n", magic)
		fmt.Printf("Expected magic number: 0x%X\n", jvm.MagicNumber)
	}
}