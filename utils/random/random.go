package random

import (
	cryptrand "crypto/rand"
	"math/big"
	rand "math/rand"
	"time"
)

// RandomString generates a random string of given length
func RandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano()) // Seed the random generator

	result := make([]byte, n)
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

func SecureRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, n)

	for i := range result {
		num, _ := cryptrand.Int(cryptrand.Reader, big.NewInt(int64(len(letters))))
		result[i] = letters[num.Int64()]
	}
	return string(result)
}

func SecureRandomInt(max int64) int64 {
	n, _ := cryptrand.Int(cryptrand.Reader, big.NewInt(max))
	return n.Int64()
}
