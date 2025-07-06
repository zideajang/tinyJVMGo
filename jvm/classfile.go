package jvm

import (
	"encoding/binary"
	"fmt"
	"os"
)

// MagicNumber represents the expected magic number for Java class files.
const MagicNumber uint32 = 0xCAFEBABE

// ReadClassMagic reads the first 4 bytes of a given file,
// interprets them as a big-endian uint32, and checks if it matches
// the Java class file magic number.
//
// It returns the read magic number and an error if one occurs.
func ReadClassMagic(filepath string) (uint32, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return 0, fmt.Errorf("error opening file '%s': %w", filepath, err)
	}
	defer file.Close() // Ensure the file is closed

	buffer := make([]byte, 4) // Buffer to read 4 bytes

	n, err := file.Read(buffer)
	if err != nil {
		return 0, fmt.Errorf("error reading from file '%s': %w", filepath, err)
	}

	if n < 4 {
		return 0, fmt.Errorf("file '%s' is too small; expected at least 4 bytes, got %d", filepath, n)
	}

	// Java class files use big-endian byte order
	magic := binary.BigEndian.Uint32(buffer)

	return magic, nil
}

// VerifyClassMagic reads and verifies the magic number of a class file.
// It returns true if the magic number is valid, false otherwise,
// and an error if an I/O issue occurs.
func VerifyClassMagic(filepath string) (bool, error) {
	magic, err := ReadClassMagic(filepath)
	if err != nil {
		return false, err
	}

	if magic != MagicNumber {
		return false, nil // Magic number is invalid but no I/O error
	}

	return true, nil // Magic number is valid
}