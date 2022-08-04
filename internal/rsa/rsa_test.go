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

func TestEncryptDecryptPk(t *testing.T) {
	var err error
	var keys *RSAKey

	if keys, err = NewKey(1024); err != nil {
		t.Errorf("Expected success during RSA keys generation, got error: %s", err)
		return
	}

	var pk []byte
	if pk, err = keys.PubKey(); err != nil {
		t.Errorf("Expected success during public key encoding, got error: %s", err)
		return
	}

	prime, _ := rand.Prime(rand.Reader, 256)
	var input []byte = []byte(prime.Text(16))
	var cipher, plain []byte
	if cipher, err = EncryptWitPubKey(pk, input); err != nil {
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

func TestEncryptPk(t *testing.T) {
	var err error
	var keys *RSAKey

	if keys, err = NewKey(1024); err != nil {
		t.Errorf("Expected success during RSA keys generation, got error: %s", err)
		return
	}

	var pk []byte
	if pk, err = keys.PubKey(); err != nil {
		t.Errorf("Expected success during public key encoding, got error: %s", err)
		return
	}

	prime, _ := rand.Prime(rand.Reader, 256)
	var input []byte = []byte(prime.Text(16))
	var cipher1, cipher2 []byte
	if cipher1, err = keys.Encrypt(input); err != nil {
		t.Errorf("Expected success during input encryption 1, got error: %s", err)
		return
	}

	if cipher2, err = EncryptWitPubKey(pk, input); err != nil {
		t.Errorf("Expected success during input encryption 2, got error: %s", err)
		return
	}

	if bytes.Equal(cipher1, cipher2) {
		t.Errorf("Expected '%s', got '%s'", string(cipher1), string(cipher2))
		return
	}
}
