[![GoDoc](https://godoc.org/github.com/lucasmenendez/gopsi?status.svg)](https://godoc.org/github.com/lucasmenendez/gopsi)
[![Go Report Card](https://goreportcard.com/badge/github.com/lucasmenendez/gopsi)](https://goreportcard.com/report/github.com/lucasmenendez/gopsi)
[![production](https://github.com/lucasmenendez/gopsi/workflows/production/badge.svg)](https://github.com/lucasmenendez/gopsi/actions?query=workflow%3Aproduction)
[![license](https://img.shields.io/github/license/lucasmenendez/gopsi)](LICENSE)

# GoPSI - Private Set Intersection in Golang

Simple Private Set Intersection implemented in pure Go. It uses SRA algorithm [[1]](#references) as encryption scheme and Bloom Filters [[2]](#references) to perform set intersection.

## Basic usage protocol

The following diagram explains the basic example of the library using [Client](./pkg/client/client.go) (alice) and [Server](./pkg/server/server.go) (bob) structs. This example is implemented into [psi_example/main.go](./examples/psi_example/main.go) file.

However you can use [SRA](./pkg/sra/sra.go) or [BloomFilters](./pkg/bloomfilter/bloomfilter.go) isolated and also design your own protocol using it.

```mermaid
sequenceDiagram
    participant alice
    participant bob

    Note over alice,bob: Request intersection
    alice->>+bob: send alice RSA public key
    bob-->>bob: generate prime number 
    bob-->>-alice: send prime number encrypted w/ alice public key

    Note over alice,bob: perform intersection
    par both encrypts its own data with the common prime
        alice-->>alice: encrypt its data with SRA
    and
        bob-->>bob: encrypt its data with SRA
    end

    bob->>alice: send encrypted data
    alice-->>alice: re-encrypt bob data with SRA
    alice->>bob: send its encrypted data and bobs re-encrypted data
    bob-->>bob: re-encrypt alice data with SRA
    bob-->>alice: send the intersection between re-encrypted data sets
```

## Docs & example
Checkout [GoDoc Documentation](https://godoc.org/github.com/lucasmenendez/gopsi)


## References

1. Adi Shamir, Ronald L. Rivest and Leonard M. Adleman, *"Mental Poker"*, April 1979. https://people.csail.mit.edu/rivest/pubs/SRA81.pdf
2. Wikipedia, *"Bloom filter"*, July 2005. https://en.wikipedia.org/wiki/Bloom_filter