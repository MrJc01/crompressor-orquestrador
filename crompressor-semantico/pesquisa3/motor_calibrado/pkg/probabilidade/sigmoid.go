package probabilidade

import (
	"math"
)

// CalcularConfiancaSigmoide converte a Distância de Hamming seca em uma Probabilidade Calibrada.
// Utiliza uma função logística invertida onde:
// Distância 0 -> Confiança próxima de 100%
// Distância igual ao K (ponto médio) -> Confiança de 50%
// Distâncias altas -> Confiança cai drasticamente para 0%
func CalcularConfiancaSigmoide(distancia int, pontoMedio float64, inclinacao float64) float64 {
	x := float64(distancia)
	
	// A fórmula da sigmoide padrão é: 1 / (1 + e^(-k(x - x0)))
	// Como queremos que a confiança DIMINUA conforme a distância AUMENTA, usamos +k
	expoente := inclinacao * (x - pontoMedio)
	confianca := 1.0 / (1.0 + math.Exp(expoente))
	
	return confianca * 100.0
}

// ScoreHardNegative aplica um peso brutal caso a distância ultrapasse um limiar de "rejeição".
func ScoreHardNegative(distancia int, limiarRejeicao int) float64 {
	if distancia > limiarRejeicao {
		return 0.0 // Rejeição absoluta. Impede deduplicação de palavras opostas.
	}
	// Usamos uma calibração dura: O ponto médio de aceitação será 4 bits.
	// Uma inclinação de 1.5 faz com que a curva caia rapidamente (Softmax-like)
	return CalcularConfiancaSigmoide(distancia, 4.0, 1.5)
}
