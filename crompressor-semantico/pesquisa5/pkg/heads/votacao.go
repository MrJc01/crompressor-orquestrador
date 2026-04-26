package heads

import (
	"github.com/MrJc01/crompressor-semantico/pkg/hamming"
)

// MultiHash representa as múltiplas assinaturas paralelas extraídas de um mesmo dado (Pesquisa 5).
type MultiHash struct {
	Entidade uint64 // Hash focado no "sujeito/objeto"
	Contexto uint64 // Hash focado em "relação/verbo"
	Visual   uint64 // Hash focado em "textura/pixel"
}

// VotacaoConsenso aplica uma regra de "Multi-Cérebros".
// Se a distância da entidade for maior que o limite, a atenção hierárquica rejeita imediatamente.
// Caso contrário, calcula-se o consenso das cabeças ativas.
func VotacaoConsenso(a, b MultiHash, limiteBitsEntidade int) (bool, float64) {
	distEntidade := hamming.Distance(a.Entidade, b.Entidade)
	
	// Atenção Hierárquica: Filtragem rápida
	if distEntidade > limiteBitsEntidade {
		return false, 0.0
	}
	
	// Se passou pelo filtro de entidade, avalia o contexto
	distContexto := hamming.Distance(a.Contexto, b.Contexto)
	
	// Consenso linear simples para o stub inicial
	media := float64(distEntidade + distContexto) / 2.0
	
	// Confiança inversamente proporcional à distância
	confianca := 100.0 - (media * 2.0)
	if confianca < 0 {
		confianca = 0
	}
	
	return media <= float64(limiteBitsEntidade), confianca
}
