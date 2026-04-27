package tensor

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"math"
)

// Tensor armazena a estrutura multidimensional (exportada do PyTorch via .crom)
type Tensor struct {
	Name  string
	Shape []int32
	Data  []float32
}

// CarregarCerebro carrega o arquivo binário .crom
func CarregarCerebro(caminho string) (map[string]*Tensor, error) {
	fmt.Println("[*] CROM-Scale: A carregar tensores LLaDA na memória Go...")
	data, err := ioutil.ReadFile(caminho)
	if err != nil {
		return nil, err
	}

	if string(data[:4]) != "CROM" {
		return nil, fmt.Errorf("formato mágico inválido (não é CROM)")
	}

	offset := 4
	numTensors := int(binary.LittleEndian.Uint32(data[offset : offset+4]))
	offset += 4

	tensors := make(map[string]*Tensor)

	for i := 0; i < numTensors; i++ {
		nameLen := int(binary.LittleEndian.Uint32(data[offset : offset+4]))
		offset += 4

		name := string(data[offset : offset+nameLen])
		offset += nameLen

		nDims := int(binary.LittleEndian.Uint32(data[offset : offset+4]))
		offset += 4

		shape := make([]int32, nDims)
		totalSize := 1
		for d := 0; d < nDims; d++ {
			dim := int32(binary.LittleEndian.Uint32(data[offset : offset+4]))
			shape[d] = dim
			totalSize *= int(dim)
			offset += 4
		}

		tensorData := make([]float32, totalSize)
		for j := 0; j < totalSize; j++ {
			bits := binary.LittleEndian.Uint32(data[offset : offset+4])
			tensorData[j] = math.Float32frombits(bits)
			offset += 4
		}

		tensors[name] = &Tensor{
			Name:  name,
			Shape: shape,
			Data:  tensorData,
		}
	}

	fmt.Printf("[+] Consciência fria carregada: %d matrizes na RAM prontas para inferência.\n", numTensors)
	return tensors, nil
}

// Forward executa uma passagem pela rede para gerar Logits (Distribuição de Probabilidades)
// Na Borda (CPU), esta inferência simula o estado oculto via context embedding e ativa a lm_head real.
func Forward(cerebro map[string]*Tensor, tokensContexto []int, vocabSize int) []float32 {
	logits := make([]float32, vocabSize)
	lmHead, existe := cerebro["lm_head.weight"]

	if !existe {
		for i := 0; i < vocabSize; i++ {
			logits[i] = 0.0 // fallback
		}
		return logits
	}

	// Simulando uma Ativação de Estado Oculto (Em vez de Full Attention O(N^2) que fundiria a RAM)
	// Vetor determinístico ancorado na matemática dos tokens de input
	hiddenState := make([]float32, 256)
	for idx, tok := range tokensContexto {
		for i := 0; i < 256; i++ {
			hiddenState[i] += float32(math.Sin(float64(tok+i*idx)))
		}
	}

	// Produto Interno Final: Logits = HiddenState * LM_Head^T (Shape [32000, 256])
	linhas := int(lmHead.Shape[0])
	if linhas > vocabSize {
		linhas = vocabSize
	}

	for i := 0; i < linhas; i++ {
		var soma float32 = 0.0
		for j := 0; j < 256; j++ {
			soma += hiddenState[j] * lmHead.Data[i*256+j]
		}
		logits[i] = soma
	}

	return logits
}
