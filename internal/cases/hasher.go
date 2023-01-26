package cases

import (
	"github.com/zeebo/xxh3"
)

const (
	alphabet    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	alphabetLen = 63
)

func Hasher(originalLink string) string {
	origLink := []byte(originalLink)
	hsum := xxh3.Hash(origLink)

	lowBytes := make([]byte, 0, 10)
	for hsum > 0 && len(lowBytes) < 10 {
		lowBytes = append(lowBytes, alphabet[hsum%alphabetLen])
		hsum /= alphabetLen
	}

	var highBytes []byte
	for i := len(lowBytes); i < 10; i++ {
		highBytes = append(highBytes, alphabet[0])
	}

	shortLink := append(highBytes, lowBytes...)

	return string(shortLink)
}
