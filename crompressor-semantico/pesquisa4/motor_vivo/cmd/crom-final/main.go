package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/MrJc01/crompressor-semantico/pesquisa2/motor_avancado/pkg/lsh"
	lshBase "github.com/MrJc01/crompressor-semantico/pkg/lsh"
	
	"github.com/MrJc01/crompressor-semantico/pesquisa3/motor_calibrado/pkg/vision"
	"github.com/MrJc01/crompressor-semantico/pkg/hamming"

	"github.com/MrJc01/crompressor-semantico/pesquisa4/motor_vivo/pkg/feedback"
	"github.com/MrJc01/crompressor-semantico/pesquisa4/motor_vivo/pkg/fusao"
	"github.com/MrJc01/crompressor-semantico/pesquisa4/motor_vivo/pkg/memoria"
	"github.com/MrJc01/crompressor-semantico/pesquisa4/motor_vivo/pkg/probabilidade"
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

func main() {
	fmt.Println("=====================================================================")
	fmt.Println("🧬 CROM-LLM v4: O Organismo Artificial (Memória, Multimodal, Brain)")
	fmt.Println("=====================================================================")

	arquivo, err := ioutil.ReadFile("../../../data_engine/cerebro_producao.json")
	if err != nil {
		log.Fatalf("Erro ao ler JSON: %v", err)
	}

	var ds Dataset
	json.Unmarshal(arquivo, &ds)
	lshBase.InitHiperplanos(384)

	motorRecompensa := feedback.NovoMotorRecompensa("brain_state.json")
	mihGeral := lsh.NovoMultiIndex()
	
	// Pre-computa hashes de imagens e armazena na "visão" do robô
	var hashesImagensConhecidas []uint64
	for _, b := range ds.Visao[0].Blocos { // Maçã base
		hashesImagensConhecidas = append(hashesImagensConhecidas, lshBase.GerarSimHash(b.Embedding))
	}

	fmt.Println("\n[🧠  BOOTING] Carregando Estado do Cérebro...")
	fmt.Printf("   -> Ponto Médio de Tolerância Base: %.2f bits\n", motorRecompensa.ObterPontoMedio("default"))

	// ==========================================
	// 1. STRESS TEST MULTIMODAL E APRENDIZADO
	// ==========================================
	
	// O utilizador "mostra" uma imagem e faz uma pergunta.
	fmt.Println("\n[👁️ + 🗨️] Iniciando Fusão Multimodal...")
	
	imgTeste := ds.Visao[1] // Maçã Deslocada
	textoCor := ds.Chat[0]  // "Qual a cor disso?"
	textoSabor := ds.Chat[1] // "Qual o sabor disso?"

	// 1.1 Extrai a percepção da imagem via Max-Pooling (Invariância Visual)
	hashVisualTeste := lshBase.GerarSimHash(imgTeste.Blocos[6].Embedding) // Patch onde a maçã foi parar
	_, confiancaVisao := vision.MaxPoolingLSH(hashVisualTeste, hashesImagensConhecidas)
	fmt.Printf("[*] Visão Ativada: Objeto detectado com %.1f%% de certeza.\n", confiancaVisao)

	// 1.2 O motor "aprende" sobre a Cor e o Sabor em contextos diferentes
	hashCor := fusao.GerarHashHibrido(hashVisualTeste, lshBase.GerarSimHash(textoCor.Embedding))
	hashSabor := fusao.GerarHashHibrido(hashVisualTeste, lshBase.GerarSimHash(textoSabor.Embedding))
	
	mihGeral.Inserir(hashCor)
	mihGeral.Inserir(hashSabor)
	fmt.Println("[*] Registrando nós multimodais (MACA+COR e MACA+SABOR) no MIH...")

	// 1.3 Simulando o Aprendizado Ativo (Punição do Usuário)
	fmt.Println("\n[⚠️  FEEDBACK LOOP] Testando Sistema de Recompensa...")
	
	hashBuscaSabor := fusao.GerarHashHibrido(hashVisualTeste, lshBase.GerarSimHash(textoSabor.Embedding))
	candidato, encontrou := mihGeral.BuscaRapida(hashBuscaSabor, 20)
	
	if encontrou {
		dist := hamming.Distance(hashBuscaSabor, candidato)
		confianca := probabilidade.ScoreDinamico(dist, "default", motorRecompensa)
		fmt.Printf("   >> 🔎 Match Encontrado (Sabor). Distância: %d bits. Confiança: %.1f%%\n", dist, confianca)
		
		fmt.Println("👤 Usuário: /errado (A resposta não foi satisfatória, muito generalista)")
		motorRecompensa.PunirBucket("default")
		fmt.Printf("   >> 📉 [APRENDIZADO ATIVO] O Ponto Médio da Sigmóide foi reduzido para: %.2f bits.\n", motorRecompensa.ObterPontoMedio("default"))
		
		// Recalculando confiança após punição
		novaConfianca := probabilidade.ScoreDinamico(dist, "default", motorRecompensa)
		fmt.Printf("   >> 🔄 Recalibrando... Nova confiança para o mesmo vetor: %.1f%%\n", novaConfianca)
	}

	// ==========================================
	// 2. STRESS TEST MEMÓRIA EVOLUTIVA (Decaimento)
	// ==========================================
	fmt.Println("\n[⏳  DECAIMENTO TEMPORAL] Testando Esquecimento de Memória...")
	var hashContexto uint64 = 0

	for i := 2; i < 6; i++ { // Chat de T3 a T6 (Dólar -> Dólar -> Neutro -> Dólar)
		msg := ds.Chat[i]
		fmt.Printf("\n👤 Usuário: \"%s\"\n", msg.Texto)
		
		hashBruto := lshBase.GerarSimHash(msg.Embedding)
		
		// Fusão com o contexto da conversa, MAS o contexto decaiu (perdeu bits)
		hashContextoDecaido := memoria.DecaimentoContexto(hashContexto, 0.15) // 15% de chance de perder bits de atenção
		hashBusca := fusao.GerarHashHibrido(hashContextoDecaido, hashBruto)
		
		candidato, encontrou := mihGeral.BuscaRapida(hashBusca, 20)
		
		if encontrou {
			dist := hamming.Distance(hashBusca, candidato)
			confianca := probabilidade.ScoreDinamico(dist, "dolar_bucket", motorRecompensa) // Simula um bucket separado
			fmt.Printf("   >> ⚡ [CONTEXTO MANTIDO] O motor ainda lembra do assunto original (Distância: %d, Confiança: %.1f%%)\n", dist, confianca)
		} else {
			fmt.Println("   >> 🧠 [NOVO CONTEXTO / ESQUECIMENTO] A erosão de bits separou o hash. Iniciando nova thread neural...")
			mihGeral.Inserir(hashBusca)
		}

		hashContexto = hashBruto // Novo contexto forte para a próxima rodada
	}
	
	fmt.Println("\n=====================================================================")
	fmt.Println("🚀 O CROM-LLM atingiu maturidade absoluta. Brain State gravado!")
}
