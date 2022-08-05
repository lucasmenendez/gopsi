package main

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/lucasmenendez/gopsi/internal/encoder"
	"github.com/lucasmenendez/gopsi/pkg/sra"
)

func main() {
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
	encodedAliceMsg := encoder.StrToInts(aliceMsg)

	// Create and encode Bob secret
	bobMsg := "testemailAddress43@gmail.com"
	encodedBobMsg := encoder.StrToInts(bobMsg)

	// Encrypt Alice original message by Alice first, and then by Bob
	var encryptedAlice []*big.Int
	var encryptedAliceBob []*big.Int
	for _, aliceWord := range encodedAliceMsg {
		encryptedWordAlice := alice.Encrypt(aliceWord)
		encryptedAlice = append(encryptedAlice, encryptedWordAlice)

		encryptedWordAliceBob := bob.Encrypt(encryptedWordAlice)
		encryptedAliceBob = append(encryptedAliceBob, encryptedWordAliceBob)
	}

	// Encrypt Bob original message by Bob, and then by Alice
	var encryptedBob []*big.Int
	var encryptedBobAlice []*big.Int
	for _, bobWord := range encodedBobMsg {
		encryptedWordBob := bob.Encrypt(bobWord)
		encryptedBob = append(encryptedBob, encryptedWordBob)

		encryptedWordBobAlice := alice.Encrypt(encryptedWordBob)
		encryptedBobAlice = append(encryptedBobAlice, encryptedWordBobAlice)
	}

	// Compare partial results
	var lenAlice, lenBob int = len(encryptedAlice), len(encryptedBob)
	var sameLen bool = lenAlice == lenBob
	fmt.Printf("Have both partial encrypted messages same len? %v\n", sameLen)

	var max int = lenAlice
	if !sameLen && lenBob > max {
		max = lenBob
	}

	var areBothEqual bool
	for i := 0; i < max; i++ {
		if !areBothEqual || i >= lenAlice || i >= lenBob {
			break
		}

		areBothEqual = encryptedAlice[i].Cmp(encryptedBob[i]) == 0
	}
	fmt.Printf("Are both partial encrypted messages equal? %v\n", areBothEqual)

	// Compare final results
	var areFinalEqual bool = true
	for i := 0; i < max; i++ {
		if !areFinalEqual || i >= lenAlice || i >= lenBob {
			break
		}

		areFinalEqual = encryptedAliceBob[i].Cmp(encryptedBobAlice[i]) == 0
	}
	fmt.Printf("Are both final encrypted messages equal? %v\n", areFinalEqual)

	// Output: Are both partial encrypted messages equal? false
	// Are both final encrypted messages equal? true
}
