package hash

import "math/rand/v2"

func Generate() string {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	hashStr := make([]rune, 16)
	for i := range hashStr {
		hashStr[i] = letterRunes[rand.IntN(len(letterRunes))]
	}
	return string(hashStr)
}
