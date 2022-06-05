package psi

import (
	"hash"
	"hash/fnv"
	"math"
)

type BloomFilter struct {
	data []bool // filter content
	m uint // number of bits of the filter
	k uint // number of hashing functions
	n uint // number of items into the filter
	hash hash.Hash64 // hash functions seed
}

func NewFilter(size int, fp float64) (filter *BloomFilter) {
	filter = &BloomFilter{
		n: uint(size),
		hash: fnv.New64a(),
	}

	filter.m = filter.numberOfBits(fp)
	filter.k = filter.numberOfHashes()
	filter.data = make([]bool, filter.m)
	
	return 
}

func (f *BloomFilter) hashPair(input []byte) (uint, uint) {
	defer f.hash.Reset()

	f.hash.Write(input)
	var hashed uint64 = f.hash.Sum64()
	var a, b uint32 = uint32(hashed >> 32), uint32(hashed)
	
	return uint(a), uint(b)
}

func (f *BloomFilter) numberOfHashes() uint {
	return uint(math.Ceil(math.Logb(2) * float64(f.m) / float64(f.n)))
}

func (f *BloomFilter) numberOfBits(fp float64) uint {
	return uint(-1 * float64(f.n) * math.Logb(fp) / math.Pow(math.Logb(2), 2))
}

func (f *BloomFilter) Add(item []byte) {
	var a, b uint = f.hashPair(item)
	for i := uint(0); i < f.k; i++ {
		var index uint = (a + b * i) % f.m
		f.data[index] = true
	}
}

func (f *BloomFilter) Test(item []byte) bool {
	var a, b uint = f.hashPair(item)

	for i := uint(0); i < f.k; i++ {
		var index uint = (a + b * i) % f.m
		if !f.data[index] {
			return false
		}
	}

	return true
}