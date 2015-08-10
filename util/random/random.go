package random

import (
	"crypto/rand"
	"encoding/base64"
	mr "math/rand"
	"time"

	log "github.com/Sirupsen/logrus"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
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

// Reference: http://stackoverflow.com/a/31832326/1343356
func GenerateUserPassword(n int) string {
	var src = mr.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
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
