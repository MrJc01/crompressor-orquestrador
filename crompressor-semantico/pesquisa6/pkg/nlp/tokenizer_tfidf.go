package nlp

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math"
	"regexp"
	"strings"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"golang.org/x/text/runes"
	"unicode"
)

// DimensaoEmbedding adaptável com base no treinamento Python
var DimensaoEmbedding int = 384

// TokenInfo armazena o índice vetorial e o peso estatístico aprendido no Python
type TokenInfo struct {
	Indice int     `json:"indice"`
	IDF    float64 `json:"idf"`
}

var Vocabulario map[string]TokenInfo

// Índice de tokens por comprimento para busca OOV rápida
var vocabPorComprimento map[int][]string

// DocumentoVetorizado representa uma conversa original exportada pelo Python
type DocumentoVetorizado struct {
	Intent       string `json:"intent"`
	Answer       string `json:"answer"`
	HashEntidade uint64 `json:"hash_entidade"`
	HashContexto uint64 `json:"hash_contexto"`
	HashVisual   uint64 `json:"hash_visual"`
}

var DatasetReal []DocumentoVetorizado

// InicializarCerebroLexical carrega o vocabulário TF-IDF e o Dataset Hashed
func InicializarCerebroLexical(pathVocab string, pathDataset string) error {
	// 1. Carregar Vocabulário
	bytesVocab, err := ioutil.ReadFile(pathVocab)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytesVocab, &Vocabulario)
	if err != nil {
		return err
	}
	DimensaoEmbedding = len(Vocabulario)
	log.Printf("[+] Vocabulário Real: %d tokens carregados.\n", len(Vocabulario))

	// Construir índice por comprimento (apenas unigramas, para fallback OOV)
	vocabPorComprimento = make(map[int][]string)
	for k := range Vocabulario {
		if !strings.Contains(k, "_") { // Só unigramas
			vocabPorComprimento[len(k)] = append(vocabPorComprimento[len(k)], k)
		}
	}

	// 2. Carregar Dataset Vetorizado
	bytesDataset, err := ioutil.ReadFile(pathDataset)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytesDataset, &DatasetReal)
	if err != nil {
		return err
	}
	log.Printf("[+] Dataset Real: %d interações carregadas com hashes pré-computados.\n", len(DatasetReal))

	return nil
}

// removerAcentos limpa os caracteres Unicode
func removerAcentos(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(t, s)
	return result
}

var tokenRegex = regexp.MustCompile(`[a-z0-9]+`)

// ExtrairTokens extrai palavras puras e gera n-grams (bigramas e trigramas)
func ExtrairTokens(texto string) []string {
	textoLimpo := strings.ToLower(removerAcentos(texto))
	palavras := tokenRegex.FindAllString(textoLimpo, -1)
	
	estimativa := len(palavras)
	if len(palavras) > 1 {
		estimativa += len(palavras) - 1
	}
	if len(palavras) > 2 {
		estimativa += len(palavras) - 2
	}
	tokens := make([]string, 0, estimativa)
	
	for _, p := range palavras {
		tokens = append(tokens, p)
	}
	
	var builder strings.Builder
	// Bigramas
	for i := 0; i < len(palavras)-1; i++ {
		builder.Reset()
		builder.WriteString(palavras[i])
		builder.WriteString("_")
		builder.WriteString(palavras[i+1])
		tokens = append(tokens, builder.String())
	}
	// Trigramas
	for i := 0; i < len(palavras)-2; i++ {
		builder.Reset()
		builder.WriteString(palavras[i])
		builder.WriteString("_")
		builder.WriteString(palavras[i+1])
		builder.WriteString("_")
		builder.WriteString(palavras[i+2])
		tokens = append(tokens, builder.String())
	}
	
	return tokens
}

// buscarOOV tenta encontrar um token semelhante no vocabulário usando o índice por comprimento
func buscarOOV(token string) (TokenInfo, bool) {
	tokenLen := len(token)
	// Procurar apenas em tokens de comprimento igual ou ±1
	for delta := 0; delta <= 1; delta++ {
		for _, l := range []int{tokenLen + delta, tokenLen - delta} {
			if l <= 0 {
				continue
			}
			candidatos, ok := vocabPorComprimento[l]
			if !ok {
				continue
			}
			for _, cand := range candidatos {
				if levenshteinDist(token, cand) <= 1 {
					return Vocabulario[cand], true
				}
			}
		}
	}
	return TokenInfo{}, false
}

// levenshteinDist calcula a distância de edição com early exit
func levenshteinDist(s, t string) int {
	if s == t {
		return 0
	}
	ls, lt := len(s), len(t)
	if ls == 0 { return lt }
	if lt == 0 { return ls }
	
	// Usar apenas uma linha para economizar memória
	prev := make([]int, lt+1)
	curr := make([]int, lt+1)
	for j := 0; j <= lt; j++ {
		prev[j] = j
	}
	for i := 1; i <= ls; i++ {
		curr[0] = i
		for j := 1; j <= lt; j++ {
			cost := 1
			if s[i-1] == t[j-1] {
				cost = 0
			}
			curr[j] = prev[j] + 1 // delete
			if curr[j-1]+1 < curr[j] {
				curr[j] = curr[j-1] + 1 // insert
			}
			if prev[j-1]+cost < curr[j] {
				curr[j] = prev[j-1] + cost // substitute
			}
		}
		prev, curr = curr, prev
	}
	return prev[lt]
}

// GerarEmbeddingTFIDF gera o vetor TF-IDF e retorna métricas de [UNK]
func GerarEmbeddingTFIDF(texto string) ([]float32, []string, int) {
	tokens := ExtrairTokens(texto)
	vetor := make([]float32, DimensaoEmbedding)
	
	// Term Frequency
	freq := make(map[string]float64)
	for _, t := range tokens {
		freq[t]++
	}

	tokensDesconhecidos := []string{}
	tokensMapeados := 0

	// TF-IDF com Resgate de OOV via índice por comprimento
	for token, count := range freq {
		if info, ok := Vocabulario[token]; ok {
			vetor[info.Indice] += float32(count * info.IDF)
			tokensMapeados++
		} else {
			// Sub-word / Typo fallback (apenas para unigramas com len > 2)
			resgatado := false
			if !strings.Contains(token, "_") && len(token) > 2 {
				if info, found := buscarOOV(token); found {
					vetor[info.Indice] += float32(count * info.IDF * 0.5)
					tokensMapeados++
					resgatado = true
				}
			}
			if !resgatado {
				tokensDesconhecidos = append(tokensDesconhecidos, token)
			}
		}
	}

	// L2 Norm
	var somaQuadrados float32
	for _, v := range vetor {
		somaQuadrados += v * v
	}
	
	if somaQuadrados > 0 {
		norma := float32(math.Sqrt(float64(somaQuadrados)))
		for i := range vetor {
			vetor[i] /= norma
		}
	}
	
	return vetor, tokensDesconhecidos, tokensMapeados
}
