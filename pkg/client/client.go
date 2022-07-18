package client

import (
	"errors"
	"math/big"

	"github.com/lucasmenendez/psi/internal/encoder"
	"github.com/lucasmenendez/psi/internal/rsa"
	"github.com/lucasmenendez/psi/pkg/sra"
)

type Client struct {
	clientKeys     *rsa.RSAKey
	PublicKey      []byte
	CommonPrime    *big.Int
	CommonPrimeEnc []byte
	sraKey         *sra.SRAKey
	Records        []*big.Int
}

func Init() (client *Client, err error) {
	client = &Client{}

	// Generate RSA keys pair
	client.clientKeys, err = rsa.NewKey(1024)
	if err != nil {
		return
	}

	// Get public key
	client.PublicKey, err = client.clientKeys.PubKey()
	return
}

func (client *Client) AddCommonPrime(encryptedPrime []byte) (err error) {
	client.CommonPrimeEnc, err = client.clientKeys.Decrypt(encryptedPrime)
	if err != nil {
		return
	}

	var sCommonPrime string = string(client.CommonPrimeEnc)
	client.CommonPrime, _ = new(big.Int).SetString(sCommonPrime, 16)

	client.sraKey, err = sra.NewKey(client.CommonPrime, 32)
	return
}

func (client *Client) LoadData(data []string) error {
	if client.sraKey == nil {
		return errors.New("common prime not defined")
	}

	client.Records = make([]*big.Int, len(data))
	for i, item := range data {
		encoded := encoder.StrToInt(item)
		client.Records[i] = client.sraKey.Encrypt(encoded)
	}

	return nil
}

func (client *Client) EncryptInput(input []*big.Int) (output []*big.Int, err error) {
	if client.sraKey == nil {
		return nil, errors.New("common prime not defined")
	}

	for _, item := range input {
		encrypted := client.sraKey.Encrypt(item)
		output = append(output, encrypted)
	}

	return
}
