package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/MrJc01/crompressor-semantico/pkg/lsh"
	"github.com/MrJc01/crompressor-semantico/pkg/vision"
	"github.com/MrJc01/crompressor-semantico/pkg/hamming"
)

type TextoData struct {
	ID        string    `json:"id"`
	Texto     string    `json:"texto"`
	Embedding []float32 `json:"embedding"`
}

type ImagemData struct {
	ID        string      `json:"id"`
	Descricao string      `json:"descricao"`
	Patches   [][]float32 `json:"patches"`
}

type Dataset struct {
	Textos  []TextoData  `json:"textos"`
	Imagens []ImagemData `json:"imagens"`
}

func main() {
	fmt.Println("==================================================")
	fmt.Println("🤖 CROM-LLM: Motor Multimodal (Visão & Chat)")
	fmt.Println("==================================================")

	// Lendo Massa de Dados
	arquivo, err := ioutil.ReadFile("../../pesquisa1/resultados/dataset_maduro.json")
	if err != nil {
		log.Fatalf("Erro ao ler JSON: %v", err)
	}

	var ds Dataset
	json.Unmarshal(arquivo, &ds)

	dimensao := len(ds.Textos[0].Embedding)
	lsh.InitHiperplanos(dimensao)

	// ==========================================
	// 1. Simulação do Chat Coerente (Cache KV O(1))
	// ==========================================
	fmt.Println("\n[🗨️  MODO CHAT] Iniciando Memória de Curto Prazo...")
	
	// Configuramos Bucketing priorizando os dois (Prefixo de 8 bits = 256 Buckets)
	cacheSemantico := lsh.NovoDicionario(8)
	limiarDeduplicacao := 25 // Ajustado para lidar com o ruído gaussiano do script Python

	// Banco simulado de "Respostas" atreladas a um Hash de pergunta
	respostasCache := make(map[uint64]string)

	for _, t := range ds.Textos {
		fmt.Printf("\nUsuário: \"%s\"\n", t.Texto)
		hashPergunta := lsh.GerarSimHash(t.Embedding)

		// Roteamento O(1): Procura no Bucket
		hashEncontrado, existeSimilar := cacheSemantico.BuscarSimilar(hashPergunta, limiarDeduplicacao)

		if existeSimilar {
			// Deduplicado! A entropia da pergunta é quase igual.
			dist := hamming.Distance(hashPergunta, hashEncontrado)
			fmt.Printf("   >> ⚡ [DEDUPLICADO] Pergunta já respondida antes. (Distância: %d bits)\n", dist)
			fmt.Printf("   >> 🤖 CROM: %s\n", respostasCache[hashEncontrado])
		} else {
			// Não encontrou no cache, precisa "Gerar" via LLM real.
			fmt.Println("   >> 🧠 [NOVO CONTEXTO] Computando resposta do zero (Simulação de inferência lenta)...")
			
			// Simulamos a resposta da rede neural...
			var resposta string
			if t.ID == "T1" {
				resposta = "A cotação do dólar hoje é de R$ 5,00."
			} else {
				resposta = "Não tenho certeza sobre isso, mas vou pesquisar."
			}

			// Salva no Bucket para as próximas interações
			cacheSemantico.Inserir(hashPergunta)
			respostasCache[hashPergunta] = resposta
			fmt.Printf("   >> 🤖 CROM: %s\n", resposta)
		}
	}

	cacheSemantico.DebugDensidade()

	// ==========================================
	// 2. Simulação de Visão (Patch-Hash)
	// ==========================================
	fmt.Println("\n[👁️  MODO VISÃO] Iniciando Reconhecimento Probabilístico...")
	limiarVisao := 10
	motorVisao := vision.NovoDicionarioConceitos(limiarVisao)

	// 2.1 Fase de "Treinamento / Ingestão"
	imgA := ds.Imagens[0]
	fmt.Printf("[*] Ingerindo a %s no Dicionário...\n", imgA.Descricao)
	for _, patch := range imgA.Patches {
		hashPatch := lsh.GerarSimHash(patch)
		// A rigor, o patch deveria ter classe separada (fundo vs maca),
		// Mas para simplificar, atrelamos todos os hashes da ImgA à classe "Maçã Vermelha".
		motorVisao.TreinarInjetar("Classe_Maca_Vermelha", hashPatch)
	}

	// 2.2 Fase de "Validação Cega"
	imgB := ds.Imagens[1]
	fmt.Printf("[*] Analisando a %s...\n", imgB.Descricao)
	
	// Geramos os 16 hashes da nova imagem
	var hashesImgB []uint64
	for _, patch := range imgB.Patches {
		hashesImgB = append(hashesImgB, lsh.GerarSimHash(patch))
	}

	// Invocamos a previsão
	previsoes := motorVisao.PreverObjeto(hashesImgB)
	
	fmt.Println(">> Resultados da Análise Patch-Hash (16 Blocos):")
	if len(previsoes) == 0 {
		fmt.Println("   ❌ Objeto desconhecido.")
	} else {
		for _, p := range previsoes {
			fmt.Printf("   🎯 %s\n", p.String())
		}
	}
}
