package shortener

import (
	"math/big"
)

const maxInt64 = int64(^uint64(0) >> 1)

var counter int64 = 0

func SetCounter(c int64) {
	counter = c;
}

func ShortURLString() string {
	res := big.NewInt(counter).Text(62)
	if counter == maxInt64 {
		counter = 0
	} else {
		counter++
	}
	return res
}
