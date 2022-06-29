package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"

	"github.com/lucasmenendez/psi/rsa"
)

func requestIntersection(w http.ResponseWriter, r *http.Request) {
	// Generate RSA keys pair
	keys, err := rsa.NewKey(1024)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error during RSA keys generation.", 500)
		return
	}

	// Get public key
	var pk []byte
	if pk, err = keys.PubKey(); err != nil {
		log.Println(err)
		http.Error(w, "Error during RSA pub key encoding.", 500)
		return
	}

	// Request new intersection sending the public key
	body := bytes.NewBuffer(pk)
	res, err := http.Post("http://localhost:8080/newIntersection", "application/octet-stream", body)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error during instersection request.", 500)
		return
	}

	// Read encrypted common prime from request response
	defer res.Body.Close()
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error during intersection response.", 500)
		return
	}

	// Decrypt common prime with private key
	encPrime, err := keys.Decrypt(content)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error during common prime decryption.", 500)
		return
	}

	// Deconding prime integer
	var prime, _ = new(big.Int).SetString(string(encPrime), 16)
	w.Write([]byte(prime.String()))
}

func main() {
	http.HandleFunc("/requestIntersection", requestIntersection)

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Client listening on port 8000")
}
