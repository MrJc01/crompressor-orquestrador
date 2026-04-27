package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/MrJc01/crompressor-semantico/pesquisa5/pkg/lsh"
	"github.com/MrJc01/crompressor-semantico/pesquisa6/pkg/heads"
	"github.com/MrJc01/crompressor-semantico/pesquisa6/pkg/memoria"
	"github.com/MrJc01/crompressor-semantico/pesquisa6/pkg/nlp"
	"github.com/MrJc01/crompressor-semantico/pesquisa6/pkg/probabilidade"
	"github.com/MrJc01/crompressor-semantico/pkg/hamming"
)

const (
	ColorReset  = "\033[0m"
	ColorUser   = "\033[36m"
	ColorBot    = "\033[32m"
	ColorSystem = "\033[33m"
	ColorStats  = "\033[90m"
)

// KnowledgeBase armazena os centróides e meta-informação
type KnowledgeBase struct {
	CentroideSaudacao uint64 `json:"centroide_saudacao"`
	RespostaSaudacao  string `json:"resposta_saudacao"`
	Versao            string `json:"versao"`
	TotalItens        int    `json:"total_itens"`
	TotalVocab        int    `json:"total_vocab"`
}

func carregarKnowledgeBase(path string) (KnowledgeBase, error) {
	var kb KnowledgeBase
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return kb, err
	}
	err = json.Unmarshal(bytes, &kb)
	return kb, err
}

func main() {
	fmt.Println(ColorSystem + "======================================================")
	fmt.Println("🚀 CROM-Chat v3: Inteligência Vetorial Real (K-NN)")
	fmt.Println("   Pesquisa 6 - Fusão Semântica + Sigmóide Adaptativa")
	fmt.Println("======================================================" + ColorReset)

	fmt.Print(ColorSystem + "[*] Carregando Cérebro PCA... " + ColorReset)
	err := lsh.CarregarMatrizPCA("../../data_engine/matriz_pca_conversacional.json")
	if err != nil {
		fmt.Printf("\nErro ao carregar PCA: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(ColorSystem + "[*] Carregando Vocabulário e Dataset Real... " + ColorReset)
	err = nlp.InicializarCerebroLexical("../../data_engine/vocabulario.json", "../../data_engine/dataset_vetorizado.json")
	if err != nil {
		fmt.Printf("\nErro ao inicializar NLP: %v\n", err)
		os.Exit(1)
	}

	// Carregar Knowledge Base (Centroid Anchoring)
	kb, errKB := carregarKnowledgeBase("../../data_engine/knowledge_base.json")
	if errKB != nil {
		fmt.Printf(ColorStats+"[!] Knowledge Base não encontrada: %v. Centroid Anchoring desativado.\n"+ColorReset, errKB)
	} else {
		fmt.Printf(ColorSystem+"[+] Knowledge Base v%s carregada (%d itens, %d vocab).\n"+ColorReset, kb.Versao, kb.TotalItens, kb.TotalVocab)
	}

	// Ancoragem base (Memória)
	var ancoragem uint64 = (1 << 1) | (1 << 2)
	baseConhecimentoZero := heads.MultiHash{Entidade: 0, Contexto: 0, Visual: 0}
	memoriaLongoPrazo := memoria.NovaMemoria(baseConhecimentoZero, ancoragem, 0.9)
	memoriaAtual := memoriaLongoPrazo.DecaimentoSeletivo(0.0)

	fmt.Println(ColorSystem + "\nPronto!\nDigite /help para ver os comandos ou apenas comece a falar." + ColorReset)

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\n" + ColorUser + "👤 Utilizador: " + ColorReset)
		if !scanner.Scan() {
			break
		}
		texto := scanner.Text()
		if strings.TrimSpace(texto) == "" {
			continue
		}

		if strings.HasPrefix(texto, "/") {
			if texto == "/exit" || texto == "/quit" {
				break
			}
			continue
		}

		start := time.Now()
		tamanhoQuery := len(strings.Fields(texto))
		
		// 1. Extração REAL do Embedding via NLP TF-IDF (com tracking OOV)
		embInput, unkTokens, matchCount := nlp.GerarEmbeddingTFIDF(texto)

		if len(unkTokens) > 0 {
			fmt.Printf(ColorStats+"[!] Tokens desconhecidos: %v\n"+ColorReset, unkTokens)
		}
		if matchCount == 0 {
			fmt.Println(ColorBot + "🤖 CROM: Conceito completamente desconhecido. Nenhum token mapeado." + ColorReset)
			continue
		}

		// 2. Projeção nas 3 Cabeças (PCA-LSH)
		iEnt, iCtx, iVis := lsh.GerarSimHashPCA_Multi(embInput)
		inputHash := heads.MultiHash{Entidade: iEnt, Contexto: iCtx, Visual: iVis}

		// 2.5. Centroid Anchoring: queries curtas (≤2 palavras) testam primeiro o centróide de saudações
		if errKB == nil && tamanhoQuery <= 2 {
			distCentroide := hamming.Distance(iEnt, kb.CentroideSaudacao)
			if distCentroide <= 8 {
				duracao := time.Since(start)
				fmt.Printf(ColorStats+"[⚙️  Pensamento: %dns | Mapeados: %d | Distância: %d bits (centróide) | Confiança: 99.0%% | ANCHORED]\n"+ColorReset,
					duracao.Nanoseconds(), matchCount, distCentroide)
				fmt.Println(ColorBot + "🤖 CROM: " + kb.RespostaSaudacao + ColorReset)
				continue
			}
		}

		// 3. Busca Vetorial K-NN (O(N) contra o Dataset Real)
		melhorDistancia := 9999
		melhorResposta := "Desculpe, não encontrei correspondência semântica."
		
		for _, doc := range nlp.DatasetReal {
			dist := hamming.Distance(inputHash.Entidade, doc.HashEntidade)
			
			if dist < melhorDistancia {
				melhorDistancia = dist
				melhorResposta = doc.Answer
			}
		}

		// 4. Guilhotina na Sigmóide Adaptativa
		confianca := probabilidade.CalcularConfiancaSigmoide(melhorDistancia, tamanhoQuery, 1.0)
		match := confianca >= 50.0

		// 5. Consenso Adaptativo (para interface/memória)
		_, _, pesos := heads.VotacaoConsensoAdaptativa(memoriaAtual, inputHash, 30, tamanhoQuery)
		memoriaAtual = memoriaLongoPrazo.DecaimentoSeletivo(0.15)
		
		duracao := time.Since(start)

		fmt.Printf(ColorStats+"[⚙️  Pensamento: %dns | Mapeados: %d | Distância: %d bits | Confiança: %.1f%% | E:%.1f C:%.1f V:%.1f]\n"+ColorReset, 
			duracao.Nanoseconds(), matchCount, melhorDistancia, confianca, pesos["entidade"], pesos["contexto"], pesos["visual"])

		fmt.Print(ColorBot + "🤖 CROM: " + ColorReset)
		if match {
			fmt.Println(melhorResposta)
		} else {
			fmt.Println("Conceito desconhecido. Novo contexto iniciado.")
		}
	}
}
