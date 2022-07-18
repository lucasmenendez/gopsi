package main

import (
	"log"
	"math/big"

	"github.com/lucasmenendez/psi/pkg/client"
	"github.com/lucasmenendez/psi/pkg/server"
)

var err error
var alice *client.Client
var bob *server.Server

func startIntances() {
	// start client instance (alice) to generate public key and share it with
	// the server
	alice, err = client.Init()
	if err != nil {
		log.Fatalln(err)
	}

	// start server instance (bob) with the client public key to generate common
	// prime, encrypt it with the key provided and share the result with the
	// client
	bob, err = server.Init(alice.PublicKey)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("generated common prime by server: %s\n", bob.CommonPrime.String())

	// decrypt the common prime on the client side from the encrypted text
	// provided
	err = alice.AddCommonPrime(bob.CommonPrimeEncrypted)
	if err != nil {
		log.Println(err)
	}
	log.Printf("decrypted common prime by client: %s\n", alice.CommonPrime.String())
}

func loadSilosData() {
	// create client (alice) data and load into the client intance to get it
	// encrypted
	var aliceData = []string{
		"Donec molestie justo eget leo convallis ullamcorper.",
		"Praesent ornare feugiat ultrices.",
		"Morbi est nisi, volutpat pellentesque eros id, lacinia tempus ante.",
	}
	err = alice.LoadData(aliceData)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("loaded data by client: %v", alice.Records)

	// create server (bob) data, encrypt and store it into the server intance
	// and share it with the client
	var bobData = []string{
		"Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
		"Donec molestie justo eget leo convallis ullamcorper.",
		"Nam eros enim, dapibus euismod sodales eget, condimentum id enim.",
		"Praesent ornare feugiat ultrices.",
		"Donec tortor velit, ornare a interdum at, viverra et urna.",
	}
	bob.LoadData(bobData)
	log.Printf("loaded data by server: %v", bob.Records)
}

func executeIntersection() {
	// re-encrypt the server encrypted records into the client and share the
	// result and the encrypted client data with the server
	var encryptedBobData []*big.Int
	encryptedBobData, err = alice.EncryptInput(bob.Records)

	// perform the intersection creating a filter with the re-encrypted server
	// data received from the client and test it with the encrypted client
	// records, re-encrypting it with the server first.
	intersection := bob.GetIntersection(encryptedBobData, alice.Records)
	log.Printf("alice items that are stored by bob: %v", intersection)
}

func main() {
	// startIntances initializes the client (alice) and server (bob) instances
	// and perform a secure common prime number exchange using RSA (read more
	// here: https://github.com/lucasmenendez/gopsi/blob/dev/internal/rsa/rsa.go).
	startIntances()

	// loadSilosData inject mocked data into client and server instances to get
	// encrypted using SRA (read more here:
	// https://github.com/lucasmenendez/gopsi/blob/dev/pkg/sra/sra.go)
	loadSilosData()

	// executeIntersection function performs two actions. First shares the
	// server encrypted data to the client to re-encrypt it into the client.
	// Then share the re-encrypted server data and encrypted client data with
	// the server.
	// Second, creates a BloomFilter (read more here:
	// https://github.com/lucasmenendez/gopsi/blob/dev/pkg/bloomfilter/bloomfilter.go)
	// with the re-encrypted server data and iterates over client encrypted
	// data, re-encrypting it and testing over the created filter.
	executeIntersection()
}
