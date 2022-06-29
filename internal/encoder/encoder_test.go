package encoder

import (
	"math/big"
	"testing"
)

func TestStrToInt(t *testing.T) {
	input := "K_4.@m"
	expected := big.NewInt(1075095052046064109)
	result := StrToInt(input)

	if expected.Cmp(result) != 0 {
		t.Errorf("Expected '%s', got '%s'", expected.String(), result.String())
		return
	}
}

func TestIntToStr(t *testing.T) {
	input := big.NewInt(1075095052046064109)
	expected := "K_4.@m"
	if result, err := IntToStr(input); err != nil {
		t.Errorf("Error decoding input: %s", err)
		return
	} else if expected != result {
		t.Errorf("Expected '%s', got '%s'", expected, result)
		return
	}
}
