package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

//Returns the bcrypt of the hashed password
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to recover hashed password: %v", err)
	}
	return string(hashedPassword), nil
}

//checks if the provided password is correct or not
func CheckPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
