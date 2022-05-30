package psi

import (
	"math/big"
	"strconv"
	"strings"
)

var padding *big.Int = big.NewInt(1000)

func Encode(str string) (encoded *big.Int) {
	slice := []byte(str)
	encoded = big.NewInt(1)
	for i := 0; i < len(slice); i++ {
		currentChar := big.NewInt(int64(slice[i]))
		encoded.Mul(encoded, padding)
		encoded.Add(encoded, currentChar)
	}
	return
}

func Decode(encoded *big.Int) (decoded string, err error) {
	var str string = strings.TrimPrefix(encoded.String(), "1")
	var bytes []byte = make([]byte, len(str)/3)

	for i := 0; i < len(str); i += 3 {
		var number int
		if number, err = strconv.Atoi(str[i : i+3]); err != nil {
			return
		}
		bytes = append(bytes, byte(number))
	}
	decoded = string(bytes)
	return
}
