package random

import (
	"math/rand"
	"time"
)

var characters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

func NewString(length int) string {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	buffer := make([]rune, length)
	for i := range buffer {
		buffer[i] = characters[rand.Intn(len(characters))]
	}
	return string(buffer)
}
