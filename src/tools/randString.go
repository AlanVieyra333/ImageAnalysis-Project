package tools

import (
	"math/rand"
	"time"
)

const LETTERS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const NUMBERS = "0123456789"

func GetRandString(n int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, n)
	for i := range b {
		b[i] = LETTERS[r.Intn(len(LETTERS))]
	}
	return string(b)
}

func GetRandInteger(n int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, n)
	for i := range b {
		b[i] = NUMBERS[r.Intn(len(NUMBERS))]
	}
	return string(b)
}
