package gojotids

import (
	"fmt"
	"strings"

	hashids "github.com/speps/go-hashids"
)

const ALPHABET = "abcdefghijklmnopqrstuvwxyz0123456789 !@#$%^&*()-=_+"

func Encode(s, salt string) (enc string, err error) {
	s = strings.ToLower(s)
	s = strings.TrimSpace(s)
	hd := hashids.NewData()
	hd.Salt = salt
	h, err := hashids.NewWithData(hd)
	if err != nil {
		return
	}
	ints := make([]int, len(s))
	it := 0
	for i := range s {
		if strings.Contains(ALPHABET, string(s[i])) {
			ints[it] = strings.Index(ALPHABET, string(s[i]))
			it++
		}
	}
	ints = ints[:it]
	fmt.Println(ints)
	enc, err = h.Encode(ints)
	return
}

func Decode(s, salt string) (dec string, err error) {
	hd := hashids.NewData()
	hd.Salt = salt
	h, err := hashids.NewWithData(hd)
	if err != nil {
		return
	}
	ints, err := h.DecodeWithError(s)
	if err != nil {
		return
	}
	fmt.Println(ints)
	dec = ""
	for _, i := range ints {
		dec += string(ALPHABET[i])
	}
	return
}
