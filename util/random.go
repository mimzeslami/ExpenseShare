package util

import (
	"math/rand"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func RandomString(n int) string {
	var letters = []rune(alphabet)
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func RandomUser() string {
	return RandomString(6)
}

func RandomEmail() string {
	return RandomString(6) + "@gmail.com"
}

func RandomDatetime() time.Time {
	return time.Now().Add(time.Duration(RandomInt(0, 1000)) * time.Hour).Local()
}
