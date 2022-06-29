package server

import (
	"log"
	"net/http"
)

func reqParseErr(w http.ResponseWriter, err error) bool {
	if err != nil {
		log.Println(err)
		http.Error(w, "request parsing error", http.StatusBadRequest)
		return true
	}

	return false
}

func encryptionErr(w http.ResponseWriter, err error) bool {
	if err != nil {
		log.Println(err)
		http.Error(w, "encryption error", http.StatusInternalServerError)
		return true
	}

	return false
}

func resEncodeErr(w http.ResponseWriter, err error) bool {
	if err != nil {
		log.Println(err)
		http.Error(w, "response encoding error", http.StatusInternalServerError)
		return true
	}

	return false
}

func internalErr(w http.ResponseWriter, err string) {
	log.Println(err)
	http.Error(w, err, http.StatusInternalServerError)
}
