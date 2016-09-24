package main

import (
	"math/rand"
)

const letters = "abcdefghijklmnopqretuvwxyz"

func GetRandomString(size int) string {
	b := make([]byte, size)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
