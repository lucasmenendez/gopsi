package psi

import (
	"crypto/rand"
	"math/big"
)

type Key struct {
	prime     *big.Int // n
	secret    *big.Int // K
	secretInv *big.Int // L
}

var bigOne = big.NewInt(1)

func GenerateKey(prime *big.Int, size int) (key *Key, err error) {
	// Store common seed (large prime).
	key = &Key{prime: prime}

	// Calculate φ(n), where n is the prime provided.
	phiP := new(big.Int).Sub(prime, bigOne)
	for {
		// Generate key.secret candidate prime number with size provided
		if key.secret, err = rand.Prime(rand.Reader, size); err != nil {
			return
		}

		// If the generated key.secret meets the property of gdc(K, φ(n)) = 1
		// store it as valid key.secret, otherwise keep trying with other prime.
		if new(big.Int).GCD(nil, nil, key.secret, phiP).Cmp(bigOne) == 0 {
			break
		}
	}

	// Calculate key.secretInv as the invers of key.secret mod φ(n)
	key.secretInv = new(big.Int).ModInverse(key.secret, phiP)
	return
}

func (key *Key) Encrypt(v *big.Int) *big.Int {
	// Generates E_K(M) = M^K mod n, where M is the plain message, K is the
	// key.secret and n is the common prime provided during the key generation.
	return new(big.Int).Exp(v, key.secret, key.prime)
}

func (key *Key) Decrypt(v *big.Int) *big.Int {
	// Generates D_K(C) = C^L mod n, where C is the cipher message, L is the
	// key.secretInv and n is the common prime provided during the key
	// generation.
	return new(big.Int).Exp(v, key.secretInv, key.prime)
}
