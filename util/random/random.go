package random

import (
	"crypto/rand"
	"fmt"

	log "github.com/Sirupsen/logrus"
)

// Reference: http://stackoverflow.com/a/25736155/1343356
func Salt() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Errorf("Error occurred while generating random string: ", err)
		return ""
	}

	return fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
