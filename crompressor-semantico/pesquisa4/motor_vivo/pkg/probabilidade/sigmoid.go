package probabilidade

import (
	"math"
	"github.com/MrJc01/crompressor-semantico/pesquisa4/motor_vivo/pkg/feedback"
)

// CalcularConfiancaSigmoide converte a Distância de Hamming seca em uma Probabilidade Calibrada.
func CalcularConfiancaSigmoide(distancia int, pontoMedio float64, inclinacao float64) float64 {
	x := float64(distancia)
	expoente := inclinacao * (x - pontoMedio)
	confianca := 1.0 / (1.0 + math.Exp(expoente))
	return confianca * 100.0
}

// ScoreDinamico consulta o Brain State para definir a rigidez da calibração.
func ScoreDinamico(distancia int, bucketID string, motor *feedback.MotorRecompensa) float64 {
	pontoMedio := motor.ObterPontoMedio(bucketID)
	// Limiar de rejeição absoluto dinâmico (ex: 3x o ponto médio)
	limiarRejeicao := int(math.Ceil(pontoMedio * 3.0))
	
	if distancia > limiarRejeicao {
		return 0.0 // Rejeição absoluta
	}
	return CalcularConfiancaSigmoide(distancia, pontoMedio, 1.5)
}
