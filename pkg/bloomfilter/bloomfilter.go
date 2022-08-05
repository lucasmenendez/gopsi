package bloomfilter

import (
	"hash"
	"hash/fnv"
	"math"
)

// BloomFilter struct contains the required parameters to create and use a
// filter such as the data bitmap (data), the optimal number of bits (m), the
// optimal number of hash functions (k) and the size of the filter (n). It also
// contains an instance of the hash function to better performance. Read more
// about Bloom Filter definition and implementation here:
// https://en.wikipedia.org/wiki/Bloom_filter.
type BloomFilter struct {
	data []bool      // filter content
	m    uint        // number of bits of the filter
	k    uint        // number of hashing functions
	n    uint        // number of items into the filter
	hash hash.Hash64 // hash function seed
}

// NewFilter functions initializes a new BloomFilter with the size and false
// positive rate provided as argument. Using this arguments, it calculates the
// optimal number of bits for the number of items to store, and optimal number
// of hash functions for this size.
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
	// the calculate number of bits (m).
	filter.data = make([]bool, filter.m)

	// Calculate the optimal number of hashes.
	filter.k = filter.numberOfHashes()
	return
}

// calcHash function generates a splitted 64-bits hash representation of the
// byte array provided as input. The hash is splitted to allow to create k
// hashes according to Kirsch-Mitzenmacher optimization, instead of create k
// single hashes.
func (f *BloomFilter) calcHash(input []byte) (uint, uint) {
	defer f.hash.Reset()

	// Create hash of 64 bits from the current item
	f.hash.Write(input)
	var hashed uint64 = f.hash.Sum64()

	// Split the hashed
	var a, b uint32 = uint32(hashed >> 32), uint32(hashed)

	return uint(a), uint(b)
}

// numberOfBits function calculates the optimal number of bits to store the
// current size of the filter (n: number of items) with the provided false
// positive rate (fp).
func (f *BloomFilter) numberOfBits(fp float64) uint {
	// Calculate the number of bits (m) of the filter by the it size (n) and the
	// false positive rate provided as argument according to the following
	// formula: m = -1 * (n * ln(fp)) / ln(2)^2
	return uint(-1 * float64(f.n) * math.Logb(fp) / math.Pow(math.Logb(2), 2))
}

// numberOfHashes function calculates the optimal number of hashes for the
// current filter size (n: number of items) and the current number of bits (m).
func (f *BloomFilter) numberOfHashes() uint {
	// Calculate the number of hash functions (k) of the filter by the number of
	// bits (m) and the size of the filter (n), according to the following
	// formula: k = (m / n) * ln(2)
	return uint(math.Ceil(math.Logb(2) * float64(f.m) / float64(f.n)))
}

// Add function allows to user to insert one (or more) items to the created
// filter. It calculates the position of each input with the number of hashes to
// mark as contained.
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

// Test function allows to user to check if the current filter has already an
// item provided as input. It performs almost the same action as Add function,
// but checking if each hashed position is already marked. If some of the
// positions are not marked, the filter does not contains the provided item.
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

// TestMultiple function allows to the user to test multiple inputs at the same
// time, using the BloomFilter.Test function.
func (f *BloomFilter) TestMultiple(items ...[]byte) (results []bool) {
	// Create bool slice to store every single test result.
	results = make([]bool, len(items))

	// Check every item provided with BloomFilter.Test function and store the result.
	for i, item := range items {
		results[i] = f.Test(item)
	}
	return
}
