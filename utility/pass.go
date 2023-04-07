package utility

import (
	"math/rand"
	"strings"
	"time"
)

func RandomPassword(length int) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	source := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(source)

	var password strings.Builder
	for i := 0; i < length; i++ {
		password.WriteByte(chars[rand.Intn(len(chars))])
	}
	return password.String()
}
