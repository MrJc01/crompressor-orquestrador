package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/MrJc01/crompressor-semantico/pesquisa2/motor_avancado/pkg/lsh"
	"github.com/MrJc01/crompressor-semantico/pesquisa2/motor_avancado/pkg/vision"
	lshBase "github.com/MrJc01/crompressor-semantico/pkg/lsh"
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

// combinacaoContextual une o Hash Atual com o Hash Anterior usando intercalação 32/32.
// Ele preserva a metade superior do Contexto (Hash Antigo) e a metade inferior da Intenção (Hash Atual).
func combinacaoContextual(hashAtual, hashAnterior uint64) uint64 {
	if hashAnterior == 0 {
		return hashAtual // Primeira mensagem, sem contexto.
	}
	contexto := hashAnterior & 0xFFFFFFFF00000000
	intencao := hashAtual & 0x00000000FFFFFFFF
	return contexto | intencao
}

func main() {
	fmt.Println("==========================================================")
	fmt.Println("🤖 CROM-LLM v2: Motor Industrial Multimodal (Pesquisa 2)")
	fmt.Println("==========================================================")

	arquivo, err := ioutil.ReadFile("../../pesquisa2/laboratorio_dados/dataset_maduro_p2.json")
	if err != nil {
		log.Fatalf("Erro ao ler JSON: %v", err)
	}

	var ds Dataset
	json.Unmarshal(arquivo, &ds)

	lshBase.InitHiperplanos(384)

	// ==========================================
	// 1. CHAT COM BUFFER DE CONTEXTO
	// ==========================================
	fmt.Println("\n[🗨️  MODO CHAT] Iniciando Memória de Contexto com Multi-Index Hashing...")
	
	mihChat := lsh.NovoMultiIndex()
	limiarChat := 15 // Aceita ruido gaussiano alto graças ao MIH
	var hashContexto uint64 = 0

	for _, msg := range ds.Chat {
		fmt.Printf("\n👤 Usuário: \"%s\"\n", msg.Texto)
		hashBruto := lshBase.GerarSimHash(msg.Embedding)
		
		// O Pulo do Gato: Buffer de Contexto!
		hashBusca := combinacaoContextual(hashBruto, hashContexto)

		// Busca nas 4 tabelas O(1)
		_, encontrou := mihChat.BuscaRapida(hashBusca, limiarChat)

		if encontrou {
			fmt.Println("   >> ⚡ [CONTEXTO RECUPERADO] Lógica deduzida via Multi-Index Hashing O(1).")
			fmt.Println("   >> 🤖 CROM: Reconheci a intenção baseada no histórico!")
		} else {
			fmt.Println("   >> 🧠 [NOVO CONTEXTO] Computando resposta original...")
			fmt.Println("   >> 🤖 CROM: Registrando novo raciocínio no MIH.")
			mihChat.Inserir(hashBusca)
		}
		
		// Atualiza o contexto para a PRÓXIMA mensagem ser influenciada por esta.
		hashContexto = hashBruto 
	}

	// ==========================================
	// 2. VISÃO COM PROBABILIDADE FOVEAL
	// ==========================================
	fmt.Println("\n[👁️  MODO VISÃO] Iniciando Análise de Densidade Foveal (Overlapping)...")
	
	motorVisao := vision.NovoDicionarioVisao(10)
	
	imgA := ds.Visao[0] // Maçã no centro
	fmt.Printf("[*] Ingerindo a %s no Motor...\n", imgA.Descricao)
	for _, bloco := range imgA.Blocos {
		hashPatch := lshBase.GerarSimHash(bloco.Embedding)
		motorVisao.Injetar("MACA_RECONHECIDA", hashPatch)
	}

	imgB := ds.Visao[1] // Pera no centro
	fmt.Printf("[*] Analisando %s (Teste de Fóvea)...\n", imgB.Descricao)
	
	var hashesTeste []uint64
	var pesosTeste []float32
	
	for _, bloco := range imgB.Blocos {
		hashesTeste = append(hashesTeste, lshBase.GerarSimHash(bloco.Embedding))
		pesosTeste = append(pesosTeste, bloco.Peso)
	}

	previsoes := motorVisao.AnalisarDensidadeFoveal(hashesTeste, pesosTeste)
	
	fmt.Println(">> Resultados da Invariância Foveal:")
	if len(previsoes) == 0 {
		fmt.Println("   ❌ Objeto desconhecido. (Rejeição bem sucedida!)")
	} else {
		for _, p := range previsoes {
			fmt.Printf("   🎯 %s\n", p.String())
		}
	}
}
