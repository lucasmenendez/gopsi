package client

import (
	"math/big"
	"reflect"
	"testing"
)

func TestInit(t *testing.T) {
	if _, err := Init(); err != nil {
		t.Fatalf("expected nil, got %s", err)
	}
}

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

	if _, err := clientB.GenEncryptedPrime(nil); err == nil {
		t.Fatal("expected error, got nil")
	} else if _, err = clientB.GenEncryptedPrime([]byte{}); err == nil {
		t.Fatal("expected error, got nil")
	}

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
	var err error
	var input = []string{"hello world"}
	clientA, _ := Init()
	clientB, _ := Init()

	var (
		nilData    []string = nil
		emptyData  []string = make([]string, 0)
		notSRAData []string = []string{""}
	)
	if _, err = clientA.Encrypt(nilData); err == nil {
		t.Fatal("expected error got nil")
	} else if _, err = clientA.Encrypt(emptyData); err == nil {
		t.Fatal("expected error got nil")
	} else if _, err = clientA.Encrypt(notSRAData); err == nil {
		t.Fatal("expected error got nil")
	}

	pubKey, _ := clientB.PubKey()
	encPrime, _ := clientA.GenEncryptedPrime(pubKey)
	clientB.SetEncryptedPrime(encPrime)

	if outputA, err := clientA.Encrypt(input); err != nil {
		t.Fatalf("expected nil, got %s", err)
	} else if result, _ := clientA.ParseIntersection(outputA); result[0] != input[0] {
		t.Fatalf("expected %v, got %v", input, result)
	} else if outputB, _ := clientB.Encrypt(input); reflect.DeepEqual(outputA, outputB) {
		t.Fatal("expected both encrypted messages different, got same")
	}
}

func TestEncryptExt(t *testing.T) {
	var err error
	var input = []string{"hello world"}

	clientA, _ := Init()
	clientB, _ := Init()

	var (
		nilData    [][]*big.Int = nil
		emptyData  [][]*big.Int = make([][]*big.Int, 0)
		notSRAData [][]*big.Int = [][]*big.Int{{new(big.Int).SetInt64(0)}}
	)
	if _, err = clientA.EncryptExt(nilData); err == nil {
		t.Fatal("expected error got nil")
	} else if _, err = clientA.EncryptExt(emptyData); err == nil {
		t.Fatal("expected error got nil")
	} else if _, err = clientA.EncryptExt(notSRAData); err == nil {
		t.Fatal("expected error got nil")
	}

	pubKey, _ := clientB.PubKey()
	encPrime, _ := clientA.GenEncryptedPrime(pubKey)
	clientB.SetEncryptedPrime(encPrime)
	encInputByA, _ := clientA.Encrypt(input)

	var encInputByAB [][]*big.Int
	if encInputByAB, err = clientB.EncryptExt(encInputByA); err != nil {
		t.Fatalf("expected nil, got %s", err)
	} else if len(encInputByA) != len(encInputByAB) {
		t.Fatalf("expected len %d, got len %d", len(encInputByA), len(encInputByAB))
	}

	var result [][]*big.Int = make([][]*big.Int, len(encInputByA))
	for i, item := range encInputByAB {
		var decrypted []*big.Int = make([]*big.Int, len(item))

		for w, word := range item {
			decrypted[w] = clientB.sraKey.Decrypt(word)
		}
		result[i] = decrypted
	}

	if !reflect.DeepEqual(encInputByA, result) {
		t.Fatalf("expected %v, got %v", encInputByA, result)
	}
}

func TestInitIntersection(t *testing.T) {
	var err error
	var input = []string{"hello world"}

	clientA, _ := Init()
	clientB, _ := Init()

	pubKey, _ := clientB.PubKey()
	encPrime, _ := clientA.GenEncryptedPrime(pubKey)
	clientB.SetEncryptedPrime(encPrime)

	if err = clientA.PrepareIntersection(nil); err == nil {
		t.Fatal("expected error got nil")
	} else if err = clientA.PrepareIntersection([][]*big.Int{}); err == nil {
		t.Fatal("expected error got nil")
	}

	encInputByA, _ := clientA.Encrypt(input)
	encInputByAB, _ := clientB.EncryptExt(encInputByA)

	if err = clientA.PrepareIntersection(encInputByAB); err != nil {
		t.Fatalf("expected nil, got %s", err)
	}
}

func TestGetIntersection(t *testing.T) {
	var err error
	var input = []string{"hello world"}

	clientA, _ := Init()
	clientB, _ := Init()

	var (
		nilData       [][]*big.Int = nil
		emptyData     [][]*big.Int = make([][]*big.Int, 0)
		noInitialized [][]*big.Int = [][]*big.Int{{new(big.Int).SetInt64(0)}}
	)
	if _, err = clientA.GetIntersection(nilData); err == nil {
		t.Fatal("expected error got nil")
	} else if _, err = clientA.GetIntersection(emptyData); err == nil {
		t.Fatal("expected error got nil")
	} else if _, err = clientA.GetIntersection(noInitialized); err == nil {
		t.Fatal("expected error got nil")
	}

	pubKey, _ := clientB.PubKey()
	encPrime, _ := clientA.GenEncryptedPrime(pubKey)
	clientB.SetEncryptedPrime(encPrime)

	if _, err = clientA.GetIntersection(noInitialized); err == nil {
		t.Fatal("expected error got nil")
	}

	encInputByA, _ := clientA.Encrypt(input)
	encInputByB, _ := clientB.Encrypt(input)

	encInputByAB, _ := clientB.EncryptExt(encInputByA)
	clientA.PrepareIntersection(encInputByAB)

	if result, err := clientA.GetIntersection(encInputByB); err != nil {
		t.Fatalf("expected nil, got %s", err)
	} else if !reflect.DeepEqual(encInputByB, result) {
		t.Fatalf("expected %v, got %v", encInputByB, result)
	}
}

func TestParseIntersection(t *testing.T) {
	var err error
	var input = []string{"hello world"}

	clientA, _ := Init()
	clientB, _ := Init()

	var (
		nilData       [][]*big.Int = nil
		emptyData     [][]*big.Int = make([][]*big.Int, 0)
		noInitialized [][]*big.Int = [][]*big.Int{{new(big.Int).SetInt64(0)}}
	)
	if _, err = clientB.ParseIntersection(nilData); err == nil {
		t.Fatal("expected error got nil")
	} else if _, err = clientB.ParseIntersection(emptyData); err == nil {
		t.Fatal("expected error got nil")
	} else if _, err = clientB.ParseIntersection(noInitialized); err == nil {
		t.Fatal("expected error got nil")
	}

	pubKey, _ := clientB.PubKey()
	encPrime, _ := clientA.GenEncryptedPrime(pubKey)
	clientB.SetEncryptedPrime(encPrime)

	encInputByA, _ := clientA.Encrypt(input)
	encInputByB, _ := clientB.Encrypt(input)

	encInputByAB, _ := clientB.EncryptExt(encInputByA)
	clientA.PrepareIntersection(encInputByAB)

	result, _ := clientA.GetIntersection(encInputByB)
	if output, err := clientB.ParseIntersection(result); err != nil {
		t.Fatalf("expected nil, got %s", err)
	} else if !reflect.DeepEqual(input, output) {
		t.Fatalf("expected %v, got %v", input, output)
	}
}
