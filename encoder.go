package psi

import (
	"math/big"
	"strconv"
	"strings"
)

var padding *big.Int = big.NewInt(1000)
var positions int = len(padding.String()) - 1

func Encode(str string) (encoded *big.Int) {
	// Initializes encode with a start mark "1".
	encoded = big.NewInt(1)

	// Encode vhar by char.
	slice := []byte(str)
	for i := 0; i < len(slice); i++ {
		// Multiply current value to create padding and make each char occupy
		// the same position size to split it then.
		encoded.Mul(encoded, padding)
		// Get char number into a BigInt.
		currentChar := big.NewInt(int64(slice[i]))
		// Sum the char number to fill the padding.
		encoded.Add(encoded, currentChar)
	}
	return
}

func Decode(encoded *big.Int) (decoded string, err error) {
	// Delete start mark of received value.
	var str string = strings.TrimPrefix(encoded.String(), "1")

	// Create a byte array to store every byte char value.
	var bytes []byte
	for i := 0; i < len(str); i += positions {
		// Get the current string position content.
		var value string = str[i : i + positions]
		// Cast the string value to integer.
		var number int
		if number, err = strconv.Atoi(value); err != nil {
			return
		}
		// Append it as byte to the current byte array.
		bytes = append(bytes, byte(number))
	}
	// Cast the byte array to string.
	decoded = string(bytes)
	return
}
