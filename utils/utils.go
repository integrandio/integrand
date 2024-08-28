package utils

import (
	"math/rand"
	"os"
	"strconv"
	"time"
	"unsafe"
)

// see if a list of strings contains a certain string
func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func GetEnvrionmentVariableString(env string, substitueVariable string) string {
	value := os.Getenv(env)
	if value == "" {
		os.Setenv(env, substitueVariable)
		return substitueVariable
	} else {
		return value
	}
}

func GetEnvrionmentVariableInt(env string, substitueValue int) int {
	subStringValue := strconv.Itoa(substitueValue)
	strValue := os.Getenv(env)
	if strValue != "" {
		value, err := strconv.Atoi(strValue)
		if err == nil {
			return value
		}
	}
	os.Setenv(env, subStringValue)
	return substitueValue
}

// Given a key value map of strings, get the value of the given key and return the string
func GetOrDefaultString(m map[string]interface{}, key string, defaultStr string) string {
	if value, ok := m[key]; ok {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return defaultStr
}

// Given a key value map of ints, get the value of the given key and return the int
func GetOrDefaultInt(m map[string]int, key string, defaultInt int) int {
	if num, ok := m[key]; ok {
		return num
	}
	return defaultInt
}

// Random String Generator stuf.....
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

func RandomString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}
