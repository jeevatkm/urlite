package util

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	return HashPasswordWithCost(password, bcrypt.DefaultCost)
}

func HashPasswordWithCost(password string, cost int) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		log.Errorf("Unable to hash give password string: %v", err)
	}
	return fmt.Sprintf(string(hashedPassword))
}

func ComparePassword(upass, spass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(upass), []byte(spass))
	if err != nil {
		log.Error("PasswordHash and Passsword doesn't match")
		return false
	}
	return true // nil means PasswordHash and Passsword is a match
}
