package helpers

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password []byte) string {
	bytes, err := bcrypt.GenerateFromPassword(password, 10)
	if err != nil {
		log.Fatal("Error generating a hashed password", err)
	}
	return string(bytes)
}

func ComparePassword(hashedPassword []byte, password []byte) bool {
	var matched bool
	err := bcrypt.CompareHashAndPassword(hashedPassword, password)

	if err != nil {
		matched = false
	} else {
		matched = true
	}
	return matched
}
