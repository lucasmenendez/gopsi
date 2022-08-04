package client

import "testing"

func TestPubKey(t *testing.T) {
	client := &Client{}
	if _, err := client.PubKey(); err == nil {
		t.Error("expected error, got nil")
	}

	client, _ = Init()
	if _, err := client.PubKey(); err != nil {
		t.Errorf("expected nil, got '%s'", err)
	}
}

func TestGenEncryptedPrime(t *testing.T) {
	clientA, _ := Init()
	clientB, _ := Init()

	aPubKey, _ := clientA.PubKey()
	encPrime, err := clientB.GenEncryptedPrime(aPubKey)
	if err != nil {
		t.Fatalf("expected nil, got %s", err)
	} else if clientB.CommonPrime == nil {
		t.Fatal("expected not nil, got nil")
	}

	result, _ := clientA.rsaKey.Decrypt(encPrime)
	if expected := clientB.CommonPrime.Text(16); string(result) != expected {
		t.Fatalf("expected %s, got %s", expected, result)
	}

	if _, err := clientB.GenEncryptedPrime(aPubKey); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestSetEncryptedPrime(t *testing.T) {
	clientA, _ := Init()
	clientB, _ := Init()

	aPubKey, _ := clientA.PubKey()
	encPrime, _ := clientB.GenEncryptedPrime(aPubKey)

	if err := clientA.SetEncryptedPrime(encPrime); err != nil {
		t.Fatalf("expected nil, got %s", err)
	} else if clientA.CommonPrime.Cmp(clientB.CommonPrime) != 0 {
		t.Fatalf("expected %d, got %d", clientA.CommonPrime.Int64(), clientB.CommonPrime.Int64())
	} else if err := clientA.SetEncryptedPrime(encPrime); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestEncrypt(t *testing.T) {

}

func TestEncryptExt(t *testing.T) {

}

func TestInitIntersection(t *testing.T) {

}

func TestGetIntersection(t *testing.T) {

}

func TestParse(t *testing.T) {

}
