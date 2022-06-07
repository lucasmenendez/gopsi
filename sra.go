package psi

import (
	"crypto/rand"
	"math/big"
)

// SRAKey struct contains the common prime number (n), the encryption key (K)
// and the decryption key (L). Read more about SRA conmutative encryption
// algorithm in the Mental Poker paper:
// https://people.csail.mit.edu/rivest/pubs/SRA81.pdf
type SRAKey struct {
	prime     *big.Int // n
	secret    *big.Int // K
	secretInv *big.Int // L
}

// GenerateKey function calculates both encryption (K) and decryptions (L) keys
// with the size provided and using the prime number (n) provided as argument.
// The key pair is calculated following SRA algorithm, where gdc(K, φ(n)) = 1
// and L = K (mod φ(n)).
func NewKey(prime *big.Int, size int) (key *SRAKey, err error) {
	var bigOne = big.NewInt(1)

	// Store common seed (large prime).
	key = &SRAKey{prime: prime}

	// Calculate φ(n), where n is the prime provided.
	phiP := new(big.Int).Sub(prime, bigOne)
	for {
		// Generate SRAKey.secret (M) candidate prime number with size provided.
		if key.secret, err = rand.Prime(rand.Reader, size); err != nil {
			return
		}

		// If the generated key.secret meets the property of gdc(K, φ(n)) = 1
		// store it as valid SRAKey.secret (M), otherwise keep trying.
		if new(big.Int).GCD(nil, nil, key.secret, phiP).Cmp(bigOne) == 0 {
			break
		}
	}

	// Calculate key.secretInv (L) as the invers of key.secret mod φ(n)
	key.secretInv = new(big.Int).ModInverse(key.secret, phiP)
	return
}

// Encrypt function encrypts the provided message (M) with the encryption key
// (K) following the formula: E(M) = M^K (mod n). Both message and result are
// big.Int.
func (key *SRAKey) Encrypt(message *big.Int) *big.Int {
	// Generates E_K(M) = M^K mod n, where M is the plain message, K is the
	// key.secret and n is the common prime provided during the key generation.
	return new(big.Int).Exp(message, key.secret, key.prime)
}

// Decrypt function decrypts the provided cipher value (C) with the decryption
// key (L) following the formula: D(C) = C^L (mod n). Both cipher value and
// result are big.Int.
func (key *SRAKey) Decrypt(cipher *big.Int) *big.Int {
	// Generates D_K(C) = C^L mod n, where C is the cipher message, L is the
	// key.secretInv and n is the common prime provided during the key
	// generation.
	return new(big.Int).Exp(cipher, key.secretInv, key.prime)
}
