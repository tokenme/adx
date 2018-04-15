package utils

import (
	"encoding/binary"
	"math"
	"math/big"
	"math/rand"
)

var Pow40 = Pow10(4)

func StrToBigInt(value string) (i *big.Int, ok bool) {
	return new(big.Int).SetString(value, 10)
}

func Uint64ToByte(i uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, i)
	return b
}

func ByteToUint64(b []byte) uint64 {
	return uint64(binary.BigEndian.Uint64(b))
}

func RangeRandUint64(from uint64, to uint64) uint64 {
	if from == to {
		return to
	}
	if from > to {
		from, to = to, from
	}
	return uint64(rand.Int63n(int64(to+1-from))) + from
}

func Pow10(n int) *big.Int {
	return new(big.Int).SetUint64(uint64(math.Pow10(n)))
}
