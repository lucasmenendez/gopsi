package server

import (
	"crypto/rand"
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
	Records     []*big.Int
	filter      *bloomfilter.BloomFilter
}

// Init function instances a Server generating a common prime number and
// initializing the server SRA key with it.
func Init() (server *Server, err error) {
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

// LoadData function receives the records to request the intersection. It
// iterates over all items encondign each item to big.Int and encrypting it with
// SRA. Then stores the encrypted records into the current client instance.
func (server *Server) LoadData(data []string) error {
	server.Records = make([]*big.Int, len(data))
	for i, item := range data {
		encoded := encoder.StrToInt(item)
		server.Records[i] = server.sraKey.Encrypt(encoded)
	}

	return nil
}

// GetIntersection function allows to the server to get the common items
// with the client records. It receives the server records re-encrypted by the
// client and the encrypted client records. First, re-encrypt the client records
// and then, compares with the its own records, re-encrypted by the client. It
// returns the common records (only encrypted by the client to allow to it to
// decrypt).
func (server *Server) GetIntersection(encRecords, input []*big.Int) []*big.Int {
	server.filter = bloomfilter.NewFilter(len(encRecords), 0.001)
	for _, record := range encRecords {
		server.filter.Add([]byte(record.Text(16)))
	}

	var common []*big.Int
	for _, item := range input {
		encrypted := server.sraKey.Encrypt(item)
		encoded := []byte(encrypted.Text(16))
		if server.filter.Test(encoded) {
			common = append(common, item)
		}
	}

	return common
}
