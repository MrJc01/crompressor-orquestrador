package vision

import (
	"fmt"
	"github.com/MrJc01/crompressor-semantico/pkg/hamming"
)

type Previsao struct {
	Classe        string
	Probabilidade float32
}

// O DicionarioConceitos simula um banco de dados pré-treinado onde associamos
// Hash Semânticos a nomes de Classes ("maca_vermelha", "fundo").
// Na vida real, isso estaria no Bucketing persistido no disco.
type DicionarioConceitos struct {
	// A chave é a string da Classe, o valor são os Hashes conhecidos que a representam.
	// Isso permite que "Maçã Vermelha" tenha vários hashes diferentes de vários ângulos.
	Conhecimento map[string][]uint64
	Limiar       int // Limiar de bits para considerar um "Match"
}

func NovoDicionarioConceitos(limiar int) *DicionarioConceitos {
	return &DicionarioConceitos{
		Conhecimento: make(map[string][]uint64),
		Limiar:       limiar,
	}
}

// TreinarInjetar associa manualmente um hash a uma classe (populando a "memória").
func (dc *DicionarioConceitos) TreinarInjetar(classe string, hash uint64) {
	dc.Conhecimento[classe] = append(dc.Conhecimento[classe], hash)
}

// PreverObjeto recebe 16 hashes (os patches de uma imagem) e cruza com a memória do CROM.
func (dc *DicionarioConceitos) PreverObjeto(patchesHash []uint64) []Previsao {
	var previsoes []Previsao
	totalPatches := float32(len(patchesHash))

	// Para cada Classe que conhecemos no dicionário...
	for nomeClasse, hashesConhecidos := range dc.Conhecimento {
		acertos := 0
		
		// Verificamos quantos patches da imagem batem com a nossa classe
		for _, ph := range patchesHash {
			// Procura no conhecimento se este patch é "parecido" com algo dessa classe
			match := false
			for _, hashClasse := range hashesConhecidos {
				dist := hamming.Distance(ph, hashClasse)
				if dist <= dc.Limiar {
					match = true
					break // Encontrou pelo menos um ângulo/patch igual
				}
			}
			if match {
				acertos++
			}
		}

		prob := (float32(acertos) / totalPatches) * 100.0
		if prob > 0 {
			previsoes = append(previsoes, Previsao{Classe: nomeClasse, Probabilidade: prob})
		}
	}

	return previsoes
}

func (p Previsao) String() string {
	return fmt.Sprintf("Classe: %s | Confiança: %.2f%%", p.Classe, p.Probabilidade)
}
