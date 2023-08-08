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
	year := RandomInt(2000, 2020)
	month := RandomInt(1, 12)
	day := RandomInt(1, 28)
	hour := RandomInt(0, 23)
	minute := RandomInt(0, 59)
	second := RandomInt(0, 59)
	return time.Date(year, time.Month(month), day, hour, minute, second, 0, time.UTC)
}
