package sra

import (
	"crypto/rand"
	"math/big"
	"testing"
)

func TestGenerateKey(t *testing.T) {
	var bigOne = big.NewInt(1)
	prime, _ := rand.Prime(rand.Reader, 256)
	phiP := new(big.Int).Sub(prime, bigOne)

	var key1, key2 *SRAKey
	var err error
	key1, err = NewKey(prime, 32)
	if err != nil {
		t.Errorf("Error generating key 1 (prime: %s", prime.String())
		return
	}

	key2, err = NewKey(prime, 32)
	if err != nil {
		t.Errorf("Error generating key 2 (prime: %s", prime.String())
		return
	}

	if key1.secret.Cmp(key2.secretInv) == 0 {
		t.Errorf("Both keys are the same.")
	}

	if new(big.Int).GCD(nil, nil, key1.secret, phiP).Cmp(bigOne) != 0 {
		t.Errorf("The generated key1 does not meet the property of gdc(K, φ(n)) = 1, where K = key1.secret.")
		return
	}

	if new(big.Int).GCD(nil, nil, key2.secret, phiP).Cmp(bigOne) != 0 {
		t.Errorf("The generated key2 does not meet the property of gdc(K, φ(n)) = 1, where K = key2.secret.")
		return
	}
}
