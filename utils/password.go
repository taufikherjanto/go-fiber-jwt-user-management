package utils

import "golang.org/x/crypto/bcrypt"

func GeneratePassword(password string) string {
	// Create hash from password
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash)
}

func ComparePassword(hashedPassword string, password string) bool {
	// compare hash andd password
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
