package utility

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
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

func GetContextValue(c *gin.Context, key string) (interface{}, error) {
	value, exists := c.Get(key)
	if !exists {
		return nil, fmt.Errorf("%s not found in context", key)
	}

	return value, nil
}
