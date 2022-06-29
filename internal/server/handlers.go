package server

import (
	"crypto/rand"
	"io/ioutil"
	"net/http"

	"github.com/lucasmenendez/psi/internal/rsa"
)

func handleIntersectionRequest(w http.ResponseWriter, r *http.Request) {
	key, err := ioutil.ReadAll(r.Body)
	if reqParseErr(w, err) {
		return
	}

	prime, _ := rand.Prime(rand.Reader, 256)
	res, err := rsa.EncryptWitPublicKey(key, []byte(prime.Text(16)))
	if encryptionErr(w, err) {
		return
	}

	n, err := w.Write(res)
	if resEncodeErr(w, err) {
		return
	}

	if n != len(res) {
		internalErr(w, "response and encrypted prime have not the same lenght")
		return
	}
}
