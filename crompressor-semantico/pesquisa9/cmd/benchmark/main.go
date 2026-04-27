package main

import (
	"fmt"
	"time"

	"github.com/MrJc01/crompressor-semantico/pesquisa9/pkg/neural"
)

func main() {
	fmt.Println("======================================================")
	fmt.Println("🏆 CROM-LLM V9 — Gauntlet de Benchmarks Automatizado")
	fmt.Println("======================================================")

	modelo := neural.NewModel(neural.DefaultConfig())
	err := modelo.Load("/home/j/Documentos/GitHub/crom/crompressor-semantico/pesquisa9/cmd/crom-chat/crom_brain.weights")
	if err != nil {
		fmt.Printf("[ERRO] Falha ao carregar modelo: %v\n", err)
		return
	}

	fmt.Printf("[OK] Modelo carregado. Vocab: %d, Parâmetros: ~374K\n\n", modelo.GetVocabSize())

	type Categoria struct {
		Nome    string
		Prompts []string
	}

	categorias := []Categoria{
		{"Factuais (Precisão)", []string{"a gravidade é", "o universo tem", "o sol converte"}},
		{"Identidade (Consistência)", []string{"eu sou o crom", "o meu nome", "fui criado como"}},
		{"Tecnologia (Conhecimento)", []string{"python é uma", "linux é um", "a internet foi"}},
		{"OOD / Inéditos (Generalização)", []string{"como fazer bolo", "receita de sopa", "quantos anos tens"}},
	}

	for _, cat := range categorias {
		fmt.Printf(">> Categoria: %s\n", cat.Nome)
		for _, prompt := range cat.Prompts {
			start := time.Now()
			result := modelo.Generate(prompt, 10)
			lat := time.Since(start)
			fmt.Printf("   [%6.2fms] %-20s → %s\n", float64(lat.Nanoseconds())/1e6, prompt, result)
		}
		fmt.Println()
	}
}
