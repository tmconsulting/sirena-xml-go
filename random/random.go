package random

import (
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

// StringWithCharset creates random string of specified length in provided charset
func StringWithCharset(length int, charset string) string {
	rand.Seed(time.Now().Unix())
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// String creates random string of specified length
func String(length int) string {
	return StringWithCharset(length, charset)
}

func Uint32() uint32 {
	return seededRand.Uint32()
}
