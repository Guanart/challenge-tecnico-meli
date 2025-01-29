package main

import (
	"math/rand"
	"time"
)

// Source: https://www.bacancytechnology.com/qanda/golang/generate-a-random-string-of-a-fixed-length-in-go
func randomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}
