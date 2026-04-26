package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/MrJc01/crompressor-semantico/pkg/hamming"
	"github.com/MrJc01/crompressor-semantico/pkg/lsh"
)

func main() {
	fmt.Println("==================================================")
	fmt.Println("🚀 CROM-SIM: Motor de Deduplicação Semântica O(1)")
	fmt.Println("==================================================")

	rand.Seed(time.Now().UnixNano())
	dimensaoEmbedding := 512

	// 1. Inicializando Dicionário de Hiperplanos
	fmt.Printf("[*] Inicializando Dicionário LSH com %d dimensões...\n", dimensaoEmbedding)
	lsh.InitHiperplanos(dimensaoEmbedding)

	// 2. Simulação de Extração de Features (A e B são similares)
	fmt.Println("[*] Simulando a extração de 'Embeddings' de duas fotos da mesma Maçã...")
	macaFotoA := make([]float32, dimensaoEmbedding)
	macaFotoB := make([]float32, dimensaoEmbedding) // B é ligeiramente diferente de A
	bananaFotoC := make([]float32, dimensaoEmbedding) // C é completamente diferente

	for i := 0; i < dimensaoEmbedding; i++ {
		baseVal := rand.Float32()
		macaFotoA[i] = baseVal
		// Foto B tem ruído/mudança de luz (até 10% de diferença no tensor)
		ruido := (rand.Float32() - 0.5) * 0.1
		macaFotoB[i] = baseVal + ruido
		
		// Foto C é outra coisa
		bananaFotoC[i] = rand.Float32()
	}

	// 3. Gerando Hash Semântico
	hashA := lsh.GerarSimHash(macaFotoA)
	hashB := lsh.GerarSimHash(macaFotoB)
	hashC := lsh.GerarSimHash(bananaFotoC)

	fmt.Printf("\nHash Foto A (Maçã):   %064b\n", hashA)
	fmt.Printf("Hash Foto B (Maçã):   %064b\n", hashB)
	fmt.Printf("Hash Foto C (Banana): %064b\n\n", hashC)

	// 4. Teste de Deduplicação usando Hamming
	distanciaAB := hamming.Distance(hashA, hashB)
	distanciaAC := hamming.Distance(hashA, hashC)

	limiarAceitavel := 5 // Toleramos até 5 bits de diferença para considerar igual

	fmt.Printf("[>] Distância de Hamming (Maçã A -> Maçã B): %d bits diferentes\n", distanciaAB)
	if distanciaAB <= limiarAceitavel {
		fmt.Println("    ✅ DEDUPLICADO: O motor reconhece como sendo a MESMA SEMÂNTICA.")
	} else {
		fmt.Println("    ❌ NÃO DEDUPLICADO: Os dados são muito diferentes.")
	}

	fmt.Printf("[>] Distância de Hamming (Maçã A -> Banana C): %d bits diferentes\n", distanciaAC)
	if distanciaAC <= limiarAceitavel {
		fmt.Println("    ✅ DEDUPLICADO: O motor reconhece como sendo a MESMA SEMÂNTICA.")
	} else {
		fmt.Println("    ❌ NÃO DEDUPLICADO: Os dados são muito diferentes.")
	}
}
