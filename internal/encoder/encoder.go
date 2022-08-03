package encoder

import (
	"math"
	"math/big"
	"strconv"
)

// charSpaces constant represents the number of bytes required to encode
// a single character.
const charSpaces int = 3

// maxWordValue constant contains the maximum value that a word can has. It
// includes the default separator at first with value "1".
const maxWordValue uint64 = 1999999999999999999

// StrToInts function encodes the string provided into a slice of *big.Int's
// iterating over input characters and storing each byte number representation.
func StrToInts(input string) (encoded []*big.Int) {
	// Calculate the padding with the charSpaces constant to multiply the
	// current word value to get free space to the next character.
	var padding = uint64(math.Pow10(charSpaces))

	// Initialize the current word with a separator with value 1 and create a
	// variable to count the length of the encoded current word, then iterate
	// over input characters as bytes.
	var currentWord uint64 = 1
	// var currentWordLen int = 1
	for _, c := range []byte(input) {
		// Calculate the new word including the current character. Compare with
		// the maximum value that a word can contain, if is greater or equal,
		// append the current word and include the current character into a new
		// one, else update the current word with the new word.
		var newWord uint64 = (currentWord * padding) + uint64(c)
		if newWord >= maxWordValue {
			encoded = append(encoded, big.NewInt(int64(currentWord)))
			currentWord = padding + uint64(c)
		} else {
			currentWord = newWord
		}
	}

	// Include the last word encoded that contains one character at least.
	encoded = append(encoded, big.NewInt(int64(currentWord)))
	return
}

// IntsToStr function decode a provided slice of *big.Int's decoding each
// character byte representation into a string.
func IntsToStr(input []*big.Int) (decoded string, err error) {
	// Create a byte array to store every byte char value.
	var bytes []byte
	for _, encoded := range input {
		// Delete start mark of received value.
		var str string = encoded.String()[1:]
		var i int
		for i < len(str) {
			// Get the current string position content.
			var value string = str[i : i+charSpaces]
			// Cast the string value to integer.
			var number int
			if number, err = strconv.Atoi(value); err != nil {
				return
			}
			// Append it as byte to the current byte array.
			bytes = append(bytes, byte(number))

			i += charSpaces
			if i > len(str) {
				i = len(str) - 1
			}
		}
	}
	// Cast the byte array to string.
	decoded = string(bytes)
	return
}
