package hamming

import "math/bits"

// Distance calcula a distância de Hamming entre duas assinaturas (hashes semânticos).
// Retorna o número de bits diferentes (0 significa idênticos).
// O cálculo O(1) usa a instrução de hardware POPCNT (OnesCount64).
func Distance(hashA, hashB uint64) int {
	return bits.OnesCount64(hashA ^ hashB)
}

// IsSimilar verifica se dois hashes semânticos estão dentro de um limiar de similaridade (Threshold).
// O threshold indica quantos bits divergentes são tolerados.
func IsSimilar(hashA, hashB uint64, threshold int) bool {
	return Distance(hashA, hashB) <= threshold
}
