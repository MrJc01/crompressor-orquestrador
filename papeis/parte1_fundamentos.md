# Parte 1 — Fundamentos Teóricos da Compressão CROM

## 1.1 O Problema da Compressão

Todo arquivo digital é uma sequência de bytes. Comprimir é encontrar uma representação menor que preserva a informação original. Existem dois limites fundamentais:

- **Lossless** (sem perdas): o arquivo original pode ser reconstruído perfeitamente
- **Lossy** (com perdas): aceita-se alguma degradação em troca de maior compressão

O Crompressor opera em **ambos os modos**: `archive` (lossless) e `edge` (lossy).

## 1.2 Entropia de Shannon

Claude Shannon provou em 1948 que existe um **limite teórico** para compressão sem perdas. A entropia H mede quanta informação real existe nos dados:

```
H = -Σ p(x) · log₂(p(x))    para cada símbolo x
```

- **H = 0 bits/byte**: dados totalmente uniformes (ex: arquivo de zeros) → compressão infinita
- **H = 8 bits/byte**: dados perfeitamente aleatórios → incompressível
- **Texto típico**: H ≈ 4-5 bits/byte → pode ser comprimido ~50%

**No Crompressor**: a função `entropy.Shannon(data)` calcula isso para cada chunk, decidindo automaticamente se vale a pena comprimir ou armazenar como literal.

## 1.3 Compressão Tradicional vs. CROM

### Algoritmos Tradicionais (gzip, zstd, brotli)
Usam **LZ77/LZ78**: procuram repetições locais dentro de uma janela deslizante.

```
Entrada:  "ABCABCABC"
LZ77:     "ABC" + referência(offset=3, len=6)
```

**Limitação**: só encontram repetições exatas dentro de uma janela limitada (32KB-128KB).

### Abordagem CROM — Compressão por Codebook
O CROM não procura repetições locais. Em vez disso:

1. **Treina** um codebook (dicionário) com padrões aprendidos via K-Means/VQ
2. **Busca** para cada chunk o padrão mais similar no codebook
3. **Armazena** apenas o ID do padrão + a diferença (delta XOR)

```
Chunk Original:   [0x48, 0x65, 0x6C, 0x6C, 0x6F]
Padrão Codebook:  [0x48, 0x65, 0x6C, 0x70, 0x21]  (ID: 42)
Delta XOR:        [0x00, 0x00, 0x00, 0x1C, 0x4E]
Armazenado:       {ID: 42, delta: [0x00, 0x00, 0x00, 0x1C, 0x4E]}
```

**Vantagem**: o codebook captura padrões **globais** do domínio, não apenas repetições locais.

## 1.4 Vector Quantization (VQ)

O coração do CROM é a **Quantização Vetorial**. Cada chunk de dados é tratado como um vetor em espaço N-dimensional (onde N = tamanho do chunk).

O treinamento do codebook funciona assim:

1. Coletar milhares de chunks de dados representativos
2. Executar K-Means clustering para encontrar os K centroides
3. Cada centroide vira um "codeword" no codebook
4. Ao comprimir, para cada chunk, encontrar o codeword mais próximo

```
Codebook com K=65536 codewords de 128 bytes cada:
- Codeword 0:     padrão típico de header HTTP
- Codeword 1:     padrão de código Go (func main)
- Codeword 42:    padrão de JSON com aspas
- ...
- Codeword 65535: padrão de binário ELF
```

## 1.5 LSH — Locality-Sensitive Hashing

Buscar o codeword mais próximo entre 65.536 candidatos é O(N×K) — lento. O CROM usa **LSH** para reduzir isso a O(1) amortizado:

1. Gera múltiplos hashes aleatórios para cada codeword
2. Armazena em hash tables
3. Para buscar, gera os mesmos hashes do chunk e consulta as tabelas
4. Candidatos com mais colisões de hash são provavelmente os mais similares

Arquivo: `internal/search/lsh.go`

## 1.6 Content-Defined Chunking (CDC)

Em vez de dividir o arquivo em blocos fixos de 128 bytes, o CDC usa uma janela deslizante com rolling hash para encontrar fronteiras "naturais" no conteúdo:

```
Inserir 1 byte no início do arquivo:
- Fixed Chunking: TODOS os chunks mudam (desalinhamento total)
- CDC: apenas o primeiro chunk muda, os demais permanecem iguais
```

Isso é crítico para sincronização P2P (crompressor-sync): apenas chunks realmente modificados são retransmitidos.

Implementações no core:
- `internal/chunker/fixed.go` — Tamanho fixo (mais rápido)
- `internal/chunker/cdc.go` — Content-Defined (mais eficiente para sync)
- `internal/chunker/fastcdc.go` — FastCDC otimizado
- `internal/chunker/acac.go` — Semântico (respeita delimitadores como `\n`)

## 1.7 O Formato .crom

O arquivo `.crom` é um container binário com a seguinte estrutura:

```
┌─────────────────────────────┐
│         Header (64B)        │  Versão, hash original, contagem de chunks
├─────────────────────────────┤
│       Block Table (var)     │  Tamanhos de cada bloco comprimido
├─────────────────────────────┤
│      Chunk Table (var)      │  Entradas: {CodebookID, DeltaOffset, DeltaSize}
├─────────────────────────────┤
│   Compressed Delta Pool     │  Deltas XOR comprimidos com Zstd
└─────────────────────────────┘
```

Definido em `pkg/format/` — suporta V1 (legacy) e V2+ (block-based streaming).

## 1.8 Resumo dos Conceitos-Chave

| Conceito | Onde no Código | Para que serve |
|----------|----------------|----------------|
| Entropia Shannon | `internal/entropy/` | Decidir se chunk é compressível |
| Vector Quantization | `internal/trainer/` | Treinar codebook |
| LSH Search | `internal/search/lsh.go` | Busca O(1) no codebook |
| CDC Chunking | `internal/chunker/` | Dividir dados em chunks |
| XOR Delta | `internal/delta/` | Diferença entre chunk e padrão |
| Zstd Pool | `internal/delta/compress.go` | Comprimir pool de deltas |
| .crom Format | `pkg/format/` | Container binário |
| AES-256-GCM | `internal/crypto/` | Criptografia opcional |

---

**Próximo**: [Parte 2 — Arquitetura do Motor](./parte2_arquitetura.md)
