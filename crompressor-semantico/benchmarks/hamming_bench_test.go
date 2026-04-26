package benchmarks

import (
	"math/rand"
	"testing"
	"github.com/MrJc01/crompressor-semantico/pkg/hamming"
)

func BenchmarkHammingDistance(b *testing.B) {
	// Preparar massa de dados pseudo-aleatória
	hashA := rand.Uint64()
	hashB := rand.Uint64()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// A meta é provar que a operação leva < 1 ns
		hamming.Distance(hashA, hashB)
	}
}

func BenchmarkIsSimilar(b *testing.B) {
	hashA := rand.Uint64()
	hashB := rand.Uint64()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hamming.IsSimilar(hashA, hashB, 10)
	}
}
