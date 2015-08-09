package random

import (
	"crypto/rand"
	"encoding/base64"

	log "github.com/Sirupsen/logrus"
)

func DomainSalt() string {
	s, err := GenerateRandomString(16) // 24 byte base64 encode output
	if err != nil {
		log.Errorf("Error occurred while simple salt random: ", err)
		return ""
	}

	return s
}

func GenerateBearerToken() string {
	t, err := GenerateRandomString(44) // 60 byte base64 encode output
	if err != nil {
		log.Errorf("Error occurred while generating bearer token: ", err)
		return ""
	}

	return t
}

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}
