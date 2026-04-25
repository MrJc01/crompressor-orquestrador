# Parte 2 — Arquitetura do Motor Crompressor

## 2.1 Visão Geral do Pipeline

O motor opera em duas direções: **Pack** (comprimir) e **Unpack** (descomprimir).

### Pipeline de Compressão (Pack)
```
Arquivo → Entropy Analysis → Chunker → LSH Search → Delta XOR → Zstd Pool → .crom
```

### Pipeline de Descompressão (Unpack)
```
.crom → Parse Header → Decompress Pool → Lookup Codebook → Apply Delta → Arquivo Original
```

## 2.2 Estrutura de Diretórios do Core

```
crompressor/
├── cmd/
│   └── crompressor/        # CLI principal (main.go)
│       └── main.go         # Entrypoint: crom pack/unpack/train/mount/...
│
├── internal/                # Lógica privada (inacessível externamente)
│   ├── autobrain/           # Seleção automática de estratégia
│   ├── chunker/             # Fixed, CDC, FastCDC, Semantic
│   ├── codebook/            # Leitura de .cromdb (mmap + WASM fallback)
│   ├── crypto/              # AES-256-GCM, Ed25519, Dilithium, Convergent
│   ├── delta/               # XOR, Diff/Patch, Zstd compress/decompress
│   ├── entropy/             # Shannon entropy calculator
│   ├── fractal/             # Geração determinística via polinômios
│   ├── merkle/              # Árvore Merkle para verificação de integridade
│   ├── metrics/             # Telemetria e contadores
│   ├── network/             # P2P (libp2p, Kademlia, Bitswap)
│   ├── remote/              # SSH fetcher para codebooks remotos
│   ├── search/              # LinearSearcher, LSHSearcher
│   ├── semantic/            # Análise semântica de conteúdo
│   ├── trainer/             # K-Means VQ training
│   └── vfs/                 # FUSE virtual filesystem (RandomReader)
│
├── pkg/                     # API pública (acessível por outros módulos)
│   ├── codebook/            # → internal/codebook (Open, OpenFromBytes)
│   ├── cromdb/              # Manipulação de .cromdb
│   ├── cromlib/             # Motor principal (Pack, Unpack, PackBytes, UnpackBytes)
│   │   └── vfs/             # VFS público
│   ├── crypto/              # → internal/crypto
│   ├── delta/               # → internal/delta
│   ├── entropy/             # → internal/entropy
│   ├── format/              # Header, ChunkEntry, Reader, Writer
│   ├── network/             # → internal/network
│   ├── sync/                # Manifest diff para sincronização
│   ├── trainer/             # → internal/trainer
│   └── vfs/                 # → internal/vfs
│
├── examples/                # Exemplos de uso
└── testdata/                # Dados de teste
```

## 2.3 O Motor de Compressão — `pkg/cromlib`

### Pack (arquivo para .crom)
Arquivo: `pkg/cromlib/packer.go`

```go
func Pack(inputPath, outputPath, codebookPath string, opts PackOptions) (*Metrics, error)
```

Fluxo interno:
1. **Entropy Analysis** — calcula Shannon do início do arquivo para escolher chunk size
2. **Abrir Codebook** — `codebook.Open(path)` via mmap (ou ReadFile em WASM)
3. **Criar LSH Searcher** — constrói índice de busca rápida
4. **Chunking** — divide dados em chunks (tamanho adaptativo baseado na entropia)
5. **Para cada chunk**:
   - Calcula entropia local
   - Se alta entropia (>3.0 em archive mode): armazena como literal
   - Senão: busca melhor match via LSH
   - Se match bom (>20% similaridade): XOR delta com padrão
   - Se match excelente (>95%): conta como "hit"
6. **Compress Pool** — comprime todos os deltas/literais com Zstd
7. **Write .crom** — header + block table + chunk table + compressed pool

### PackBytes (in-memory para WASM)
Arquivo: `pkg/cromlib/pack_bytes.go`

```go
func PackBytes(input []byte, codebookData []byte, opts PackOptions) ([]byte, *Metrics, error)
```

Mesma lógica, mas sem filesystem. Recebe e retorna `[]byte`.

### Unpack (arquivo .crom para original)
Arquivo: `pkg/cromlib/unpacker.go`

```go
func Unpack(inputPath, outputPath, codebookPath string, opts UnpackOptions) error
```

Fluxo interno:
1. **Parse Header** — versão, hash original, chunk count
2. **Stream Blocks** — para cada bloco:
   - Decrypt (se criptografado)
   - Decompress (Zstd)
   - Para cada chunk entry: lookup codebook + apply delta
3. **Verify SHA-256** — compara hash reconstruído com o original

### Modos Especiais
- **Passthrough**: arquivos de alta entropia (>7.0) são armazenados sem compressão
- **Fractal**: chunks gerados via polinômios determinísticos (sem codebook)
- **Variational (Fuzziness)**: modo lossy que escolhe codewords vizinhos aleatoriamente

## 2.4 O Codebook (.cromdb)

### Estrutura Binária
```
┌──────────────────────────┐
│    Header (32 bytes)     │  Magic, version, codeword_count, codeword_size
├──────────────────────────┤
│   Codeword 0 (N bytes)  │
│   Codeword 1 (N bytes)  │
│   ...                    │
│   Codeword K (N bytes)   │
└──────────────────────────┘
```

### Treinamento
```go
// internal/trainer/trainer.go
func Train(opts TrainOptions) (*TrainResult, error)
```

1. Coleta amostras de dados do domínio
2. Divide em chunks
3. Executa K-Means iterativo (convergência por distorção mínima)
4. Salva centroides como codewords

### Acesso em Runtime
- **Linux/macOS**: `mmap` — mapeia arquivo direto na memória virtual (zero-copy)
- **WASM**: `ReadFile` + buffer — carrega inteiro na memória JS

## 2.5 O Sistema de Criptografia

```
internal/crypto/
├── aes.go          # AES-256-GCM (encrypt/decrypt)
├── convergent.go   # Convergent encryption (dedup-safe)
├── dilithium.go    # Post-quantum signatures
├── key.go          # PBKDF2 key derivation
└── sign.go         # Ed25519 signatures
```

- **Pack com criptografia**: cada bloco é criptografado individualmente com AES-256-GCM
- **Convergent Encryption**: o hash do conteúdo é a chave → arquivos idênticos geram o mesmo ciphertext (permite deduplicação mesmo criptografado)
- **Dilithium**: assinaturas pós-quânticas para verificação em rede P2P

## 2.6 O VFS (Virtual Filesystem)

```go
// internal/vfs/
func Mount(cromFile, codebookPath, mountPoint, password string, maxMB int) error
```

Monta um arquivo `.crom` como um sistema de arquivos FUSE:
- Acesso aleatório O(1) via `RandomReader`
- LRU cache de blocos descomprimidos
- Suporta leitura parcial (streaming de vídeo, seek em PDF)

## 2.7 CLI — Comandos Disponíveis

```bash
crompressor pack    <input> <output.crom> <codebook.cromdb>   # Comprimir
crompressor unpack  <input.crom> <output> <codebook.cromdb>   # Descomprimir
crompressor train   <data_dir> <output.cromdb>                 # Treinar codebook
crompressor mount   <input.crom> <mountpoint> <codebook>       # Montar VFS FUSE
crompressor info    <file.crom>                                # Mostrar metadados
crompressor bench   <input> <codebook.cromdb>                  # Benchmark
```

---

**Próximo**: [Parte 3 — O Ecossistema](./parte3_ecossistema.md)
