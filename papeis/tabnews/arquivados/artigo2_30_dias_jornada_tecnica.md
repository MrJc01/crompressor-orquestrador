# 30 dias construindo um motor de compressão do zero — a jornada real, as decisões e os erros

**TL;DR**: Do primeiro `git init` até 11 repositórios com 545 arquivos Go, 22 suites de teste, CI/CD e WASM. Sem formação acadêmica em compressão. Autodidata. Essa é a timeline real.

---

## Dia 0 — A ideia (26 de Março)

Tudo começou com uma pergunta: **e se o compressor já conhecesse seus dados antes de comprimir?**

Gzip e zstd constroem um dicionário novo para cada arquivo. E se o dicionário fosse pré-treinado, compartilhado, e otimizado para um domínio específico?

Isso não é uma ideia nova. A compressão baseada em dicionário existe desde os anos 90 (zstd --dict, brotli shared dictionaries). A diferença é que eu queria levar ao extremo: o dicionário como um **artefato soberano** — quem tem o dicionário, controla a compressão.

```bash
# Primeiro commit
git init crompressor
echo "start" > README.md && git add . && git commit -m "start"
```

## Semana 1 — Protótipo (26 Mar – 2 Abr)

### O que funcionou
- **Go como linguagem**: compilação rápida, concorrência nativa, deploy estático
- **Chunker fixo**: dividir arquivo em blocos de 128 bytes — simples e funcional
- **Busca linear**: para cada chunk, comparar com todos os codewords do codebook

### O que não funcionou
- **K-Means real**: muito lento para treinar com Go puro
- **Formato V1**: sem block table, carregava todo o arquivo na memória

### Decisão crucial: chunk size = 128 bytes
Por quê? Menor que 64 → overhead de metadados domina. Maior que 256 → menos chances de match. 128 é o sweet spot para Hamming distance em dados típicos.

## Semana 2 — LSH e V2 (2 – 9 Abr)

### LSH (Locality-Sensitive Hashing)
A busca linear era O(N×K) — com K=65536 codewords, inviável. LSH reduz para O(1) amortizado via hash tables espaciais.

```go
// internal/search/lsh.go
func (ls *LSHSearcher) FindBestMatch(chunk []byte) (MatchResult, error) {
    // 1. Gera hash do chunk
    // 2. Consulta bucket correspondente
    // 3. Compara só com candidatos do bucket
    // 4. Se bucket vazio → fallback para busca linear
}
```

### Formato V2 — Block Table
V1 carregava tudo na memória. V2 processa bloco por bloco:

```
V1: [header] [todos os chunks] [todos os deltas] → precisa ler tudo
V2: [header] [block_table] [block1] [block2] ... → streaming possível
```

## Semana 3 — P2P, Security, Video (9 – 17 Abr)

Nessa semana, o escopo explodiu. Em vez de um compressor, comecei a construir um ecossistema:

| Repo | Motivação |
|------|-----------|
| crompressor-security | "E se alguém interceptar o .crom?" |
| crompressor-sinapse | "E se dois nós trocarem chunks?" |
| crompressor-video | "E se comprimir frames como chunks 2D?" |

### O erro: expandir antes de estabilizar
O motor principal ainda tinha Hit Rate 0% (descobri isso só agora, honestamente). Mas eu já estava construindo proxy TCP, red-team simulators e codec de vídeo.

**Lição aprendida**: estabilize o core antes de escalar o ecossistema.

## Semana 4 — CromGPT e Pesquisa Neural (17 – 22 Abr)

### A tese
Se um codebook pode representar padrões de dados, pode representar pesos de uma rede neural.

**CromLinear**: substituir `nn.Linear` (multiplicação de matrizes) por lookup em codebook + interpolação.

```python
# nn.Linear: y = Wx + b (milhões de multiplicações)
# CromLinear: y = codebook[quantize(x)] + b (lookup + soma)
```

### Resultado
- Modelo de 125M parâmetros treinado na Wikipedia PT
- Gera texto coerente em português
- Loss comparável ao baseline nn.Linear
- Mas: **PTQ de modelos SOTA ainda não funciona** (Hybrid PTQ em andamento)

**O nome técnico que eu esquecia**: isso se chama **Product Quantization (PQ)** ou **Vector Quantization for Neural Networks**. Google (AQLM) e Meta (QuIP#) publicaram variações em 2024-2025.

## Semana 5 — Hardening (24 Abr)

### O que eu fiz nas últimas 12 horas
- Removidos 150MB de binários commitados
- Reestruturado para layout `cmd/pkg/internal`
- Extraído GUI, Sync e WASM para repos separados
- Criados 7 `pkg/` wrappers para desacoplamento
- CI/CD GitHub Actions em 5 repos
- PackBytes + UnpackBytes para WASM in-memory
- Demo HTML drag-and-drop para WASM

### Estado final
```
11/11 repos limpos
22/22 test suites passando
 5/5  CI workflows
 0    erros de build
```

---

## Os erros que cometi

| Erro | Impacto | Lição |
|------|---------|-------|
| Linguagem exagerada nos artigos | Perda de credibilidade | Benchmarks > buzzwords |
| Expandir antes de estabilizar | Core com Hit Rate 0% | Fix the engine first |
| Não testar com dados reais | Benchmarks inflados | Use dados do próprio sistema |
| Commitar binários de 150MB | Repo inutilizável | .gitignore desde o dia 0 |
| 11 repos para 1 pessoa | Overhead de manutenção | Monorepo pode ser melhor |

---

## O que vem depois

1. **Corrigir o Hit Rate** — investigar métricas de similaridade melhores que Hamming
2. **Benchmarks em cenários reais** — dados binários, backups incrementais, dados IoT
3. **Paper formal** — preparando para Zenodo e ArXiv
4. **Vídeo demonstrativo** — demo WASM, CLI, VFS em ação

Se você leu até aqui: obrigado. E se quiser ajudar: https://github.com/MrJc01/crompressor

---

*MrJ — [crom.run](https://crom.run) — Abril 2026*
