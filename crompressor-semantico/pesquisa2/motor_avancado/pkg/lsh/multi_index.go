package lsh

import (
	"sync"
	"github.com/MrJc01/crompressor-semantico/pkg/hamming"
)

// MultiIndex implementa o "Estado da Arte" em busca LSH para bilhões de registros.
// Quebra um hash de 64 bits em 4 blocos de 16 bits.
// Permite encontrar vizinhos com pequenas distâncias de Hamming em tempo real,
// pois se 2 hashes têm apenas 3 bits de diferença, pelo menos 1 dos 4 blocos de 16 bits será IDÊNTICO.
type MultiIndex struct {
	mu      sync.RWMutex
	Tabelas [4]map[uint16][]uint64
}

func NovoMultiIndex() *MultiIndex {
	mi := &MultiIndex{}
	for i := 0; i < 4; i++ {
		mi.Tabelas[i] = make(map[uint16][]uint64)
	}
	return mi
}

// segmentar divide o uint64 em 4 blocos de 16 bits
func segmentar(hash uint64) [4]uint16 {
	return [4]uint16{
		uint16(hash >> 48),
		uint16((hash >> 32) & 0xFFFF),
		uint16((hash >> 16) & 0xFFFF),
		uint16(hash & 0xFFFF),
	}
}

// Inserir guarda o hash em 4 Buckets diferentes simultaneamente.
func (mi *MultiIndex) Inserir(hash uint64) {
	mi.mu.Lock()
	defer mi.mu.Unlock()

	segmentos := segmentar(hash)
	for i := 0; i < 4; i++ {
		mi.Tabelas[i][segmentos[i]] = append(mi.Tabelas[i][segmentos[i]], hash)
	}
}

// BuscaRapida procura candidatos que compartilhem pelo menos 1 bloco de 16 bits exato com o alvo.
// Depois calcula a distância real de Hamming apenas para esses candidatos seletos.
func (mi *MultiIndex) BuscaRapida(hashAlvo uint64, limiarBits int) (uint64, bool) {
	mi.mu.RLock()
	defer mi.mu.RUnlock()

	segmentos := segmentar(hashAlvo)
	candidatosVerificados := make(map[uint64]bool)

	// Busca nas 4 tabelas
	for i := 0; i < 4; i++ {
		if lista, existe := mi.Tabelas[i][segmentos[i]]; existe {
			for _, candidato := range lista {
				// Evita recalcular a distância se o mesmo candidato já foi encontrado por outra tabela
				if !candidatosVerificados[candidato] {
					if hamming.Distance(hashAlvo, candidato) <= limiarBits {
						return candidato, true // Deduplicação/Match!
					}
					candidatosVerificados[candidato] = true
				}
			}
		}
	}
	return 0, false
}
