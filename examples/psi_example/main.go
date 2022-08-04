package main

import (
	"fmt"
	"log"
	"math/big"

	"github.com/lucasmenendez/psi/pkg/client"
)

var err error
var alice *client.Client
var bob *client.Client

var aliceData = []string{
	"at.iaculis@google.couk",
	"luctus.et@outlook.couk",
	"sem@aol.edu",
	"donec@outlook.net",
	"nisi@outlook.com",
	"nunc.pulvinar@google.ca",
	"curabitur.dictum@protonmail.edu",
}
var bobData = []string{
	"neque.et@outlook.ca",
	"vehicula.aliquet@yahoo.couk",
	"sem@aol.edu",
	"ut.pellentesque@hotmail.org",
	"non.enim@google.com",
	"justo.praesent@hotmail.couk",
	"nunc.pulvinar@google.ca",
	"amet.consectetuer@hotmail.com",
	"lacinia.sed.congue@aol.com",
	"donec@outlook.net",
}

var encryptedAliceData [][]*big.Int
var encryptedBobData [][]*big.Int

// startIntances initializes the client (alice) and server (bob) instances and
// perform a secure common prime number exchange using RSA (read more here:
// https://github.com/lucasmenendez/gopsi/blob/dev/internal/rsa/rsa.go).
func startIntances() {
	fmt.Println("\nSTARTING SILOS INSTANCES")
	fmt.Println("------------------------")
	// start client instance (alice) to generate public key and share it with
	// the server (bob)
	alice, err = client.Init()
	if err != nil {
		log.Fatalln(err)
	}

	// get client (alice) public key byte slice to share it with the server
	// (bob)
	var alicePubKey []byte
	alicePubKey, err = alice.PubKey()

	// start server instance (bob) with the client public key to generate common
	// prime, encrypt it with the key provided and share the result with the
	// client
	bob, err = client.Init()
	if err != nil {
		log.Fatalln(err)
	}
	var encPrime []byte
	if encPrime, err = bob.GenEncryptedPrime(alicePubKey); err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("[bob] common prime: %s\n", bob.CommonPrime.String())

	// decrypt the common prime on the client side from the encrypted text
	// provided
	err = alice.SetEncryptedPrime(encPrime)
	if err != nil {
		log.Println(err)
	}

	fmt.Printf("[alice] common prime: %s\n", alice.CommonPrime.String())
}

// loadSilosData inject mocked data into client and server instances to get
// encrypted using SRA (read more here:
// https://github.com/lucasmenendez/gopsi/blob/dev/pkg/sra/sra.go)
func loadSilosData(aliceData, bobData []string) {
	fmt.Println("\nLOADING SILOS DATA")
	fmt.Println("------------------")

	// create client (alice) data and load into the client intance to get it
	// encrypted
	encryptedAliceData, err = alice.Encrypt(aliceData)
	if err != nil {
		log.Fatalln("[alice] error encrypting data:", err)
	}
	fmt.Printf("[alice] %d items encrypted. Raw items:\n", len(encryptedAliceData))
	for i, d := range aliceData {
		fmt.Printf("\t%d. %v\n", i, d)
	}

	// create server (bob) data, encrypt and store it into the server intance
	// and share it with the client
	if encryptedBobData, err = bob.Encrypt(bobData); err != nil {
		log.Fatalln("[bob] error encrypting data:", err)
	}
	fmt.Printf("[bob] %d items encrypted. Raw items:\n", len(encryptedBobData))
	for i, d := range bobData {
		fmt.Printf("\t%d. %v\n", i, d)
	}
}

// executeIntersection function performs two actions. First shares the server
// encrypted data to the client to re-encrypt it into the client. Then share
// the re-encrypted server data and encrypted client data with the server.
// Second, creates a BloomFilter (read more here:
// https://github.com/lucasmenendez/gopsi/blob/dev/pkg/bloomfilter/bloomfilter.go)
// with the re-encrypted server data and iterates over client encrypted data,
// re-encrypting it and testing over the created filter.
func executeIntersection() {
	fmt.Println("\nEXECUTING INTERSECTION")
	fmt.Println("----------------------")
	// re-encrypt the server encrypted data into the client and share the
	// result and the encrypted client data with the server
	var encryptedBobDataByAlice [][]*big.Int
	encryptedBobDataByAlice, err = alice.EncryptExt(encryptedBobData)
	if err != nil {
		log.Fatalln(err)
	}

	// initialize the intersection creating a filter with the re-encrypted
	// server data received from the client
	bob.PrepareIntersection(encryptedBobDataByAlice)

	// perform intersection re-encrypyting client data and comparing with
	// the re-encrypted data of the server
	intersection, err := bob.GetIntersection(encryptedAliceData)
	if err != nil {
		log.Fatalln(err)
	}

	// Parse received results from the server on the client.
	var results []string
	results, err = alice.ParseIntersection(intersection)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("[alice] %d common items received:\n", len(intersection))
	for _, d := range results {
		var index int
		for i, c := range aliceData {
			index = i
			if c == d {
				break
			}
		}
		fmt.Printf("\t%d. %v\n", index, d)
	}
}

func main() {
	// request intersection
	startIntances()

	// perform intersection
	loadSilosData(aliceData, bobData)
	executeIntersection()
}
