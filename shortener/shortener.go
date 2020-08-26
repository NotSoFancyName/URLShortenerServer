// Package shortener implements two func
// First one is used to produce a short URL according to the counter
// Second one sets a counter with a provided number
package shortener

import (
	"math/big"
)

const maxInt64 = int64(^uint64(0) >> 1)

var counter int64 = 1

// Should be used to set counter during init
func SetCounter(c int64) {
	counter = c
}

// Increments counter by one and returns counter in base62 representation
func ShortURLString() string {
	res := big.NewInt(counter).Text(62)
	if counter == maxInt64 {
		counter = 0
	} else {
		counter++
	}
	return res
}
