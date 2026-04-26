package vision

import (
	"github.com/MrJc01/crompressor-semantico/pkg/hamming"
	"github.com/MrJc01/crompressor-semantico/pesquisa3/motor_calibrado/pkg/probabilidade"
)

// MaxPoolingLSH imita a técnica convolucional de MaxPooling.
// Se um objeto se move da posição 5 para a 6, o Hash do bloco 5 muda.
// O Max-Pooling avalia os vizinhos e pega a MENOR Distância de Hamming encontrada,
// garantindo que deslocamentos espaciais não afundem o score Foveal.
func MaxPoolingLSH(patchTeste uint64, classeConhecida []uint64) (melhorDistancia int, confianca float64) {
	menorDistancia := 64 // Pior caso

	for _, hashConhecido := range classeConhecida {
		dist := hamming.Distance(patchTeste, hashConhecido)
		if dist < menorDistancia {
			menorDistancia = dist
		}
	}

	// Calcula a confiança bayesiana do "vencedor" do pooling.
	// O ponto médio é brando (ex: 6 bits) para permitir invariância visual.
	conf := probabilidade.CalcularConfiancaSigmoide(menorDistancia, 6.0, 1.2)
	
	return menorDistancia, conf
}
