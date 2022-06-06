package psi

import (
	"hash"
	"hash/fnv"
	"math"
)

type BloomFilter struct {
	data []bool      // filter content
	m    uint        // number of bits of the filter
	k    uint        // number of hashing functions
	n    uint        // number of items into the filter
	hash hash.Hash64 // hash function seed
}

func NewFilter(size int, fp float64) (filter *BloomFilter) {
	// Initializes the Bloom Filter with the size provided and a initialized
	// hash function seed.
	filter = &BloomFilter{
		n:    uint(size),
		hash: fnv.New64a(),
	}

	// Calculate the required number of bits of the filter by the data size and
	// the false positive rate provided.
	filter.m = filter.numberOfBits(fp)
	// Initialize the filter data bit map creating a slice of bool to store
	// the calulate number of bits (m).
	filter.data = make([]bool, filter.m)

	// Calculate the optimal number of hashes.
	filter.k = filter.numberOfHashes()
	return
}

func (f *BloomFilter) calcHash(input []byte) (uint, uint) {
	defer f.hash.Reset()

	// Create hash of 64 bits from the current item
	f.hash.Write(input)
	var hashed uint64 = f.hash.Sum64()

	// Split the hashed
	var a, b uint32 = uint32(hashed >> 32), uint32(hashed)

	return uint(a), uint(b)
}

func (f *BloomFilter) numberOfBits(fp float64) uint {
	// Calculate the number of bits (m) of the filter by the it size (n) and the
	// false positive rate provided as argument according to the following
	// formula: m = -1 * (n * ln(fp)) / ln(2)^2
	return uint(-1 * float64(f.n) * math.Logb(fp) / math.Pow(math.Logb(2), 2))
}

func (f *BloomFilter) numberOfHashes() uint {
	// Calculate the number of hash functions (k) of the filter by the number of
	// bits (m) and the size of the filter (n), according to the following
	// formula: k = (m / n) * ln(2)
	return uint(math.Ceil(math.Logb(2) * float64(f.m) / float64(f.n)))
}

func (f *BloomFilter) Add(items ...[]byte) {
	for _, item := range items {
		// For each item provided, calculate both hash parts to generate k hash
		// functions according to the Kirsch-Mitzenmacher optimization
		// (https://www.eecs.harvard.edu/~michaelm/postscripts/tr-02-05.pdf).
		var a, b uint = f.calcHash(item)

		// Set to 1 (true) every bit map position calculates with the hash parts
		// generated.
		for i := uint(0); i < f.k; i++ {
			var index uint = (a + i*b) % f.m
			f.data[index] = true
		}
	}
}

func (f *BloomFilter) Test(item []byte) bool {
	// Calculate both hash parts for the item provided to generate k hash
	// functions and get its byte positions.
	var a, b uint = f.calcHash(item)

	// Calculate every item byte position, if some of byte position are false
	// (0) into the bitmap (f.data), it does not contains the item. If every
	// byte position are true (1), the bitma probably contains the item.
	for i := uint(0); i < f.k; i++ {
		var index uint = (a + b*i) % f.m
		if !f.data[index] {
			return false
		}
	}

	return true
}

func (f *BloomFilter) TestMultiple(items ...[]byte) (results []bool) {
	// Create bool slice to store every single test result.
	results = make([]bool, len(items))

	// Check every item provided with f.Test function and store the result.
	for i, item := range items {
		results[i] = f.Test(item)
	}
	return
}
