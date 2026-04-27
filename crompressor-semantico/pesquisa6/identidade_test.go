package pesquisa6_test

import (
	"strings"
	"testing"

	"github.com/MrJc01/crompressor-semantico/pesquisa5/pkg/lsh"
	"github.com/MrJc01/crompressor-semantico/pesquisa6/pkg/nlp"
	"github.com/MrJc01/crompressor-semantico/pesquisa6/pkg/probabilidade"
	"github.com/MrJc01/crompressor-semantico/pkg/hamming"
)

func TestIdentidadeCROM(t *testing.T) {
	// Carregar Cerebro
	err := lsh.CarregarMatrizPCA("data_engine/matriz_pca_conversacional.json")
	if err != nil {
		t.Fatalf("Erro ao carregar PCA: %v", err)
	}

	err = nlp.InicializarCerebroLexical("data_engine/vocabulario.json", "data_engine/dataset_vetorizado.json")
	if err != nil {
		t.Fatalf("Erro ao carregar Lexical: %v", err)
	}

	testes := []struct {
		query                   string
		esperadoConfiancaMinima float64
		esperadoMapeadosMinimo  int
	}{
		{"oi", 80.0, 1},
		{"o que é o universo?", 80.0, 2},
	}

	for _, tt := range testes {
		embInput, unkTokens, matchCount := nlp.GerarEmbeddingTFIDF(tt.query)
		if matchCount < tt.esperadoMapeadosMinimo {
			t.Errorf("Query '%s': Esperava pelo menos %d mapeados, obteve %d. UNK: %v", tt.query, tt.esperadoMapeadosMinimo, matchCount, unkTokens)
		}

		iEnt, _, _ := lsh.GerarSimHashPCA_Multi(embInput)

		melhorDistancia := 9999
		for _, doc := range nlp.DatasetReal {
			dist := hamming.Distance(iEnt, doc.HashEntidade)
			if dist < melhorDistancia {
				melhorDistancia = dist
			}
		}

		tamanhoQuery := len(strings.Fields(tt.query))
		confianca := probabilidade.CalcularConfiancaSigmoide(melhorDistancia, tamanhoQuery, 1.0)

		if confianca < tt.esperadoConfiancaMinima {
			t.Errorf("Query '%s': Confiança muito baixa (%.2f%%). Distância foi %d bits.", tt.query, confianca, melhorDistancia)
		} else {
			t.Logf("Query '%s' -> OK: Mapeados=%d, Dist=%d, Confianca=%.2f%%", tt.query, matchCount, melhorDistancia, confianca)
		}
	}
}
