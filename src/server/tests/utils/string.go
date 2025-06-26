package utils_test

import (
	"math/rand"
	"strings"
	"unicode"
)

func GenerateMatchingString(substr string) string {
	var b strings.Builder
	for _, r := range substr {
		if rand.Intn(2) == 1 {
			b.WriteRune(unicode.ToUpper(r))
		} else {
			b.WriteRune(unicode.ToLower(r))
		}
	}
	mixed := b.String()

	letters := "abcdefghijklmnopqrstuvwxyz"
	randLen := func() int { return rand.Intn(6) }

	randStr := func(n int) string {
		sb := make([]byte, n)
		for i := range sb {
			sb[i] = letters[rand.Intn(len(letters))]
		}
		return string(sb)
	}

	prefix := randStr(randLen())
	suffix := randStr(randLen())

	return prefix + mixed + suffix
}

func GenerateRandomEmail() string {
	letters := "abcdefghijklmnopqrstuvwxyz"
	domain := "example.com"

	randStr := func(n int) string {
		sb := make([]byte, n)
		for i := range sb {
			sb[i] = letters[rand.Intn(len(letters))]
		}
		return string(sb)
	}

	username := randStr(8)
	return username + "@" + domain
}

func GenerateRandomString(length int) string {
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	sb := make([]byte, length)
	for i := range sb {
		sb[i] = letters[rand.Intn(len(letters))]
	}
	return string(sb)
}
