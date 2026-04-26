package memoria

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// DecaimentoContexto simula a meia-vida da memória orgânica.
// Uma taxa de 0.1 significa que 10% dos bits do hash de contexto 
// perderão sua coerência (sofrerão bit-flipping ou zeroing out aleatório),
// tornando a memória menos influente na busca quanto mais tempo passa.
func DecaimentoContexto(hashAnterior uint64, taxaEsquecimento float64) uint64 {
	if hashAnterior == 0 {
		return 0
	}

	hashDecaido := hashAnterior
	for i := 0; i < 64; i++ {
		// Se o dado bit for sorteado na taxa de esquecimento, ele perde integridade (vira 0)
		if rand.Float64() < taxaEsquecimento {
			mascara := uint64(1) << i
			hashDecaido = hashDecaido &^ mascara // Limpa o bit
		}
	}
	return hashDecaido
}
