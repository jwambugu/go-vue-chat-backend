package factory

import "github.com/brianvoe/gofakeit/v6"

// RandomString generates a random string with length n
func RandomString(n int) string {
	return gofakeit.LetterN(uint(n))
}
