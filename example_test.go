package psi_test

import (
	"crypto/rand"
	"fmt"

	"github.com/lucasmenendez/psi"
)

func Example() {
	var err error

	// Agree a prime seed
	prime, _ := rand.Prime(rand.Reader, 256)

	// Create Alice key pair
	var alice *psi.KeyPair
	if alice, err = psi.GenerateKeyPair(rand.Reader, prime, 32); err != nil {
		fmt.Println(err)
	}

	// Create Bob key pair
	var bob *psi.KeyPair
	if bob, err = psi.GenerateKeyPair(rand.Reader, prime, 32); err != nil {
		fmt.Println(err)
	}

	// Create and encode Alice secret
	aliceMsg := "testemailAddress43@gmail.com"
	encodedAliceMsg := psi.Encode(aliceMsg)

	// Create and encode Bob secret
	bobMsg := "testemailAddress43@gmail.com"
	encodedBobMsg := psi.Encode(bobMsg)

	// Encrypt Alice original message by Alice first, and then by Bob
	encryptedAlice := alice.EncryptInt(encodedAliceMsg)
	encryptedAliceBob := bob.EncryptInt(encryptedAlice)

	// Encrypt Bob original message by Bob, and then by Alice
	encryptedBob := bob.EncryptInt(encodedBobMsg)
	encryptedBobAlice := alice.EncryptInt(encryptedBob)

	// Compare partial results
	arePartialEqual := encryptedAlice.Cmp(encryptedBob) == 0
	fmt.Printf("Are both partial encrypted messages equal? %v\n", arePartialEqual)

	// Compare final results
	areFinalEqual := encryptedAliceBob.Cmp(encryptedBobAlice) == 0
	fmt.Printf("Are both final encrypted messages equal? %v\n", areFinalEqual)

	// Output: Are both partial encrypted messages equal? false
	// Are both final encrypted messages equal? true
}
