package main

import (
	"fmt"
	"github.com/MrJc01/crompressor-semantico/pesquisa5/pkg/lsh"
	"github.com/MrJc01/crompressor-semantico/pesquisa6/pkg/nlp"
	"github.com/MrJc01/crompressor-semantico/pkg/hamming"
)

func main() {
	lsh.CarregarMatrizPCA("pesquisa6/data_engine/matriz_pca_conversacional.json")
	nlp.InicializarCerebroLexical("pesquisa6/data_engine/vocabulario.json", "pesquisa6/data_engine/dataset_vetorizado.json")

	texto := "meu nome é jorge"
	embInput := nlp.GerarEmbeddingTFIDF(texto)
	iEnt, _, _ := lsh.GerarSimHashPCA_Multi(embInput)

	fmt.Printf("Texto: %s\nHash: %016x\n", texto, iEnt)
	for _, doc := range nlp.DatasetReal {
		dist := hamming.Distance(iEnt, doc.HashEntidade)
		fmt.Printf("Doc: %-30s Dist: %d\n", doc.Intent, dist)
	}
}
