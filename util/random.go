package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

var r *rand.Rand

func init() {
	source := rand.NewSource(time.Now().UnixNano())
	r = rand.New(source)
}

// RandomInt returns a random integer between min and max.
func RandomInt(min, max int64) int64 {
	return min + r.Int63n(max-min+1)
}

// RandomString returns a random string of length n.
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[r.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomOwner returns a random owner name.
func RandomOwner() string {
	return RandomString(6)
}

// RandomMoney returns a random amount of money.
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// RandomCurrency returns a random currency code.
func RandomCurrency() string {
	currencies := []string{USD, EUR, CAD}
	n := len(currencies)

	return currencies[r.Intn(n)]
}

// RandomEmail returns a random email address.
func RandomEmail() string {
	return RandomString(6) + "@gmail.com"
}
