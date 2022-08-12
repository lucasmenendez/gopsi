package encoder

import (
	"math/big"
)

// StrToInts function encodes the string provided into a slice of *big.Int's
// iterating over input characters and storing each byte number representation.
func StrToInts(input string) (encoded []*big.Int) {
	encoded = make([]*big.Int, len(input))
	for i, char := range []byte(input) {
		encoded[i] = new(big.Int).SetUint64(uint64(char))
	}

	return
}

// IntsToStr function decode a provided slice of *big.Int's decoding each
// character byte representation into a string.
func IntsToStr(input []*big.Int) string {
	var decoded []byte = make([]byte, len(input))
	for i, char := range input {
		decoded[i] = byte(char.Uint64())
	}

	return string(decoded)
}
