package hw02unpackstring

import (
	"errors"
	"strconv"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	return unpack(str)
}

func unpack(str string) (string, error) {
	var res = []rune{}
	var prevRune *rune
	for _, v := range str {
		if prevRune == nil {
			if unicode.IsDigit(v) {
				return "", ErrInvalidString
			} else {
				prevRune = &v
			}
		} else {
			if unicode.IsDigit(v) {
				n, err := strconv.Atoi(string(v))
				if err != nil {
					return "", ErrInvalidString
				}
				for range n {
					res = append(res, *prevRune)
				}
				prevRune = nil
			} else {
				res = append(res, *prevRune)
				prevRune = &v
			}
		}
	}
	if prevRune != nil {
		res = append(res, *prevRune)
	}
	return string(res), nil
}
