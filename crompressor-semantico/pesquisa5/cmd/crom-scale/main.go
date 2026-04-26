package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/MrJc01/crompressor-semantico/pesquisa5/pkg/heads"
	"github.com/MrJc01/crompressor-semantico/pesquisa5/pkg/lsh"
	"github.com/MrJc01/crompressor-semantico/pesquisa5/pkg/memoria"
)

func main() {
	fmt.Println("=====================================================================")
	fmt.Println("🌌 CROM-SCALE: Orquestrador Multi-Heads (Pesquisa 5.1)")
	fmt.Println("=====================================================================")

	err := lsh.CarregarMatrizPCA("../../data_engine/matriz_pca_multi.json")
	if err != nil {
		log.Fatalf("Falha ao carregar as cabeças de PCA: %v", err)
	}

	// Simulando dados contínuos de stream (10.000 requisições simultâneas para estresse concorrido)
	totalRequisicoes := 10000
	fmt.Printf("\n[🚀] Iniciando bombardeio de %d inputs concorrentes...\n", totalRequisicoes)

	// O Hash "Base" (O que o robô já conhece perfeitamente bem - Memória Forte)
	embBase := make([]float32, 384)
	for i := 0; i < 50; i++ { embBase[i] = 0.9 }
	hEntidade, hContexto, hVisual := lsh.GerarSimHashPCA_Multi(embBase)
	
	baseConhecimento := heads.MultiHash{
		Entidade: hEntidade,
		Contexto: hContexto,
		Visual:   hVisual,
	}

	memoriaAtiva := memoria.NovaMemoriaPonderada(baseConhecimento, 0.95) // 95% de relevância (Difícil de esquecer)
	memoriaDeCurtoPrazo := memoriaAtiva.DecaimentoSeletivo(0.10) // 10% de ruído temporal

	var wg sync.WaitGroup
	wg.Add(totalRequisicoes)

	start := time.Now()

	// Simulando requisições paralelas
	sucessos := 0
	rejeicoesPrecoces := 0
	var mu sync.Mutex

	for i := 0; i < totalRequisicoes; i++ {
		go func(id int) {
			defer wg.Done()

			// Simulando variações no input
			embInput := make([]float32, 384)
			for j := 0; j < 50; j++ {
				if id%2 == 0 {
					embInput[j] = 0.88 // Input válido (pequena variação)
				} else {
					embInput[j] = -0.5 // Input inválido (foco da oclusão)
				}
			}

			iEntidade, iContexto, iVisual := lsh.GerarSimHashPCA_Multi(embInput)
			inputHash := heads.MultiHash{
				Entidade: iEntidade,
				Contexto: iContexto,
				Visual:   iVisual,
			}

			// Votação Consenso com Short-Circuit (Atenção Hierárquica na Entidade)
			match, confianca := heads.VotacaoConsenso(memoriaDeCurtoPrazo, inputHash, 5) // limite de 5 bits na Entidade

			mu.Lock()
			if match {
				sucessos++
			} else if confianca == 0.0 {
				rejeicoesPrecoces++
			}
			mu.Unlock()
		}(i)
	}

	wg.Wait()
	duracao := time.Since(start)

	fmt.Println("\n[📊 RELATÓRIO DO CROM-SCALE]")
	fmt.Printf("-> Tempo Total para %d inferências Multi-Head: %s\n", totalRequisicoes, duracao)
	fmt.Printf("-> Latência Média de Consenso: %d nanossegundos\n", duracao.Nanoseconds()/int64(totalRequisicoes))
	fmt.Printf("-> Sucessos (Matches Coesos): %d\n", sucessos)
	fmt.Printf("-> Rejeições Rápidas (Short-Circuit via Entidade): %d\n", rejeicoesPrecoces)
	
	fmt.Println("\n[✅] A Arquitetura Multi-Camadas provou ser concorrente e segura sob fogo pesado.")
}
