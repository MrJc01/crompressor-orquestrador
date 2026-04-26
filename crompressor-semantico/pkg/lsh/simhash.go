package lsh

import "math/rand"

// DicionarioHiperplanos armazena 64 hiperplanos aleatórios utilizados para a projeção de LSH.
var DicionarioHiperplanos [][]float32

// InitHiperplanos inicializa a matriz de projeção LSH. Em produção, esses pesos
// viriam do Cérebro CROM pré-treinado. Para o PoC, usamos vetores aleatórios.
func InitHiperplanos(dimensao int) {
	DicionarioHiperplanos = make([][]float32, 64)
	for i := 0; i < 64; i++ {
		plano := make([]float32, dimensao)
		for j := 0; j < dimensao; j++ {
			// Valores aleatórios entre -1.0 e 1.0
			plano[j] = rand.Float32()*2 - 1.0
		}
		DicionarioHiperplanos[i] = plano
	}
}

// dotProduct faz o produto escalar entre dois vetores.
func dotProduct(a, b []float32) float32 {
	var sum float32 = 0.0
	for i := 0; i < len(a) && i < len(b); i++ {
		sum += a[i] * b[i]
	}
	return sum
}

// GerarSimHash projeta um vetor denso em uma assinatura de 64 bits (Hash Semântico).
// A genialidade do SimHash é que vetores de entrada similares no espaço flutuante
// caem do mesmo lado dos hiperplanos a maior parte das vezes, gerando bits iguais.
func GerarSimHash(embedding []float32) uint64 {
	var hash uint64 = 0
	
	// Para cada um dos 64 bits possíveis...
	for i := 0; i < 64; i++ {
		if len(DicionarioHiperplanos) > i {
			// Se o vetor cruzar o hiperplano positivamente, ligamos o bit (1)
			if dotProduct(embedding, DicionarioHiperplanos[i]) > 0 {
				hash |= (1 << i)
			}
		}
	}
	return hash
}
