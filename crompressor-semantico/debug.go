package main

import (
	"fmt"
	"github.com/MrJc01/crompressor-semantico/pesquisa5/pkg/lsh"
	"github.com/MrJc01/crompressor-semantico/pkg/hamming"
)

func main() {
	lsh.CarregarMatrizPCA("pesquisa6/data_engine/matriz_pca_conversacional.json")

	embBase := make([]float32, 384)
	for i := 0; i < 50; i++ { embBase[i] = 0.5 }
	
	embInput := make([]float32, 384)
	for i := 0; i < 42; i++ { embInput[i] = 0.5 }
	for i := 42; i < 80; i++ { embInput[i] = 0.3 }

	hE1, _, _ := lsh.GerarSimHashPCA_Multi(embBase)
	hE2, _, _ := lsh.GerarSimHashPCA_Multi(embInput)

	dist := hamming.Distance(hE1, hE2)
	fmt.Printf("Distancia Entidade: %d\n", dist)
}
