# Parte 5 — Cronologia Completa do Projeto Crompressor

## Timeline

### Semana 1 — Fundação (26 Mar – 2 Abr 2026)

**26 de Março**: Primeiro commit do `crompressor`. O motor começa como um protótipo simples: chunker fixo + busca linear em codebook.

**Decisões iniciais**:
- Linguagem Go (performance + compilação estática)
- Formato binário customizado (.crom)
- Codebook como arquivo separado (.cromdb)

**2 de Abril**: `crompressor-ia` criado — primeiros experimentos com LLMs locais em CPU.

### Semana 2 — Expansão (2 – 9 Abr)

**6 de Abril**: `crompressor-projetos` criado — hub central com 15 demos interativos.

**7 de Abril**: `crompressor-sinapse` criado — protocolo de transporte P2P. Primeira tentativa de comunicação entre nós CROM.

**9 de Abril**: `crompressor-security` criado — proxy TCP, simuladores de ataque, first red-team session. Descoberta de vulnerabilidades no framing do protocolo.

### Semana 3 — Motor V2 (9 – 17 Abr)

**Formato V2**: Implementação de block table para streaming. O V1 carregava todo o delta pool na memória; V2 processa bloco por bloco.

**LSH Search**: Substituição do LinearSearcher por LSHSearcher. Performance de busca vai de O(N×K) para O(1) amortizado.

**11 de Abril**: `crompressor-neuronio` criado — início da pesquisa CromGPT.

**17 de Abril**: `crompressor-video` criado — primeiros experimentos com codec de vídeo semântico.

### Semana 4 — Pesquisa Neural (17 – 22 Abr)

**Pesquisa 0 (5D Active Inference)**: 6 laboratórios explorando percepção dimensional.

**Pesquisa 1 (CromGPT)**: Treinamento de transformer 125M com CromLinear na Wikipedia PT. Convergência alcançada, amostras de texto coerentes geradas.

**22 de Abril**: Modelo CromGPT 125M validado qualitativamente. Início da Pesquisa 3 (PTQ de modelos SOTA).

### Semana 5 — Hybrid PTQ + Matemática (22 – 24 Abr)

**23 de Abril**: Tentativa de PTQ full-model falha (modelo perde coerência). Decisão de migrar para **Hybrid PTQ** — preservar Attention/LM-Head em FP16.

**Kernels C++ AVX**: Implementação de kernel customizado para forward pass do CromLinear em CPU.

**24 de Abril (manhã)**: `crompressor-matematica` recebe 7 papéis formais. Papel 5 com 15 Testes de Impossibilidade. Papel 6 responde ao peer-review sobre Bifurcação de Shannon.

### Semana 5 — Hardening do Ecossistema (24 Abr, tarde/noite)

**16:00 — Reestruturação Profissional do Core**:
- Remoção de 150MB de binários commitados
- Reorganização para layout `cmd/pkg/internal`
- Extração do GUI para `crompressor-gui`
- Extração do Sync para `crompressor-sync`
- Criação do `crompressor-wasm`

**20:00 — Sessão de Wiring** (esta sessão):
- Criação de 7 `pkg/` wrappers públicos no core
- Migração de imports em 3 satélites (gui, sync, wasm)
- `PackBytes` e `UnpackBytes` in-memory para WASM
- Fix de testes quebrados (vfs, video)
- Fix do module path do security
- Reorganização dos simuladores do security
- CI/CD em 5 repos
- READMEs profissionais
- Demo HTML para WASM
- CLI para sync
- GitHub descriptions para todos os 11 repos
- 14 commits pushados

## Métricas do Projeto

| Métrica | Valor |
|---------|-------|
| Repositórios | 11 |
| Commits totais | ~261 |
| Arquivos Go | ~545 |
| Linhas de código (Go) | ~72.000 |
| Papéis de pesquisa | 7 (matemática) + 3 (neural) |
| Duração do projeto | 30 dias (26 Mar – 24 Abr) |
| Contribuidores | 1 |

## Decisões Arquiteturais Históricas

| Data | Decisão | Alternativa Rejeitada | Motivo |
|------|---------|----------------------|--------|
| 26 Mar | Go em vez de Rust | Rust | Velocidade de iteração, stdlib robusta |
| 26 Mar | Codebook externo (.cromdb) | Codebook embutido no .crom | Reutilização entre arquivos |
| ~5 Abr | LSH em vez de busca linear | HNSW, KD-tree | Simplicidade + O(1) amortizado |
| ~10 Abr | libp2p para P2P | gRPC, custom TCP | Ecossistema maduro (Kademlia, GossipSub) |
| ~15 Abr | FUSE VFS | API de leitura parcial | UX transparente (qualquer programa lê .crom) |
| 22 Abr | CromLinear (VQ layers) | Pruning, destilação | Alinhamento com a tese do codebook |
| 23 Abr | Hybrid PTQ | Full PTQ | Preservar Attention = preservar coerência |
| 24 Abr | pkg/ wrappers | Mover tudo para pkg (sem internal) | Manter encapsulamento + compatibilidade |
| 24 Abr | replace directives | Publicar tags semver | Desenvolvimento local ágil |

## O Que Define o Crompressor

1. **Compressão é conhecimento** — o codebook é um "modelo do mundo" treinado para um domínio
2. **Soberania** — nós que compartilham o mesmo codebook formam uma rede de confiança
3. **Dualidade compressão/IA** — os mesmos codebooks servem para comprimir dados e para inferência neural
4. **Modularidade radical** — 11 repos independentes, cada um compilável separadamente
5. **Do browser ao bare-metal** — WASM (browser), GUI (desktop), CLI (servidor), VFS (kernel)

---

**Fim do Guia de Estudo. Bom estudo! 🧬**
