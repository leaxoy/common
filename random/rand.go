package random

import (
	"math/rand"
)

const (
	DefaultCharset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
)

func StringWithCharset(size int, charset string) string {
	length := len(charset)
	bytes := make([]byte, 0, size)
	for i := 0; i < size; i++ {
		bytes = append(bytes, charset[rand.Intn(length)])
	}
	return string(bytes)
}

func String(size int) string {
	return StringWithCharset(size, DefaultCharset)
}
