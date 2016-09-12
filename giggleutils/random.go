package giggleutils

import (
	"math/rand"
	"time"
)

const (
	cChars = "abcdefghijklmnopqrstuvwxyz0123456789"
)

func RandomString(length int) string {
	rand.Seed(time.Now().UTC().UnixNano())

	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = cChars[rand.Intn(len(cChars))]
	}

	return string(result)
}
