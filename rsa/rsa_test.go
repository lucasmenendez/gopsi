package rsa

import (
	"bytes"
	"crypto/rand"
	"math/big"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	var err error
	var keys *RSAKey

	if keys, err = NewKey(1024); err != nil {
		t.Errorf("Expected success during RSA keys generation, got error: %s", err)
		return
	}

	prime, _ := rand.Prime(rand.Reader, 256)
	var input []byte = []byte(prime.Text(16))
	var cipher, plain []byte
	if cipher, err = keys.Encrypt(input); err != nil {
		t.Errorf("Expected success during input encryption, got error: %s", err)
		return
	}

	if plain, err = keys.Decrypt(cipher); err != nil {
		t.Errorf("Expected success during input decryption, got error: %s", err)
		return
	}

	var output, _ = new(big.Int).SetString(string(plain), 16)
	if prime.Cmp(output) != 0 {
		t.Errorf("Expected '%s', got '%s'", string(input), string(plain))
		return
	}
}

func TestEncryptPublicKey(t *testing.T) {
	var err error
	var key *RSAKey

	if key, err = NewKey(1024); err != nil {
		t.Errorf("Expected success during RSA key generation, got error: %s", err)
		return
	}

	var pub []byte
	if pub, err = key.PubKey(); err != nil {
		t.Errorf("Expected success during PublicKey encoding, got error: %s", err)
		return
	}

	prime, _ := rand.Prime(rand.Reader, 256)
	var input []byte = []byte(prime.Text(16))
	var cipher1, cipher2 []byte
	if cipher1, err = key.Encrypt(input); err != nil {
		t.Errorf("Expected success during input encryption, got error: %s", err)
		return
	}

	if cipher2, err = EncryptWitPublicKey(pub, input); err != nil {
		t.Errorf("Expected success during input encryption, got error: %s", err)
		return
	}

	if !bytes.Equal(cipher1, cipher2) {
		t.Errorf("Expected '%s', got '%s'", cipher1, cipher2)
		return
	}
}
