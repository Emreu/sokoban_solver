package world

import (
	"hash/fnv"
	"hash/maphash"
	"math/rand"
	"testing"
)

const hashSize = 10000
const hashMin = 4
const hashMax = 20

var hashTestData [][]byte

func init() {
	rnd := rand.New(rand.NewSource(0))
	length := hashMax - hashMin
	for i := 0; i < hashSize; i++ {
		data := make([]byte, hashMin+rnd.Intn(length))
		for j := range data {
			data[j] = byte(rnd.Int63())
		}
		hashTestData = append(hashTestData, data)
	}
}

func BenchmarkMaphash(b *testing.B) {
	var hash maphash.Hash

	for i := 0; i < b.N; i++ {
		for _, data := range hashTestData {
			hash.Reset()
			hash.Write(data)
			_ = hash.Sum64()
		}
	}
}

func BenchmarkFNV(b *testing.B) {
	hash := fnv.New64()

	for i := 0; i < b.N; i++ {
		for _, data := range hashTestData {
			hash.Reset()
			hash.Write(data)
			_ = hash.Sum64()
		}
	}
}

func BenchmarkFNVa(b *testing.B) {
	hash := fnv.New64a()

	for i := 0; i < b.N; i++ {
		for _, data := range hashTestData {
			hash.Reset()
			hash.Write(data)
			_ = hash.Sum64()
		}
	}
}
