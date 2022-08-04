package client

import (
	"crypto/rand"
	"errors"
	"math/big"

	"github.com/lucasmenendez/psi/internal/encoder"
	"github.com/lucasmenendez/psi/internal/rsa"
	"github.com/lucasmenendez/psi/pkg/bloomfilter"
	"github.com/lucasmenendez/psi/pkg/sra"
)

// Client struct contains all required parameters to perform a private set
// intersection over another knowed Client.
type Client struct {
	CommonPrime *big.Int
	sraKey      *sra.SRAKey
	rsaKey      *rsa.RSAKey
	filter      *bloomfilter.BloomFilter
}

// Init function instances a Client generating a new RSA key pair.
func Init() (client *Client, err error) {
	client = &Client{}

	// Generate RSA keys pair
	client.rsaKey, err = rsa.NewKey(1024)
	return
}

// PubKey function returns the current client instance RSA public key byte slice
// to be shared to the other client. It allows to share a common prime securely.
func (client *Client) PubKey() ([]byte, error) {
	if client.rsaKey == nil {
		return nil, errors.New("client not initialized")
	}

	return client.rsaKey.PubKey()
}

// GenEncryptedPrime function generates a common prime number to share with
// other client and encrypts it with the RSA public key provided. It also try to
// initialize the SRA key with the common prime generated.
func (client *Client) GenEncryptedPrime(extKey []byte) ([]byte, error) {
	var err error
	if client.sraKey != nil && client.CommonPrime != nil {
		return nil, errors.New("common prime already defined, create a new instance")
	}

	var commonPrime *big.Int
	if commonPrime, err = rand.Prime(rand.Reader, 256); err != nil {
		return nil, err
	}

	var encryptedPrime []byte
	var cpBytes []byte = []byte(commonPrime.Text(16))
	if encryptedPrime, err = rsa.EncryptWitPubKey(extKey, cpBytes); err != nil {
		return nil, err
	} else if client.sraKey, err = sra.NewKey(commonPrime, 32); err != nil {
		return nil, err
	}

	client.CommonPrime = commonPrime
	return encryptedPrime, nil
}

// SetEncryptedPrime function receives the common prime encrypted with the
// current client public key, decrypts it with it private key and stores it into
// the current client instance to request the intersection. It also initializes
// the client SRA key with the received and decrypted common prime.
func (client *Client) SetEncryptedPrime(encryptedPrime []byte) (err error) {
	if client.sraKey != nil && client.CommonPrime != nil {
		err = errors.New("common prime already defined, create a new instance")
		return
	}

	var encodedCommonPrime []byte
	encodedCommonPrime, err = client.rsaKey.Decrypt(encryptedPrime)
	if err != nil {
		return
	}

	var sCommonPrime string = string(encodedCommonPrime)
	if commonPrime, ok := new(big.Int).SetString(sCommonPrime, 16); !ok {
		err = errors.New("error decoding decrypted common prime")
	} else if client.sraKey, err = sra.NewKey(commonPrime, 32); err == nil {
		client.CommonPrime = commonPrime
	}

	return
}

// Encrypt function receives the data of the current client to encrypt it with
// the SRA key. It iterates over all items enconding each item to big.Int and
// encrypting it. Then returns the encrypted data.
func (client *Client) Encrypt(data []string) (output [][]*big.Int, err error) {
	if client.sraKey == nil {
		err = errors.New("common prime not defined")
		return
	}

	output = make([][]*big.Int, len(data))
	for i, item := range data {
		var encoded []*big.Int = encoder.StrToInts(item)
		var encrypted []*big.Int = make([]*big.Int, len(encoded))
		for w, word := range encoded {
			encrypted[w] = client.sraKey.Encrypt(word)
		}

		output[i] = encrypted
	}

	return
}

// EncryptExt functions allows to the current client to encrypt the encrypted
// data of another client. It allows to the another client to perform the
// intersection using its re-encrypted data (the output) and, after re-encrypt
// it, the current client encrypted data.
func (client *Client) EncryptExt(input [][]*big.Int) (output [][]*big.Int, err error) {
	if client.sraKey == nil {
		return nil, errors.New("common prime not defined")
	}

	output = make([][]*big.Int, len(input))
	// Iterate over input items and its words encrypting it.
	for i, item := range input {
		var encrypted []*big.Int = make([]*big.Int, len(item))
		for w, word := range item {
			encrypted[w] = client.sraKey.Encrypt(word)
		}
		output[i] = encrypted
	}

	return
}

// PrepareIntersection function receives the current client re-encrypted data
// (from another client) and creates a Bloom Filter with its content to be ready
// to calculate the intersection.
func (client *Client) PrepareIntersection(encryptedData [][]*big.Int) {
	// Initialize the filter.
	client.filter = bloomfilter.NewFilter(len(encryptedData), 0.0001)

	// Iterate over each encrypted data item flatting it into a single slice of
	// bytes with the string representation of all of its words. Then adds the
	// result to the initialized filter.
	for _, item := range encryptedData {
		var record []byte
		for _, word := range item {
			record = append(record, []byte(word.Text(16))...)
		}
		client.filter.Add(record)
	}
}

// GetIntersection function allows to the current client to get the common items
// with the data from another client. It receives the data re-encrypted by the
// the encrypted data of the external client, re-encrypts it with the current
// client SRA key and compares with the its own data using the bloom filter. It
// returns the common data (only encrypted by the client to allow to it to
// decrypt).
func (client *Client) GetIntersection(input [][]*big.Int) ([][]*big.Int, error) {
	if client.filter == nil {
		return nil, errors.New("intersection not initialized")
	}

	var common [][]*big.Int
	for _, item := range input {
		var record []byte
		for _, word := range item {
			var encrypted *big.Int = client.sraKey.Encrypt(word)
			record = append(record, []byte(encrypted.Text(16))...)
		}

		if client.filter.Test(record) {
			common = append(common, item)
		}
	}

	return common, nil
}

// ParseIntersection function decrypts and decodes the received intersection
// result from another client. It returns an error if the common prime is not
// defined or if the decoding process fails.
func (client *Client) ParseIntersection(results [][]*big.Int) ([]string, error) {
	var err error
	if client.sraKey == nil {
		err = errors.New("common prime not defined")
		return nil, err
	}

	// Iterate over intersection result items and its words decrypting and
	// decoding it.
	var output []string = make([]string, len(results))
	for i, item := range results {
		var decrypted []*big.Int = make([]*big.Int, len(item))
		for w, word := range item {
			decrypted[w] = client.sraKey.Decrypt(word)
		}

		var decoded string
		if decoded, err = encoder.IntsToStr(decrypted); err != nil {
			return nil, err
		}

		output[i] = decoded
	}

	return output, nil
}
