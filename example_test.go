package psi

import (
	"crypto/rand"
	"fmt"

	"github.com/lucasmenendez/psi/internal/encoder"
	"github.com/lucasmenendez/psi/pkg/sra"
)

func Example() {
	var err error

	// Agree a prime seed
	prime, _ := rand.Prime(rand.Reader, 256)

	// Create Alice key pair
	var alice *sra.SRAKey
	if alice, err = sra.NewKey(prime, 32); err != nil {
		fmt.Println(err)
	}

	// Create Bob key pair
	var bob *sra.SRAKey
	if bob, err = sra.NewKey(prime, 32); err != nil {
		fmt.Println(err)
	}

	// Create and encode Alice secret
	aliceMsg := "testemailAddress43@gmail.com"
	encodedAliceMsg := encoder.StrToInt(aliceMsg)

	// Create and encode Bob secret
	bobMsg := "testemailAddress43@gmail.com"
	encodedBobMsg := encoder.StrToInt(bobMsg)

	// Encrypt Alice original message by Alice first, and then by Bob
	encryptedAlice := alice.Encrypt(encodedAliceMsg)
	encryptedAliceBob := bob.Encrypt(encryptedAlice)

	// Encrypt Bob original message by Bob, and then by Alice
	encryptedBob := bob.Encrypt(encodedBobMsg)
	encryptedBobAlice := alice.Encrypt(encryptedBob)

	// Compare partial results
	arePartialEqual := encryptedAlice.Cmp(encryptedBob) == 0
	fmt.Printf("Are both partial encrypted messages equal? %v\n", arePartialEqual)

	// Compare final results
	areFinalEqual := encryptedAliceBob.Cmp(encryptedBobAlice) == 0
	fmt.Printf("Are both final encrypted messages equal? %v\n", areFinalEqual)

	// Output: Are both partial encrypted messages equal? false
	// Are both final encrypted messages equal? true
}
