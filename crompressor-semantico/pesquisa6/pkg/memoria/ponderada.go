package memoria

import (
	"github.com/MrJc01/crompressor-semantico/pesquisa6/pkg/heads"
	"math/rand"
)

// Ponderada introduz a funcionalidade de Anchoring (Pesquisa 6)
type Ponderada struct {
	Base       heads.MultiHash
	Ancoragem  uint64  // Máscara de bits na Entidade que representam preferências/nome do utilizador
	Relevancia float64 // 0.0 a 1.0
}

// NovaMemoria inicializa a memória com bits protegidos
func NovaMemoria(hash heads.MultiHash, ancoragem uint64, relevancia float64) *Ponderada {
	return &Ponderada{
		Base:       hash,
		Ancoragem:  ancoragem,
		Relevancia: relevancia,
	}
}

// DecaimentoSeletivo aplica erosão temporal, mas protege os bits Ancorados (Longo Prazo)
func (m *Ponderada) DecaimentoSeletivo(taxa float64) heads.MultiHash {
	// A relevância da conversa atual abranda a corrosão global da entidade
	chanceEntidade := taxa * (1.0 - m.Relevancia)
	chanceContexto := taxa * 1.5 // Diálogos perdem contexto de 'tom' muito depressa
	chanceVisual   := taxa * 2.0 // Elementos visuais não mencionados são esquecidos quase de imediato
	
	corrompido := heads.MultiHash{
		Entidade: corromperBitsAncorados(m.Base.Entidade, chanceEntidade, m.Ancoragem),
		Contexto: corromperBits(m.Base.Contexto, chanceContexto),
		Visual:   corromperBits(m.Base.Visual, chanceVisual),
	}
	return corrompido
}

// corromperBitsAncorados garante que assuntos fundamentais do utilizador nunca decaiam
func corromperBitsAncorados(hash uint64, prob float64, ancoragem uint64) uint64 {
	resultado := hash
	for i := 0; i < 64; i++ {
		bitMask := uint64(1 << i)
		if (ancoragem & bitMask) != 0 {
			// Bit ancorado está protegido (memória profunda/constante)
			continue
		}
		if rand.Float64() < prob {
			resultado ^= bitMask
		}
	}
	return resultado
}

// corromperBits simula perda sináptica genérica
func corromperBits(hash uint64, prob float64) uint64 {
	resultado := hash
	for i := 0; i < 64; i++ {
		if rand.Float64() < prob {
			resultado ^= (1 << i)
		}
	}
	return resultado
}
