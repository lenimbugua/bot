package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"
const numeric = "0123456789"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)
	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

func RandomTitle() string {
	return RandomString(20)
}
func RandomSlug() string {
	return RandomString(10)
}
func RandomDescription() string {
	return RandomString(20)
}

// Random phone number
func RandomDigits(n int) string {
	var sb strings.Builder
	k := len(numeric)
	for i := 0; i < n; i++ {
		c := numeric[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

func RandomPhoneNumber() string {
	return "+" + RandomDigits(10)
}

func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}
