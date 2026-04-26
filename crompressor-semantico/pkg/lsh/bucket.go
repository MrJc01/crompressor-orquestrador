package lsh

import (
	"fmt"
	"sync"
	
	"github.com/MrJc01/crompressor-semantico/pkg/hamming"
)

// DicionarioBuckets é o núcleo do Motor de Busca O(1).
// Em vez de iterar sobre 1 milhão de hashes, usamos um Multi-Index Hashing.
// A chave (uint16) representa os primeiros N bits do Hash (Prefixo).
// Assim, isolamos a busca em um pequeno "Bucket", garantindo alta velocidade.
type DicionarioBuckets struct {
	mu      sync.RWMutex
	Buckets map[uint16][]uint64
	
	// Configurável: quantos bits do prefixo vamos usar como chave (ex: 8 bits = 256 buckets, 16 bits = 65536 buckets)
	PrefixoBits int
}

func NovoDicionario(prefixoBits int) *DicionarioBuckets {
	return &DicionarioBuckets{
		Buckets:     make(map[uint16][]uint64),
		PrefixoBits: prefixoBits,
	}
}

// extrairPrefixo pega os X primeiros bits do Hash de 64 bits para usar como chave do Bucket.
func (d *DicionarioBuckets) extrairPrefixo(hash uint64) uint16 {
	// Deslocamos para a direita (64 - N) bits para ficar apenas com os N bits mais significativos
	deslocamento := 64 - d.PrefixoBits
	return uint16(hash >> deslocamento)
}

// Inserir guarda o hash no seu Bucket correspondente.
func (d *DicionarioBuckets) Inserir(hash uint64) {
	d.mu.Lock()
	defer d.mu.Unlock()
	
	chave := d.extrairPrefixo(hash)
	d.Buckets[chave] = append(d.Buckets[chave], hash)
}

// BuscarSimilar procura dentro do Bucket correspondente por um hash que esteja
// dentro da distância de Hamming aceitável (limiar).
func (d *DicionarioBuckets) BuscarSimilar(hashAlvo uint64, limiar int) (uint64, bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	
	chave := d.extrairPrefixo(hashAlvo)
	candidatos, existe := d.Buckets[chave]
	
	if !existe {
		return 0, false
	}
	
	// Itera apenas sobre a fração ínfima de hashes que estão no mesmo Bucket
	for _, candidato := range candidatos {
		if hamming.Distance(hashAlvo, candidato) <= limiar {
			// Deduplicação! Encontramos um similar.
			return candidato, true
		}
	}
	
	return 0, false
}

// Densidade debuga quantos hashes existem num determinado bucket
func (d *DicionarioBuckets) DebugDensidade() {
	d.mu.RLock()
	defer d.mu.RUnlock()
	fmt.Printf("[Bucketing] Dicionário ativo com %d PrefixBits (Max %d Buckets). Buckets ocupados: %d\n", 
		d.PrefixoBits, 1<<d.PrefixoBits, len(d.Buckets))
}
