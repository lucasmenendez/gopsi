package psi

import (
	"crypto/rand"
	"math/big"
)

type Key struct {
	prime     *big.Int
	secret    *big.Int
	secretInv *big.Int
}

var bigOne = big.NewInt(1)

func GenerateKey(prime *big.Int, size int) (key *Key, err error) {
	// Store common seed (large prime).
	key = &Key{prime: prime}

	phiP := new(big.Int).Sub(prime, bigOne)
	for {
		if key.secret, err = rand.Prime(rand.Reader, size); err != nil {
			return
		}
		if new(big.Int).GCD(nil, nil, key.secret, phiP).Cmp(bigOne) == 0 {
			break
		}
	}

	key.secretInv = new(big.Int).ModInverse(key.secret, phiP)
	return
}

func (key *Key) Encrypt(v *big.Int) *big.Int {
	return new(big.Int).Exp(v, key.secret, key.prime)
}

func (key *Key) Decrypt(v *big.Int) *big.Int {
	return new(big.Int).Exp(v, key.secretInv, key.prime)
}