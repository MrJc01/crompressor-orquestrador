package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/MrJc01/crompressor-semantico/pkg/hamming"
	"github.com/MrJc01/crompressor-semantico/pkg/lsh"
)

type DadoReal struct {
	ID        int       `json:"id"`
	Texto     string    `json:"texto"`
	Embedding []float32 `json:"embedding"`
}

func main() {
	fmt.Println("=====================================================")
	fmt.Println("🧠 CROM LLM-TEXT: Inferência O(1) com Dados Reais")
	fmt.Println("=====================================================")

	// 1. Lendo os dados gerados pelo laboratório Python
	caminhoArquivo := "../resultados/dados_reais.json"
	arquivo, err := ioutil.ReadFile(caminhoArquivo)
	if err != nil {
		log.Fatalf("Erro ao ler %s: %v", caminhoArquivo, err)
	}

	var dados []DadoReal
	err = json.Unmarshal(arquivo, &dados)
	if err != nil {
		log.Fatalf("Erro no unmarshal do JSON: %v", err)
	}

	if len(dados) == 0 {
		log.Fatal("Nenhum dado encontrado no JSON.")
	}

	dimensaoEmbedding := len(dados[0].Embedding)
	fmt.Printf("[*] Foram carregadas %d frases com dimensão %d.\n", len(dados), dimensaoEmbedding)

	// 2. Inicializar o LSH com semente determinística (fixa) para que o Hash seja estável
	// (Na lib original não passamos semente pro Init, mas pra teste vamos apenas chamar Init)
	fmt.Println("[*] Gerando Hiperplanos para Hash Semântico de Texto...")
	lsh.InitHiperplanos(dimensaoEmbedding)

	// 3. Gerar Hashes
	hashes := make(map[int]uint64)
	fmt.Println("\n--- Hashes Gerados ---")
	for _, d := range dados {
		hashVal := lsh.GerarSimHash(d.Embedding)
		hashes[d.ID] = hashVal
		fmt.Printf("ID %d: %-35s -> %064b\n", d.ID, `"`+d.Texto+`"`, hashVal)
	}

	// 4. Teste de Deduplicação de Contexto (Memory-Bound LLM test)
	fmt.Println("\n--- Comparação de Similaridade (Deduplicação de Cache KV) ---")
	limiar := 10 // tolerância de bits

	// Vamos comparar "o carro veloz" (ID 1) com os outros
	idAlvo := 1
	textoAlvo := dados[0].Texto
	hashAlvo := hashes[idAlvo]

	fmt.Printf("[+] Analisando o prompt: \"%s\"\n", textoAlvo)

	for _, d := range dados {
		if d.ID == idAlvo {
			continue
		}
		dist := hamming.Distance(hashAlvo, hashes[d.ID])
		status := "❌ Novo Contexto (Não Deduplicado)"
		if dist <= limiar {
			status = "✅ DEDUPLICADO (Reusar Cache KV O(1))"
		}
		fmt.Printf(" -> vs \"%s\": Distância = %2d | Status: %s\n", d.Texto, dist, status)
	}
}
