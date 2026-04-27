package heads

import (
	"github.com/MrJc01/crompressor-semantico/pkg/hamming"
	"github.com/MrJc01/crompressor-semantico/pesquisa6/pkg/probabilidade"
)

type MultiHash struct {
	Entidade uint64
	Contexto uint64
	Visual   uint64
}

// VotacaoConsensoAdaptativa aplica a "Cognição Adaptativa" da Pesquisa 6
// O motor calcula a entropia da resposta e ajusta o peso vetorial em tempo real.
func VotacaoConsensoAdaptativa(a, b MultiHash, limiteBitsEntidade int, tamanhoQuery int) (bool, float64, map[string]float64) {
	distEntidade := hamming.Distance(a.Entidade, b.Entidade)
	
	// Short-circuit (Atenção Hierárquica preservada da Pesquisa 5)
	if distEntidade > limiteBitsEntidade {
		// Retornamos os pesos básicos para a UI não mostrar zeros
		return false, 0.0, map[string]float64{"entidade": 1.0, "contexto": 0.0, "visual": 0.0}
	}
	
	distContexto := hamming.Distance(a.Contexto, b.Contexto)
	distVisual := hamming.Distance(a.Visual, b.Visual)
	
	// Pesos Base (Adaptive Head Weights)
	pesoEntidade := 1.0
	pesoContexto := 1.0
	pesoVisual := 0.5 // Menor relevância em chat puramente textual
	
	// Ajuste Adaptativo
	// Ex: Se o sujeito é idêntico (distEntidade == 0) mas o contexto mudou muito, 
	// o usuário pode estar apenas trocando de 'foco/tom' no mesmo assunto.
	if distEntidade <= 2 && distContexto > 10 {
		pesoContexto = 0.4 // Diminui o peso do contexto para não causar um falso negativo sobre a mesma entidade
	}
	
	// Se for uma pergunta fortemente visual, aumentaria o pesoVisual (simulado)
	if distVisual < 5 {
		pesoVisual = 1.2
	}

	pesoTotal := pesoEntidade + pesoContexto + pesoVisual
	distPonderada := (float64(distEntidade)*pesoEntidade + float64(distContexto)*pesoContexto + float64(distVisual)*pesoVisual) / pesoTotal
	
	confianca := probabilidade.CalcularConfiancaSigmoide(int(distPonderada), tamanhoQuery, 1.0)

	
	pesosFinais := map[string]float64{
		"entidade": pesoEntidade,
		"contexto": pesoContexto,
		"visual": pesoVisual,
	}

	return distPonderada <= float64(limiteBitsEntidade), confianca, pesosFinais
}
