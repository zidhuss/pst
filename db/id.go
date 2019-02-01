package db

import (
	"bytes"
	"fmt"
)

const chars = "012345689abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var charsLen = uint64(len(chars))

func intToID(n uint64) string {
	id := []byte{}

	for n/charsLen > 0 {
		id = append(id, byte(chars[n%charsLen]))
		n /= charsLen
	}

	id = append(id, byte(chars[n%charsLen]))

	// Reverse string
	for i, j := 0, len(id)-1; i < j; i, j = i+1, j-1 {
		id[i], id[j] = id[j], id[i]
	}

	return string(id)
}

func idToInt(id string) (uint64, error) {
	b := []byte(id)

	// Reverse string
	for i, j := 0, len(id)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}

	e := uint64(1)
	n := uint64(0)
	for _, x := range b {
		idx := bytes.IndexByte([]byte(chars), x)
		if idx == -1 {
			return 0, fmt.Errorf("IDtoInt: id contains illegal character %c", x)
		}
		n += uint64(idx) * e
		e *= charsLen
	}
	return n, nil
}
