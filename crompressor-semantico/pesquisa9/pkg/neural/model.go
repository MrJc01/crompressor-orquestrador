package neural

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"
	"strings"
)

type Config struct {
	EmbeddingDim   int
	HiddenDim      int
	ContextRadius  int
	LearningRate   float64
	Epochs         int
	MaskRatio      float64
	DiffusionSteps int
	Temperature    float64
	TopK           int
	MaxTokens      int
}

func DefaultConfig() Config {
	return Config{
		EmbeddingDim: 64, HiddenDim: 192, ContextRadius: 4,
		LearningRate: 0.005, Epochs: 300, MaskRatio: 0.4,
		DiffusionSteps: 10, Temperature: 0.5, TopK: 8, MaxTokens: 12,
	}
}

const MaskToken = "[MASK]"

type Model struct {
	Config    Config
	Vocab     []string
	WordToIdx map[string]int
	IdxToWord map[int]string
	MaskIdx   int

	Emb, PosEmb [][]float64
	W1          [][]float64
	B1          []float64
	W2          [][]float64
	B2          []float64
	LossHistory []float64
}

func NewModel(c Config) *Model {
	return &Model{Config: c, WordToIdx: map[string]int{}, IdxToWord: map[int]string{}}
}

func (m *Model) BuildVocab(corpus string) {
	seen := map[string]bool{MaskToken: true}
	m.Vocab = []string{MaskToken}
	for _, w := range strings.Fields(strings.ToLower(corpus)) {
		if !seen[w] {
			seen[w] = true
			m.Vocab = append(m.Vocab, w)
		}
	}
	sort.Strings(m.Vocab[1:])
	for i, w := range m.Vocab {
		m.WordToIdx[w] = i
		m.IdxToWord[i] = w
	}
	m.MaskIdx = 0
	fmt.Printf("[+] Vocabulário: %d palavras\n", len(m.Vocab))
}

func (m *Model) InitWeights() {
	V := len(m.Vocab)
	E := m.Config.EmbeddingDim
	H := m.Config.HiddenDim
	R := m.Config.ContextRadius
	W := 2*R + 1
	inputDim := W * E

	m.Emb = xm(V, E, E)
	m.PosEmb = xm(W, E, E)
	m.W1 = xm(inputDim, H, inputDim)
	m.B1 = make([]float64, H)
	m.W2 = xm(H, V, H)
	m.B2 = make([]float64, V)

	total := V*E + W*E + inputDim*H + H + H*V + V
	fmt.Printf("[+] Parâmetros: %d | Window=%d | Emb[%d×%d] → Hidden[%d] → Out[%d]\n",
		total, W, V, E, H, V)
}

func xm(r, c, fan int) [][]float64 {
	s := math.Sqrt(2.0 / float64(fan))
	mat := make([][]float64, r)
	for i := range mat {
		mat[i] = make([]float64, c)
		for j := range mat[i] {
			mat[i][j] = rand.NormFloat64() * s
		}
	}
	return mat
}

func (m *Model) forward(tokens []int, pos int) ([]float64, []float64, []float64) {
	E := m.Config.EmbeddingDim
	H := m.Config.HiddenDim
	R := m.Config.ContextRadius
	V := len(m.Vocab)
	W := 2*R + 1

	input := make([]float64, W*E)
	for w := 0; w < W; w++ {
		srcPos := pos - R + w
		var tokIdx int
		if srcPos < 0 || srcPos >= len(tokens) {
			tokIdx = m.MaskIdx
		} else {
			tokIdx = tokens[srcPos]
		}
		for j := 0; j < E; j++ {
			input[w*E+j] = m.Emb[tokIdx][j] + m.PosEmb[w][j]
		}
	}

	hidden := make([]float64, H)
	for j := 0; j < H; j++ {
		s := m.B1[j]
		for i := 0; i < len(input); i++ {
			s += input[i] * m.W1[i][j]
		}
		if s > 0 {
			hidden[j] = s
		}
	}

	logits := make([]float64, V)
	for j := 0; j < V; j++ {
		s := m.B2[j]
		for i := 0; i < H; i++ {
			s += hidden[i] * m.W2[i][j]
		}
		logits[j] = s
	}
	return logits, hidden, input
}

func softmax(logits []float64) []float64 {
	mx := logits[0]
	for _, v := range logits[1:] {
		if v > mx {
			mx = v
		}
	}
	p := make([]float64, len(logits))
	s := 0.0
	for i, v := range logits {
		p[i] = math.Exp(v - mx)
		s += p[i]
	}
	for i := range p {
		p[i] /= s
	}
	return p
}

func clip(g, mx float64) float64 {
	if g > mx {
		return mx
	}
	if g < -mx {
		return -mx
	}
	return g
}

func (m *Model) trainStep(tokens []int, pos, target int) float64 {
	V := len(m.Vocab)
	E := m.Config.EmbeddingDim
	H := m.Config.HiddenDim
	R := m.Config.ContextRadius
	lr := m.Config.LearningRate
	W := 2*R + 1
	const MG = 5.0

	logits, hidden, input := m.forward(tokens, pos)
	probs := softmax(logits)
	loss := -math.Log(math.Max(probs[target], 1e-10))

	dL := make([]float64, V)
	copy(dL, probs)
	dL[target] -= 1.0

	dH := make([]float64, H)
	for i := 0; i < H; i++ {
		for j := 0; j < V; j++ {
			dH[i] += dL[j] * m.W2[i][j]
			m.W2[i][j] -= lr * clip(dL[j]*hidden[i], MG)
		}
	}
	for j := 0; j < V; j++ {
		m.B2[j] -= lr * clip(dL[j], MG)
	}

	for i := 0; i < H; i++ {
		if hidden[i] <= 0 {
			dH[i] = 0
		}
	}

	dInput := make([]float64, len(input))
	for i := 0; i < len(input); i++ {
		for j := 0; j < H; j++ {
			dInput[i] += dH[j] * m.W1[i][j]
			m.W1[i][j] -= lr * clip(dH[j]*input[i], MG)
		}
	}
	for j := 0; j < H; j++ {
		m.B1[j] -= lr * clip(dH[j], MG)
	}

	for w := 0; w < W; w++ {
		srcPos := pos - R + w
		var tokIdx int
		if srcPos < 0 || srcPos >= len(tokens) {
			tokIdx = m.MaskIdx
		} else {
			tokIdx = tokens[srcPos]
		}
		for j := 0; j < E; j++ {
			g := clip(dInput[w*E+j], MG)
			m.Emb[tokIdx][j] -= lr * g
			m.PosEmb[w][j] -= lr * g
		}
	}
	return loss
}

func (m *Model) Train(corpus string) {
	sents := splitSentences(corpus)
	baseLR := m.Config.LearningRate
	fmt.Printf("[*] Corpus: %d frases | %d épocas | mask=%.0f%% | lr=%.4f→%.5f (cosine)\n",
		len(sents), m.Config.Epochs, m.Config.MaskRatio*100, baseLR, baseLR*0.1)

	m.LossHistory = nil
	for epoch := 0; epoch < m.Config.Epochs; epoch++ {
		progress := float64(epoch) / float64(m.Config.Epochs)
		m.Config.LearningRate = baseLR * (0.1 + 0.9*0.5*(1.0+math.Cos(math.Pi*progress)))

		rand.Shuffle(len(sents), func(i, j int) { sents[i], sents[j] = sents[j], sents[i] })
		totalLoss := 0.0
		count := 0

		for _, sent := range sents {
			words := strings.Fields(strings.ToLower(sent))
			if len(words) < 3 {
				continue
			}
			idxs := make([]int, len(words))
			ok := true
			for i, w := range words {
				idx, has := m.WordToIdx[w]
				if !has {
					ok = false
					break
				}
				idxs[i] = idx
			}
			if !ok {
				continue
			}

			masked := make([]int, len(idxs))
			copy(masked, idxs)
			for i := range masked {
				if rand.Float64() < m.Config.MaskRatio {
					masked[i] = m.MaskIdx
				}
			}

			for i := range masked {
				if masked[i] == m.MaskIdx {
					totalLoss += m.trainStep(masked, i, idxs[i])
					count++
				}
			}
		}
		if count > 0 {
			avg := totalLoss / float64(count)
			m.LossHistory = append(m.LossHistory, avg)
			if epoch%25 == 0 || epoch == m.Config.Epochs-1 {
				bar := strings.Repeat("█", int(math.Min(avg*5, 40)))
				fmt.Printf("    Época %3d/%d | Loss: %.4f | lr=%.5f |%s\n",
					epoch+1, m.Config.Epochs, avg, m.Config.LearningRate, bar)
			}
		}
	}
	m.Config.LearningRate = baseLR
	if len(m.LossHistory) > 0 {
		fmt.Printf("[+] Treino concluído! Loss: %.4f → %.4f\n",
			m.LossHistory[0], m.LossHistory[len(m.LossHistory)-1])
	}
}

func (m *Model) Generate(prompt string, n int) string {
	if n <= 0 {
		n = m.Config.MaxTokens
	}
	pw := strings.Fields(strings.ToLower(prompt))
	seq := make([]int, len(pw)+n)
	for i, w := range pw {
		idx, ok := m.WordToIdx[w]
		if ok {
			seq[i] = idx
		} else {
			seq[i] = m.MaskIdx
		}
	}
	for i := len(pw); i < len(seq); i++ {
		seq[i] = m.MaskIdx
	}

	for step := 0; step < m.Config.DiffusionSteps; step++ {
		var maskedPos []int
		for i := len(pw); i < len(seq); i++ {
			if seq[i] == m.MaskIdx {
				maskedPos = append(maskedPos, i)
			}
		}
		if len(maskedPos) == 0 {
			break
		}

		type Pred struct {
			Pos, Tok int
			Conf     float64
		}
		preds := make([]Pred, len(maskedPos))
		for pi, pos := range maskedPos {
			logits, _, _ := m.forward(seq, pos)
			for i := range logits {
				logits[i] /= m.Config.Temperature
			}
			probs := softmax(logits)
			tok := m.sampleTopK(probs)
			preds[pi] = Pred{pos, tok, probs[tok]}
		}

		sort.Slice(preds, func(i, j int) bool { return preds[i].Conf > preds[j].Conf })
		reveal := int(math.Ceil(float64(len(preds)) * float64(step+1) / float64(m.Config.DiffusionSteps)))
		if reveal > len(preds) {
			reveal = len(preds)
		}
		for i := 0; i < reveal; i++ {
			seq[preds[i].Pos] = preds[i].Tok
		}
	}

	var result []string
	for i := len(pw); i < len(seq); i++ {
		w := m.IdxToWord[seq[i]]
		if w != MaskToken {
			result = append(result, w)
		}
	}
	return strings.Join(result, " ")
}

func (m *Model) sampleTopK(probs []float64) int {
	K := m.Config.TopK
	if K <= 0 || K >= len(probs) {
		K = len(probs)
	}
	type IP struct {
		I int
		P float64
	}
	ps := make([]IP, len(probs))
	for i, p := range probs {
		ps[i] = IP{i, p}
	}
	sort.Slice(ps, func(i, j int) bool { return ps[i].P > ps[j].P })
	top := ps[:K]
	sum := 0.0
	for _, p := range top {
		sum += p.P
	}
	r := rand.Float64() * sum
	cum := 0.0
	for _, p := range top {
		cum += p.P
		if r <= cum {
			return p.I
		}
	}
	return top[0].I
}

func (m *Model) GetVocabSize() int { return len(m.Vocab) }

func splitSentences(corpus string) []string {
	corpus = strings.ReplaceAll(corpus, ".", ".\n")
	corpus = strings.ReplaceAll(corpus, "!", "!\n")
	corpus = strings.ReplaceAll(corpus, "?", "?\n")
	var r []string
	for _, l := range strings.Split(corpus, "\n") {
		l = strings.TrimSpace(l)
		if len(strings.Fields(l)) >= 3 {
			r = append(r, l)
		}
	}
	return r
}

// ---- Save/Load de Pesos Treinados ----

// Save serializa vocab + pesos em formato binário
func (m *Model) Save(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)

	// Magic + Config
	binary.Write(w, binary.LittleEndian, []byte("CROM"))
	binary.Write(w, binary.LittleEndian, int32(m.Config.EmbeddingDim))
	binary.Write(w, binary.LittleEndian, int32(m.Config.HiddenDim))
	binary.Write(w, binary.LittleEndian, int32(m.Config.ContextRadius))

	// Vocab
	binary.Write(w, binary.LittleEndian, int32(len(m.Vocab)))
	for _, word := range m.Vocab {
		b := []byte(word)
		binary.Write(w, binary.LittleEndian, int32(len(b)))
		w.Write(b)
	}

	// Pesos
	writeMatrix(w, m.Emb)
	writeMatrix(w, m.PosEmb)
	writeMatrix(w, m.W1)
	writeVec(w, m.B1)
	writeMatrix(w, m.W2)
	writeVec(w, m.B2)

	w.Flush()
	info, _ := f.Stat()
	fmt.Printf("[+] Pesos salvos: %s (%.1f KB)\n", path, float64(info.Size())/1024)
	return nil
}

// Load carrega vocab + pesos de ficheiro binário
func (m *Model) Load(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	r := bufio.NewReader(f)

	// Magic
	magic := make([]byte, 4)
	r.Read(magic)
	if string(magic) != "CROM" {
		return fmt.Errorf("formato inválido (magic: %s)", string(magic))
	}

	// Config
	var e, h, cr int32
	binary.Read(r, binary.LittleEndian, &e)
	binary.Read(r, binary.LittleEndian, &h)
	binary.Read(r, binary.LittleEndian, &cr)
	m.Config.EmbeddingDim = int(e)
	m.Config.HiddenDim = int(h)
	m.Config.ContextRadius = int(cr)

	// Vocab
	var vSize int32
	binary.Read(r, binary.LittleEndian, &vSize)
	m.Vocab = make([]string, vSize)
	m.WordToIdx = make(map[string]int)
	m.IdxToWord = make(map[int]string)
	for i := 0; i < int(vSize); i++ {
		var bLen int32
		binary.Read(r, binary.LittleEndian, &bLen)
		b := make([]byte, bLen)
		r.Read(b)
		m.Vocab[i] = string(b)
		m.WordToIdx[string(b)] = i
		m.IdxToWord[i] = string(b)
	}
	m.MaskIdx = m.WordToIdx[MaskToken]

	// Pesos
	W := 2*m.Config.ContextRadius + 1
	V := len(m.Vocab)
	E := m.Config.EmbeddingDim
	H := m.Config.HiddenDim

	m.Emb = readMatrix(r, V, E)
	m.PosEmb = readMatrix(r, W, E)
	m.W1 = readMatrix(r, W*E, H)
	m.B1 = readVec(r, H)
	m.W2 = readMatrix(r, H, V)
	m.B2 = readVec(r, V)

	fmt.Printf("[+] Pesos carregados: %s | vocab=%d params=%d\n", path,
		len(m.Vocab), V*E+W*E+(W*E)*H+H+H*V+V)
	return nil
}

func writeMatrix(w *bufio.Writer, mat [][]float64) {
	for _, row := range mat {
		for _, v := range row {
			binary.Write(w, binary.LittleEndian, float32(v))
		}
	}
}

func writeVec(w *bufio.Writer, vec []float64) {
	for _, v := range vec {
		binary.Write(w, binary.LittleEndian, float32(v))
	}
}

func readMatrix(r *bufio.Reader, rows, cols int) [][]float64 {
	mat := make([][]float64, rows)
	for i := range mat {
		mat[i] = make([]float64, cols)
		for j := range mat[i] {
			var v float32
			binary.Read(r, binary.LittleEndian, &v)
			mat[i][j] = float64(v)
		}
	}
	return mat
}

func readVec(r *bufio.Reader, size int) []float64 {
	vec := make([]float64, size)
	for i := range vec {
		var v float32
		binary.Read(r, binary.LittleEndian, &v)
		vec[i] = float64(v)
	}
	return vec
}
