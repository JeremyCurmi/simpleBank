package utils

import (
	"math/rand"
	"strings"
	"time"
)

const alp = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt generates a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alp)

	for i := 0; i < n; i++ {
		char := alp[rand.Intn(k)]
		sb.WriteByte(char)
	}
	return sb.String()
}

func RandomUserID() int64 {
	return RandomInt(0, 1000)
}

func RandomAccountName() string {
	return RandomString(8)
}

func RandomMoney() int64 {
	return RandomInt(0, 10000)
}

func RandomCurrency() string {
	currencies := []string{"EUR", "USD", "CAD"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
