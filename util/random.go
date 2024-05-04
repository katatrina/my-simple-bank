package util

import (
	"math/rand"
	"strings"
)

// RandomInt returns a random integer between min and max.
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString returns a random string of length n.
func RandomString(n int) string {
	var alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomOwner returns a random owner name.
func RandomOwner() string {
	return RandomString(6)
}

// RandomMoney returns a random amount of money between 0 and 1000.
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// RandomCurrency returns a random currency code.
func RandomCurrency() string {
	currencies := []string{USD, CAD, EUR}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

func RandomEmail() string {
	return RandomString(6) + "@gmail.com"
}
