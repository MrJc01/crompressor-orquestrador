package memoria

import (
	"github.com/MrJc01/crompressor-semantico/pesquisa5/pkg/heads"
	"math/rand"
)

// Ponderada gerencia a retenção baseada em relevância estatística
type Ponderada struct {
	EntidadeBase uint64
	ContextoBase uint64
	VisualBase   uint64
	Relevancia   float64 // 0.0 a 1.0 (1.0 = inesquecível)
}

// NovaMemoriaPonderada inicializa o bloco orgânico
func NovaMemoriaPonderada(hash heads.MultiHash, relevancia float64) *Ponderada {
	return &Ponderada{
		EntidadeBase: hash.Entidade,
		ContextoBase: hash.Contexto,
		VisualBase:   hash.Visual,
		Relevancia:   relevancia,
	}
}

// DecaimentoSeletivo corrói apenas as cabeças com baixa relevância
func (m *Ponderada) DecaimentoSeletivo(taxaCorrosao float64) heads.MultiHash {
	// A relevância protege a Entidade de esquecimento rápido
	chanceCorrosaoEntidade := taxaCorrosao * (1.0 - m.Relevancia)
	chanceCorrosaoContexto := taxaCorrosao // Contexto é volátil por natureza
	chanceCorrosaoVisual := taxaCorrosao * 1.5 // Visual decai muito rápido se não houver reforço

	return heads.MultiHash{
		Entidade: corromperBits(m.EntidadeBase, chanceCorrosaoEntidade),
		Contexto: corromperBits(m.ContextoBase, chanceCorrosaoContexto),
		Visual:   corromperBits(m.VisualBase, chanceCorrosaoVisual),
	}
}

// corromperBits simula a perda de sinapses (inversão de bits)
func corromperBits(hash uint64, probabilidade float64) uint64 {
	var resultado uint64 = hash
	for i := 0; i < 64; i++ {
		if rand.Float64() < probabilidade {
			resultado ^= (1 << i) // Inverte o bit
		}
	}
	return resultado
}
