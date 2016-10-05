package main

import (
	"math/rand"
)

const letters = "abcdefghijklmnopqretuvwxyz"

func GetRandomStringBytes(size int) []byte {
	b := make([]byte, size)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return b
}
