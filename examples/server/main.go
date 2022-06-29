package main

import "github.com/lucasmenendez/psi/internal/server"

func main() {
	server := server.Init(8000)
	server.Start()
}
