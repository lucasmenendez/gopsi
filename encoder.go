package psi

import (
	"math/big"
	"strconv"
)

// charSpaces constant represents the number of bytes required to encode
// a single character.
var charSpaces int = 3

// Encode function process the provided string and return a concatenated bytes
// numeric values into a single big.Int. It includes a starter mark with value
// 1 to ensure that every character occupies the same space.
func Encode(str string) (encoded *big.Int) {
	// Calculate the padding value (p) to multiply the result in each iteration
	// to get free space to store the next character, taking the reference of
	// the number of spaces by character: padding = 10^charSpaces
	var base10, exp *big.Int = big.NewInt(10), big.NewInt(int64(charSpaces))
	var padding *big.Int = new(big.Int).Exp(base10, exp, nil)

	// Initializes encode with a start mark "1".
	encoded = big.NewInt(1)

	// Encode char by char.
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

// Decode function process the provided big.Int value, splitting it into slice
// of bytes value representations and casting it to string.
func Decode(encoded *big.Int) (decoded string, err error) {
	// Delete start mark of received value.
	var str string = encoded.String()[1:]

	// Create a byte array to store every byte char value.
	var bytes []byte
	for i := 0; i < len(str); i += charSpaces {
		// Get the current string position content.
		var value string = str[i : i+charSpaces]
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
