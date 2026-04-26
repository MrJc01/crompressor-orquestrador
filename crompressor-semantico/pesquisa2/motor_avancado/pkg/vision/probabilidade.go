package vision

import (
	"fmt"
	"github.com/MrJc01/crompressor-semantico/pkg/hamming"
)

type PrevisaoFoveal struct {
	Classe       string  `json:"objeto"`
	Confianca    float32 `json:"confianca"`
}

type DicionarioVisao struct {
	Conhecimento map[string][]uint64
	Limiar       int
}

func NovoDicionarioVisao(limiar int) *DicionarioVisao {
	return &DicionarioVisao{
		Conhecimento: make(map[string][]uint64),
		Limiar:       limiar,
	}
}

func (dv *DicionarioVisao) Injetar(classe string, hash uint64) {
	dv.Conhecimento[classe] = append(dv.Conhecimento[classe], hash)
}

// AnalisarDensidadeFoveal implementa a Invariância baseada em Overlapping Patches.
// Recebe uma lista de pesos onde os patches centrais (fóvea) têm peso maior.
func (dv *DicionarioVisao) AnalisarDensidadeFoveal(patchesHash []uint64, pesos []float32) []PrevisaoFoveal {
	var previsoes []PrevisaoFoveal
	var somaPesosTotais float32 = 0
	
	for _, p := range pesos {
		somaPesosTotais += p
	}

	for nomeClasse, hashesConhecidos := range dv.Conhecimento {
		var scoreDensidade float32 = 0
		
		for i, ph := range patchesHash {
			match := false
			for _, hashClasse := range hashesConhecidos {
				dist := hamming.Distance(ph, hashClasse)
				if dist <= dv.Limiar {
					match = true
					break
				}
			}
			
			if match {
				scoreDensidade += pesos[i] // Se o centro bate, o score sobe violentamente.
			}
		}

		confianca := (scoreDensidade / somaPesosTotais) * 100.0
		if confianca > 0 {
			previsoes = append(previsoes, PrevisaoFoveal{Classe: nomeClasse, Confianca: confianca})
		}
	}

	return previsoes
}

func (p PrevisaoFoveal) String() string {
	return fmt.Sprintf("Objeto: %s | Confiança Ponderada: %.2f%%", p.Classe, p.Confianca)
}
