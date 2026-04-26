package lsh

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type MatrizPCAMulti struct {
	Metadados struct {
		Algoritmo       string `json:"algoritmo"`
		DimensaoEntrada int    `json:"dimensao_entrada"`
		DimensaoSaida   int    `json:"dimensao_saida"`
		TreinadoEm      string `json:"treinado_em"`
	} `json:"metadados"`
	Heads struct {
		Entidade [][]float64 `json:"entidade"`
		Contexto [][]float64 `json:"contexto"`
		Visual   [][]float64 `json:"visual"`
	} `json:"heads"`
}

var pcaEntidade [][]float32
var pcaContexto [][]float32
var pcaVisual [][]float32

func converterParaFloat32(matriz [][]float64, dimensao int) [][]float32 {
	res := make([][]float32, dimensao)
	for i, linha := range matriz {
		res[i] = make([]float32, len(linha))
		for j, val := range linha {
			res[i][j] = float32(val)
		}
	}
	return res
}

func CarregarMatrizPCA(caminho string) error {
	bytes, err := ioutil.ReadFile(caminho)
	if err != nil {
		return err
	}
	var dados MatrizPCAMulti
	if err := json.Unmarshal(bytes, &dados); err != nil {
		return err
	}

	dim := dados.Metadados.DimensaoSaida
	pcaEntidade = converterParaFloat32(dados.Heads.Entidade, dim)
	pcaContexto = converterParaFloat32(dados.Heads.Contexto, dim)
	pcaVisual = converterParaFloat32(dados.Heads.Visual, dim)
	
	log.Printf("[+] LSH Multi-Head Inicializado: %d hiperplanos por cabeça carregados.\n", dim)
	return nil
}

func gerarHash(embedding []float32, hiperplanos [][]float32) uint64 {
	var hash uint64 = 0
	for i, hiperplano := range hiperplanos {
		var dot float32 = 0.0
		for j := 0; j < len(embedding) && j < len(hiperplano); j++ {
			dot += embedding[j] * hiperplano[j]
		}
		if dot > 0 {
			hash |= (1 << i)
		}
	}
	return hash
}

// GerarSimHashPCA_Multi retorna as três cabeças de atenção simultaneamente
func GerarSimHashPCA_Multi(embedding []float32) (uint64, uint64, uint64) {
	if pcaEntidade == nil {
		log.Fatal("PCA Hiperplanos não carregados. Chame CarregarMatrizPCA() primeiro.")
	}
	hEntidade := gerarHash(embedding, pcaEntidade)
	hContexto := gerarHash(embedding, pcaContexto)
	hVisual := gerarHash(embedding, pcaVisual)
	
	return hEntidade, hContexto, hVisual
}
