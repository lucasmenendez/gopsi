package sdk

import (
	"crypto/rand"
	"errors"
	"math/big"

	"github.com/lucasmenendez/psi/internal/encoder"
	"github.com/lucasmenendez/psi/internal/rsa"
	"github.com/lucasmenendez/psi/pkg/bloomfilter"
	"github.com/lucasmenendez/psi/pkg/sra"
)

// Server struct contains all required parameters to allows to a client to
// request a private set intersection. It generates a common prime number
// (required by SRA protocol), allows to share it securely using RSA public key
// from the client and perform the intersection.
type Server struct {
	CommonPrime *big.Int
	sraKey      *sra.SRAKey
	Data        [][]*big.Int
	filter      *bloomfilter.BloomFilter
}

// Init function instances a Server generating a common prime number and
// initializing the server SRA key with it.
func InitServer() (server *Server, err error) {
	server = &Server{}
	if server.CommonPrime, err = rand.Prime(rand.Reader, 256); err != nil {
		return
	}

	server.sraKey, err = sra.NewKey(server.CommonPrime, 32)
	if err != nil {
		return
	}

	return
}

// EncryptedPrime function encrypts the generated common prime with the RSA
// public key provided by the client.
func (server *Server) EncryptedPrime(clientKey []byte) ([]byte, error) {
	var prime []byte = []byte(server.CommonPrime.Text(16))
	return rsa.EncryptWitPublicKey(clientKey, prime)
}

// LoadData function receives the data to request the intersection. It
// iterates over all items encondign each item to big.Int and encrypting it with
// SRA. Then stores the encrypted data into the current client instance.
func (server *Server) LoadData(data []string) error {
	server.Data = make([][]*big.Int, len(data))
	for i, item := range data {
		var encrypted []*big.Int

		var encoded []*big.Int = encoder.StrToInts(item)
		for _, word := range encoded {
			var encryptedWord = server.sraKey.Encrypt(word)
			encrypted = append(encrypted, encryptedWord)
		}

		server.Data[i] = encrypted
	}

	return nil
}

// InitIntersection function receives the re-encrypted server data from the
// client and creates a Bloom Filter with its content to be ready to calculate
// the intersection.
func (server *Server) InitIntersection(encryptedData [][]*big.Int) {
	// Initialize the filter.
	server.filter = bloomfilter.NewFilter(len(encryptedData), 0.001)

	// Iterate over each encrypted data item flatting it into a single slice of
	// bytes with the string representation of all of its words. Then adds the
	// result to the initialized filter.
	for _, item := range encryptedData {
		var record []byte
		for _, word := range item {
			record = append(record, []byte(word.Text(16))...)
		}
		server.filter.Add(record)
	}
}

// GetIntersection function allows to the server to get the common items
// with the client data. It receives the server data re-encrypted by the
// client and the encrypted client data. First, re-encrypt the client data
// and then, compares with the its own data, re-encrypted by the client. It
// returns the common data (only encrypted by the client to allow to it to
// decrypt).
func (server *Server) GetIntersection(input [][]*big.Int) ([][]*big.Int, error) {
	if server.filter == nil {
		return nil, errors.New("intersection not initialized")
	}

	var common [][]*big.Int
	for _, item := range input {
		var record []byte
		for _, word := range item {
			var encrypted *big.Int = server.sraKey.Encrypt(word)
			record = append(record, []byte(encrypted.Text(16))...)
		}

		if server.filter.Test(record) {
			common = append(common, item)
		}
	}

	return common, nil
}
