package probabilidade

import "math"

// CalcularConfiancaSigmoide converte a Distância de Hamming numa probabilidade não-linear rígida
// Agora utiliza a heurística logarítmica baseada no tamanho da query para ajustar o pMid.
func CalcularConfiancaSigmoide(distancia int, tamanhoQuery int, k float64) float64 {
	dist := float64(distancia)
	
	// Ajuste do pMid com decaimento logarítmico sugerido: pMid = 12 + (6 / log2(tamanhoQuery + 1))
	// Tratamos tamanhoQuery = 0 como 1 para não explodir
	tq := float64(tamanhoQuery)
	if tq < 1.0 {
		tq = 1.0
	}
	
	pMid := 12.0 + (6.0 / math.Log2(tq+1.0))
	
	// Função Sigmóide Invertida: f(x) = 1 / (1 + e^(k * (x - pMid)))
	prob := 1.0 / (1.0 + math.Exp(k*(dist-pMid)))
	
	return prob * 100.0
}
