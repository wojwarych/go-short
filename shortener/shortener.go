package shortener

import (
	"math"
)

const alphanumerics = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func Shorten(n uint) string {
	var short string
	if n == 0 {
		return string(alphanumerics[0])
	}
	for n > 0 {
		short += string(alphanumerics[n%62])
		n /= 62
	}
	return reverse(short)
}

func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func Decoder(s string) int {
	r := []rune(s)
	var num int
	runeLen := len(r)
	base := float64(len(alphanumerics))
	for i := 0; i < runeLen; i++ {
		pow := float64(runeLen - (i + 1))
		num += indexOf(r[i]) * int(math.Pow(base, pow))
	}
	return num
}

func indexOf(c rune) int {
	for i, v := range alphanumerics {
		if v == c {
			return i
		}
	}
	return -1
}
