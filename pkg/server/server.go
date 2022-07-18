package server

import (
	"crypto/rand"
	"errors"
	"math/big"

	"github.com/lucasmenendez/psi/internal/encoder"
	"github.com/lucasmenendez/psi/internal/rsa"
	"github.com/lucasmenendez/psi/pkg/bloomfilter"
	"github.com/lucasmenendez/psi/pkg/sra"
)

type Server struct {
	CommonPrime          *big.Int
	CommonPrimeEncrypted []byte
	sraKey               *sra.SRAKey
	Records              []*big.Int
	filter               *bloomfilter.BloomFilter
}

func Init(clientKey []byte) (server *Server, err error) {
	server = &Server{}
	if server.CommonPrime, err = rand.Prime(rand.Reader, 256); err != nil {
		return
	}

	server.sraKey, err = sra.NewKey(server.CommonPrime, 32)
	if err != nil {
		return
	}

	var prime []byte = []byte(server.CommonPrime.Text(16))
	server.CommonPrimeEncrypted, err = rsa.EncryptWitPublicKey(clientKey, prime)
	return
}

func (server *Server) LoadData(data []string) error {
	if server.sraKey == nil {
		return errors.New("common prime not defined")
	}

	server.Records = make([]*big.Int, len(data))
	for i, item := range data {
		encoded := encoder.StrToInt(item)
		server.Records[i] = server.sraKey.Encrypt(encoded)
	}

	return nil
}

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
