[![GoDoc](https://godoc.org/github.com/lucasmenendez/gopsi?status.svg)](https://godoc.org/github.com/lucasmenendez/gopsi)
[![Go Report Card](https://goreportcard.com/badge/github.com/lucasmenendez/gopsi)](https://goreportcard.com/report/github.com/lucasmenendez/gopsi)
[![test](https://github.com/lucasmenendez/gopsi/workflows/test/badge.svg)](https://github.com/lucasmenendez/gopsi/actions?query=workflow%3Atest)
[![license](https://img.shields.io/github/license/lucasmenendez/gopsi)](LICENSE)

# GoPSI - Private Set Intersection in Golang

Simple Private Set Intersection implemented in pure Go. It uses SRA algorithm [[1]](#references) as encryption scheme and Bloom Filters [[2]](#references) to perform set intersection.

## Examples and Docs
Two full examples are already implemented:
- Simple SRA encryption: [code](examples/sra_example/main.go)
- PSI algorithm example: [docs](examples/psi_example/README.md) [code](examples/psi_example/main.go).

Checkout [GoDoc Documentation](https://godoc.org/github.com/lucasmenendez/gopsi)


## References

1. Adi Shamir, Ronald L. Rivest and Leonard M. Adleman, *"Mental Poker"*, April 1979. https://people.csail.mit.edu/rivest/pubs/SRA81.pdf
2. Wikipedia, *"Bloom filter"*, July 2005. https://en.wikipedia.org/wiki/Bloom_filter