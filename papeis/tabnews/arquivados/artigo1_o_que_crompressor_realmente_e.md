# O que o Crompressor realmente é (e o que ele NÃO é) — Benchmarks reais, código aberto, sem buzzwords

**TL;DR**: Publiquei 3 artigos sobre o Crompressor com linguagem exagerada. Vocês tinham razão em questionar. Aqui estão os benchmarks reais, honestos, com dados do meu próprio sistema. O CROM não bate o ZSTD em compressão pura — mas faz coisas que ZSTD não faz.

---

## Antes de tudo: vocês tinham razão

Nos meus artigos anteriores, usei termos como "entropia do universo", "computação termodinâmica" e "compilação de realidade". O @wsobrinho resumiu perfeitamente:

> "Existe uma centelha técnica aí, mas hoje o texto está maior que a prova. Abaixe o volume das promessas e aumente o peso das evidências."

Aceito. Este artigo é a correção. Sem poesia, sem metáforas cósmicas. Só código, números e honestidade.

---

## O que o Crompressor é, em uma frase

É um compressor lossless baseado em **codebook** (dicionário pré-treinado). Em vez de procurar repetições locais dentro do arquivo (como gzip/zstd fazem), ele procura padrões globais em um dicionário externo (.cromdb) e armazena apenas a diferença (XOR delta).

```
Compressão tradicional (ZSTD):
  arquivo → encontra repetições locais → substitui por referências

Compressão CROM:
  arquivo → divide em chunks de 128 bytes
           → busca cada chunk no codebook (dicionário treinado)
           → armazena: ID do padrão + diferença XOR
```

## Benchmarks reais — dados do meu sistema

Codebook de 16.384 codewords, treinado nos próprios dados.

| Dataset | Original | GZIP-9 | ZSTD-19 | CROM |
|---------|----------|--------|---------|------|
| Código Go (todo o CROM, 2.5 MB) | 100% | 27.5% | **10.1%** | 30.8% |
| Logs do sistema (6.9 MB) | 100% | 3.8% | **2.3%** | 22.4% |
| Documentação Markdown (5.8 MB) | 100% | 29.3% | **20.4%** | 42.6% |
| go.sum hashes (446 KB) | 100% | 38.1% | **7.8%** | 27.5% |
| dmesg kernel (260 KB) | 100% | 20.1% | **12.3%** | 40.3% |

**Veredicto**: ZSTD ganha em compressão pura. Sem exceção. Em todos os cenários testados.

**Hit Rate do codebook**: 0% — nenhum chunk encontrou match acima de 20% de similaridade no codebook.

Sim, você leu certo. Zero por cento.

---

## Então por que continuar?

Porque **compressão pura não é a proposta do CROM**. A proposta é um conjunto de capacidades que nenhum compressor tradicional oferece:

### 1. Acesso aleatório O(1) — VFS/FUSE
Montar um .crom como sistema de arquivos virtual. Ler byte 500.000 de um arquivo de 10GB sem descomprimir os outros 9.999.999.500 bytes.

```bash
crompressor mount arquivo.crom /mnt/virtual codebook.cromdb
cat /mnt/virtual/arquivo.txt  # lê só o bloco necessário
```

Gzip não faz isso. Zstd não faz isso. Você precisa descomprimir tudo primeiro.

### 2. Sincronização P2P — só transferir o que mudou
Dois nós que compartilham o mesmo codebook podem sincronizar arquivos enviando apenas os **IDs de chunks que mudaram**, não os bytes:

```
Arquivo original: 1000 chunks
Versão modificada: 5 chunks diferentes
Transferência CROM: 5 × (ID + delta) ≈ poucos KB
Transferência tradicional: re-enviar arquivo inteiro
```

### 3. Soberania via codebook
O codebook é um **shared secret**. Sem ele, o arquivo .crom é incompreensível. Não é criptografia (existe AES-256 separado para isso), mas é uma camada de obscuridade que tem valor prático:

- Data center A treina codebook nos seus dados
- Envia .crom para data center B (que tem o mesmo codebook)
- Um intermediário que intercepta o .crom não consegue reconstruir sem o codebook

### 4. Criptografia convergente
Mesmo criptografado, dois arquivos idênticos geram o mesmo ciphertext. Isso permite **deduplicação em dados criptografados** — algo que ZSTD + AES separados não conseguem.

---

## O problema real: Hit Rate 0%

O motor usa Hamming Distance (contagem de bits diferentes) para medir similaridade. Com texto UTF-8, mesmo chunks "parecidos" têm muitos bits divergentes. O resultado: tudo cai como literal.

**Isso é um bug ou uma limitação de design?** É uma limitação. O codebook-based approach funciona melhor com:
- Dados binários com padrões exatos (protocolos, headers, structs)
- Domínios ultra-específicos onde o codebook "conhece" os padrões
- Múltiplos arquivos do mesmo domínio (cross-file dedup)

**Próximo passo técnico**: investigar métricas de similaridade melhores que Hamming puro (ex: byte-level Jaccard, cosseno sobre n-grams).

---

## Como rodar você mesmo

```bash
# Instalar
git clone https://github.com/MrJc01/crompressor
cd crompressor && go build -o crompressor ./cmd/crompressor/

# Treinar codebook nos seus dados
./crompressor train -i /seus/dados/ -o meu.cromdb -s 8192

# Comprimir
./crompressor pack -i arquivo.txt -o arquivo.crom -c meu.cromdb

# Descomprimir (verificação SHA-256 automática)
./crompressor unpack -i arquivo.crom -o restaurado.txt -c meu.cromdb

# Benchmark
./crompressor benchmark -i arquivo.txt -c meu.cromdb
```

---

## O ecossistema (11 repos, tudo aberto)

| Repo | O quê | Stars |
|------|-------|-------|
| [crompressor](https://github.com/MrJc01/crompressor) | Motor principal (122 arquivos Go) | ⭐ |
| [crompressor-wasm](https://github.com/MrJc01/crompressor-wasm) | Motor no browser via WebAssembly | — |
| [crompressor-sync](https://github.com/MrJc01/crompressor-sync) | Sincronização P2P (libp2p) | — |
| [crompressor-gui](https://github.com/MrJc01/crompressor-gui) | Interface desktop nativa | — |
| [crompressor-video](https://github.com/MrJc01/crompressor-video) | Codec de vídeo experimental | — |
| [crompressor-security](https://github.com/MrJc01/crompressor-security) | Red-team e simuladores | — |
| [crompressor-matematica](https://github.com/MrJc01/crompressor-matematica) | Provas formais (7 papéis) | — |
| [crompressor-neuronio](https://github.com/MrJc01/crompressor-neuronio) | CromGPT (pesquisa neural) | — |
| [crompressor-ia](https://github.com/MrJc01/crompressor-ia) | Edge AI (llama.cpp) | — |
| [crompressor-sinapse](https://github.com/MrJc01/crompressor-sinapse) | Protocolo de transporte | — |
| [crompressor-projetos](https://github.com/MrJc01/crompressor-projetos) | Labs e demos | — |

---

## Pedindo ajuda à comunidade

Estou preparando papers formais para Zenodo e ArXiv. Mas antes, preciso de ajuda:

1. **Revisão de código**: o motor está em `pkg/cromlib/compiler.go`. Alguém com experiência em compressão pode identificar otimizações?
2. **Métricas de similaridade**: Hamming puro não é ideal para texto. Sugestões de métricas melhores?
3. **Benchmarks em cenários reais**: se você tem datasets específicos (backups, logs de produção, dados IoT), posso rodar o CROM e publicar os resultados.
4. **Validação das provas**: os 7 papéis em `crompressor-matematica` precisam de peer-review.

Se quiser contribuir: https://github.com/MrJc01/crompressor

---

## Próximos artigos desta série

- **Parte 2**: "30 dias construindo um motor de compressão — a jornada técnica"
- **Parte 3**: "Codebooks para IA — quando compressão encontra redes neurais"
- **Parte 4**: "Publicando pesquisa open-source — Zenodo + ArXiv para brasileiros"

---

*MrJ — [crom.run](https://crom.run) — Abril 2026*
