package main

import (
	"math/rand"
	"strings"
	"time"
)

var src = rand.NewSource(time.Now().UnixNano())
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
const (
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
)
func RandomString(length int) string {
	sb := strings.Builder{}
	sb.Grow(length)
	for i, cache, remain := length - 1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return sb.String()
}

func RandomArray(length int) []int {
	array := make([]int, length)
	for i := 0; i < length; i++{
		array[i] = rand.Intn(100)
	}
	return array
}
