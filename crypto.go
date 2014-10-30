package main

import (
	"code.google.com/p/go.crypto/bcrypt"
)

// GenerateHash generates bcrypt hash from a plaintext input value.
func GenerateHash(input string) ([]byte, error) {
	hex := []byte(input)

	// Hash the password using the default cost
	hashedPassword, err := bcrypt.GenerateFromPassword(hex, bcrypt.DefaultCost)
	if err != nil {
		return hashedPassword, err
	}
	return hashedPassword, nil
}

// CompareHash compares a bcrypt hash with a plaintext input value.
// Returns true if the hash matches the input value.
func CompareHash(digest []byte, input string) bool {
	hex := []byte(input)
	if err := bcrypt.CompareHashAndPassword(digest, hex); err == nil {
		return true
	}
	return false
}
