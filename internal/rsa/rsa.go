package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"errors"
)

type RSAKey struct {
	pub  *rsa.PublicKey
	priv *rsa.PrivateKey
}

func NewKey(size int) (key *RSAKey, err error) {
	key = &RSAKey{}
	if key.priv, err = rsa.GenerateKey(rand.Reader, size); err != nil {
		return
	}

	key.priv.Precompute()
	if err = key.priv.Validate(); err != nil {
		return
	}

	key.pub = &key.priv.PublicKey
	return
}

func (key *RSAKey) Encrypt(msg []byte) (cipher []byte, err error) {
	cipher, err = rsa.EncryptOAEP(sha1.New(), rand.Reader, key.pub, msg, nil)
	return
}

func (key *RSAKey) Decrypt(cipher []byte) (msg []byte, err error) {
	msg, err = rsa.DecryptOAEP(sha1.New(), rand.Reader, key.priv, cipher, nil)
	return
}

func (key *RSAKey) PubKey() (pub []byte, err error) {
	if pub, err = x509.MarshalPKIXPublicKey(key.pub); err != nil {
		return
	}

	return
}

func EncryptWitPublicKey(pub, msg []byte) ([]byte, error) {
	var ok bool
	var key *RSAKey = &RSAKey{}
	if candidate, err := x509.ParsePKIXPublicKey(pub); err != nil {
		return nil, err
	} else if key.pub, ok = candidate.(*rsa.PublicKey); !ok {
		return nil, errors.New("error casting public key to *rsa.PublicKey")
	}

	return key.Encrypt(msg)
}
