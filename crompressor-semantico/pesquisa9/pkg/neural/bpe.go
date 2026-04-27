package neural

import (
	"fmt"
	"regexp"
	"strings"
)

// BPERule representa uma regra de fusão "A B -> AB"
type BPERule struct {
	Pair   string
	Merged string
}

// BPETokenizer implementa o Byte-Pair Encoding
type BPETokenizer struct {
	Vocab      map[string]int
	IdxToToken map[int]string
	Merges     []BPERule
	MaxVocab   int
	MaskIdx    int
}

func NewBPETokenizer(maxVocab int) *BPETokenizer {
	return &BPETokenizer{
		Vocab:      make(map[string]int),
		IdxToToken: make(map[int]string),
		Merges:     []BPERule{},
		MaxVocab:   maxVocab,
	}
}

// extrairPalavras usa regex para separar palavras
func extrairPalavras(texto string) []string {
	texto = strings.ToLower(texto)
	re := regexp.MustCompile(`\b\w+\b`)
	return re.FindAllString(texto, -1)
}

// Train treina o modelo BPE no corpus fornecido
func (b *BPETokenizer) Train(corpus string) {
	fmt.Println("[*] Treinando BPE Tokenizer...")
	
	// 1. Contar frequências das palavras
	palavras := extrairPalavras(corpus)
	wordFreqs := make(map[string]int)
	for _, p := range palavras {
		wordFreqs[p]++
	}

	// 2. Inicializar vocábulos com caracteres espaçados + </w>
	vocab := make(map[string]int)
	alfabeto := make(map[string]bool)
	
	for word, freq := range wordFreqs {
		var chars []string
		for _, c := range word {
			chars = append(chars, string(c))
			alfabeto[string(c)] = true
		}
		chars = append(chars, "</w>")
		alfabeto["</w>"] = true
		vocab[strings.Join(chars, " ")] = freq
	}

	// Pré-popular dicionário de tokens
	tokenID := 0
	b.addToken("[MASK]", &tokenID)
	b.MaskIdx = b.Vocab["[MASK]"]
	
	// Adicionar alfabeto base
	for char := range alfabeto {
		b.addToken(char, &tokenID)
	}

	// 3. Iterações de fusão
	for len(b.Vocab) < b.MaxVocab {
		// Contar pares
		pairs := make(map[string]int)
		for word, freq := range vocab {
			symbols := strings.Split(word, " ")
			for i := 0; i < len(symbols)-1; i++ {
				pair := symbols[i] + " " + symbols[i+1]
				pairs[pair] += freq
			}
		}

		if len(pairs) == 0 {
			break // Nenhuma fusão possível
		}

		// Encontrar par mais frequente
		bestPair := ""
		maxFreq := -1
		for pair, freq := range pairs {
			if freq > maxFreq {
				maxFreq = freq
				bestPair = pair
			}
		}

		parts := strings.Split(bestPair, " ")
		merged := parts[0] + parts[1]

		b.Merges = append(b.Merges, BPERule{Pair: bestPair, Merged: merged})
		b.addToken(merged, &tokenID)

		// Aplicar fusão no vocab
		newVocab := make(map[string]int)
		for word, freq := range vocab {
			// Replace isolado (não usar strings.Replace direto para não cruzar fronteiras)
			// Uma forma robusta:
			symbols := strings.Split(word, " ")
			var newSymbols []string
			i := 0
			for i < len(symbols) {
				if i < len(symbols)-1 && symbols[i] == parts[0] && symbols[i+1] == parts[1] {
					newSymbols = append(newSymbols, merged)
					i += 2
				} else {
					newSymbols = append(newSymbols, symbols[i])
					i++
				}
			}
			newVocab[strings.Join(newSymbols, " ")] = freq
		}
		vocab = newVocab
	}
	
	fmt.Printf("[+] BPE treinado. Vocabulário: %d tokens.\n", len(b.Vocab))
}

func (b *BPETokenizer) addToken(token string, idCounter *int) {
	if _, exists := b.Vocab[token]; !exists {
		b.Vocab[token] = *idCounter
		b.IdxToToken[*idCounter] = token
		*idCounter++
	}
}

// Encode converte uma string num array de IDs de tokens
func (b *BPETokenizer) Encode(texto string) []int {
	palavras := extrairPalavras(texto)
	var ids []int

	for _, p := range palavras {
		var chars []string
		for _, c := range p {
			chars = append(chars, string(c))
		}
		chars = append(chars, "</w>")
		
		wordSpace := strings.Join(chars, " ")
		
		// Aplicar fusões por ordem de criação
		for _, rule := range b.Merges {
			parts := strings.Split(rule.Pair, " ")
			merged := rule.Merged
			
			symbols := strings.Split(wordSpace, " ")
			var newSymbols []string
			i := 0
			for i < len(symbols) {
				if i < len(symbols)-1 && symbols[i] == parts[0] && symbols[i+1] == parts[1] {
					newSymbols = append(newSymbols, merged)
					i += 2
				} else {
					newSymbols = append(newSymbols, symbols[i])
					i++
				}
			}
			wordSpace = strings.Join(newSymbols, " ")
		}

		// Obter IDs
		for _, sym := range strings.Split(wordSpace, " ") {
			if id, ok := b.Vocab[sym]; ok {
				ids = append(ids, id)
			} else {
				// Fallback raro: devia estar no alfabeto base, se não, ignora ou [MASK]
				ids = append(ids, b.MaskIdx)
			}
		}
	}
	return ids
}

// Decode reconstrói a string a partir dos IDs
func (b *BPETokenizer) Decode(ids []int) string {
	var sb strings.Builder
	for _, id := range ids {
		token, ok := b.IdxToToken[id]
		if ok {
			if token == "[MASK]" {
				sb.WriteString("[MASK] ")
			} else if strings.HasSuffix(token, "</w>") {
				sb.WriteString(strings.TrimSuffix(token, "</w>") + " ")
			} else {
				sb.WriteString(token)
			}
		}
	}
	return strings.TrimSpace(sb.String())
}

// DecodeTokens retorna a representação das strings isoladas para debugging
func (b *BPETokenizer) DecodeTokens(ids []int) []string {
	var res []string
	for _, id := range ids {
		token, ok := b.IdxToToken[id]
		if ok {
			res = append(res, token)
		} else {
			res = append(res, "?")
		}
	}
	return res
}
