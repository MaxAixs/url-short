package handlers

import (
	"math/rand"
	"time"
)

func aliasIsEmpty(alias string) bool {
	if alias == "" {
		return true
	}

	return false
}

func generateRandomAlias(size int) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz")

	b := make([]rune, size)
	for i := range b {
		b[i] = chars[rnd.Intn(len(chars))]
	}

	return string(b)
}
