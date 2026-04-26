package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/MrJc01/crompressor-semantico/pesquisa2/motor_avancado/pkg/lsh"
	lshBase "github.com/MrJc01/crompressor-semantico/pkg/lsh"
	"github.com/MrJc01/crompressor-semantico/pesquisa3/motor_calibrado/pkg/probabilidade"
	"github.com/MrJc01/crompressor-semantico/pesquisa3/motor_calibrado/pkg/vision"
	"github.com/MrJc01/crompressor-semantico/pkg/hamming"
)

type BlocoImagem struct {
	Indice    int       `json:"indice"`
	Peso      float32   `json:"peso"`
	Embedding []float32 `json:"embedding"`
}

type ImagemData struct {
	ID        string        `json:"id"`
	Descricao string        `json:"descricao"`
	Blocos    []BlocoImagem `json:"blocos"`
}

type TextoData struct {
	ID        string    `json:"id"`
	Texto     string    `json:"texto"`
	Embedding []float32 `json:"embedding"`
}

type Dataset struct {
	Chat  []TextoData  `json:"chat"`
	Visao []ImagemData `json:"visao"`
}

func combinacaoContextual(hashAtual, hashAnterior uint64) uint64 {
	if hashAnterior == 0 {
		return hashAtual
	}
	return (hashAnterior & 0xFFFFFFFF00000000) | (hashAtual & 0x00000000FFFFFFFF)
}

func main() {
	fmt.Println("===================================================================")
	fmt.Println("🔥 CROM-LLM v3: Stress Test (Calibração Bayesiana e Max-Pooling)")
	fmt.Println("===================================================================")

	arquivo, err := ioutil.ReadFile("../../../laboratorio_dados/cerebro_real.json")
	if err != nil {
		log.Fatalf("Erro ao ler JSON: %v", err)
	}

	var ds Dataset
	json.Unmarshal(arquivo, &ds)
	lshBase.InitHiperplanos(384)

	// ==========================================
	// 1. STRESS TEST CHAT (RESET TRIGGER)
	// ==========================================
	fmt.Println("\n[🗨️  MODO CHAT] Iniciando Stress Test de Contexto Móvel...")
	mihChat := lsh.NovoMultiIndex()
	var hashContexto uint64 = 0

	for _, msg := range ds.Chat {
		fmt.Printf("\n👤 Usuário: \"%s\"\n", msg.Texto)
		hashBruto := lshBase.GerarSimHash(msg.Embedding)
		hashBusca := combinacaoContextual(hashBruto, hashContexto)

		var mockResposta string
		switch msg.ID {
		case "T1", "T2":
			mockResposta = "Paris fica na França, na Europa."
		case "T3":
			mockResposta = "Entendido, registrando a correção sobre a localização."
		case "T4":
			mockResposta = "A previsão para hoje é de chuva forte."
		case "T5":
			mockResposta = "Amanhã o tempo deve abrir e fazer sol."
		case "T6", "T7":
			mockResposta = "O dólar está operando em alta a R$ 5,05."
		default:
			mockResposta = "Estou gerando uma resposta complexa sobre isso."
		}

		// Busca MIH
		candidato, encontrou := mihChat.BuscaRapida(hashBusca, 20) // Margem folgada no MIH para testarmos a sigmóide depois

		if encontrou {
			dist := hamming.Distance(hashBusca, candidato)
			confianca := probabilidade.ScoreHardNegative(dist, 14) // TAREFA 4: Filtro de Rejeição (Limiar 14 bits)
			
			fmt.Printf("   >> 🔎 Hash Match Encontrado. Distância: %d bits.\n", dist)
			
			if confianca > 60.0 {
				fmt.Printf("   >> ⚡ [DEDUPLICADO] Probabilidade Logística: %.1f%%. Reutilizando resposta do Cache O(1)!\n", confianca)
				fmt.Printf("   >> 🤖 CROM: %s\n", mockResposta)
				hashContexto = hashBruto
			} else {
				// RESET TRIGGER!
				fmt.Printf("   >> ⚠️ [RESET TRIGGER ATIVADO] Distância excedeu limite de segurança (Confiança: %.1f%%).\n", confianca)
				fmt.Println("   >> 🤖 CROM: Mudança abrupta de assunto detectada. Limpando memória e gerando contexto novo...")
				fmt.Printf("   >> 🤖 CROM: %s\n", mockResposta)
				mihChat.Inserir(hashBruto) // Insere apenas a intenção pura, quebra o histórico poluído
				hashContexto = hashBruto
			}
		} else {
			fmt.Println("   >> 🧠 [NOVO CONTEXTO] Gerando resposta original do modelo...")
			fmt.Printf("   >> 🤖 CROM: %s\n", mockResposta)
			mihChat.Inserir(hashBusca)
			hashContexto = hashBruto
		}
	}

	// ==========================================
	// 2. STRESS TEST VISÃO (MAX-POOLING)
	// ==========================================
	fmt.Println("\n[👁️  MODO VISÃO] Teste de Invariância com Max-Pooling e Ruído...")
	
	imgBase := ds.Visao[0]
	imgRuido := ds.Visao[1]
	imgDeslocada := ds.Visao[2]

	var classeConhecida []uint64
	for _, b := range imgBase.Blocos {
		classeConhecida = append(classeConhecida, lshBase.GerarSimHash(b.Embedding))
	}
	
	fmt.Println("\n[*] Analisando Imagem 2 (Maçã com 20% de Ruído):")
	for _, b := range imgRuido.Blocos {
		hashTeste := lshBase.GerarSimHash(b.Embedding)
		_, conf := vision.MaxPoolingLSH(hashTeste, classeConhecida)
		if b.Indice == 5 { // Analisando apenas o centro foveal para debug
			fmt.Printf("    -> Patch Foveal Central (Idx %d): Confiança de Reconhecimento = %.2f%%\n", b.Indice, conf)
		}
	}

	fmt.Println("\n[*] Analisando Imagem 3 (Maçã Deslocada Lateralmente):")
	for _, b := range imgDeslocada.Blocos {
		hashTeste := lshBase.GerarSimHash(b.Embedding)
		dist, conf := vision.MaxPoolingLSH(hashTeste, classeConhecida)
		if b.Indice == 6 { // A maçã moveu pro patch 6 (não está mais no 5)
			fmt.Printf("    -> Patch Deslocado (Idx %d): Max-Pooling achou vizinho com dist %d. Confiança = %.2f%%\n", b.Indice, dist, conf)
		}
	}
}
