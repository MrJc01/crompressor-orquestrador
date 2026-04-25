# Parte 3 — O Ecossistema dos 11 Repositórios

## 3.1 Mapa Geral

O Crompressor não é um único programa — é um **ecossistema de 11 repositórios** com responsabilidades distintas. Cada repo é um módulo Go independente que compila separadamente.

```
                    ┌─────────────────────┐
                    │   crompressor       │
                    │   (Motor Principal) │
                    │   122 arquivos Go   │
                    └──────┬──────────────┘
                           │ pkg/ wrappers
              ┌────────────┼────────────┬──────────────┐
              ▼            ▼            ▼              ▼
        ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐
        │   gui    │ │   sync   │ │   wasm   │ │  video   │
        │ Desktop  │ │   P2P    │ │ Browser  │ │  Codec   │
        └──────────┘ └──────────┘ └──────────┘ └──────────┘

        ┌──────────┐ ┌──────────┐ ┌──────────┐
        │ security │ │ sinapse  │ │ matemat. │
        │ Red-Team │ │Transport │ │  Provas  │
        └──────────┘ └──────────┘ └──────────┘

        ┌──────────┐ ┌──────────┐ ┌──────────┐
        │neuronio  │ │    ia    │ │ projetos │
        │ CromGPT  │ │ Edge AI  │ │   Labs   │
        └──────────┘ └──────────┘ └──────────┘
```

## 3.2 Cada Repositório em Detalhe

### 🧬 crompressor (Core Engine)
**O quê**: Motor de compressão principal. Contém toda a lógica de pack/unpack, codebook, chunker, search, delta, crypto, VFS, rede.
**Linguagem**: Go (15.011 linhas)
**Dependências**: stdlib Go + zstd + libp2p + bazil/fuse
**API pública**: 11 packages em `pkg/`

### 🖥️ crompressor-gui (Interface Desktop)
**O quê**: Interface gráfica nativa usando Lorca (Chrome headless) + WebSocket + Vite.
**Depende de**: `pkg/vfs`, `pkg/trainer`, `pkg/network`, `pkg/crypto` do core
**Como roda**: `go run ./cmd/gui/` abre janela Chrome com frontend web
**Funcionalidades**: Compressão drag-and-drop, visualização de codebook, controle de nó P2P

### 🔄 crompressor-sync (Sincronização P2P)
**O quê**: Camada de sincronização descentralizada usando libp2p.
**Depende de**: `pkg/codebook`, `pkg/crypto`, `pkg/delta` do core
**Protocolos**: Bitswap (transferência de chunks), GossipSub (anúncios), Kademlia (descoberta)
**Conceito-chave**: Soberania via Codebook — apenas nós com o mesmo codebook podem trocar dados
**CLI**: `cromsync start --codebook trained.cromdb --port 4001`

### 🌐 crompressor-wasm (WebAssembly)
**O quê**: O motor compilado para WASM, roda no browser.
**Depende de**: `pkg/cromlib` (PackBytes/UnpackBytes), `pkg/entropy` do core
**Binário**: 11MB .wasm
**API JS**: `cromPack()`, `cromUnpack()`, `cromAnalyze()`, `cromInfo()`
**Demo**: `demo/index.html` — interface drag-and-drop com métricas visuais

### 🎬 crompressor-video (Codec de Vídeo)
**O quê**: Codec experimental que usa VQ para comprimir frames de vídeo.
**Conceito**: Extrai frames via ffmpeg → chunka cada frame 2D → busca padrões visuais no codebook
**Depende de**: `pkg/format` do core

### 🔐 crompressor-security (Camada de Segurança)
**O quê**: Red-team, simuladores de ataque, proxy TCP, ferramentas de pentest.
**Componentes**: Client SDK (crommobile), proxy alpha (TCP tunnel), alien sniffer, race condition PoCs
**Conceito**: Validar que o protocolo CROM resiste a ataques reais

### ⚡ crompressor-sinapse (Protocolo de Transporte)
**O quê**: Implementação de baixo nível do protocolo de transporte P2P.
**Conceito**: Camada de rede "raw" abaixo do sync — gerenciamento de conexões, heartbeat, framing

### 📐 crompressor-matematica (Fundação Matemática)
**O quê**: Provas formais da teoria por trás do CROM.
**Conteúdo**: 7 papéis (papel0-papel6) cobrindo Rate-Distortion, Bifurcação de Shannon, testes de impossibilidade
**Validação**: Suite de testes Go que verifica as provas computacionalmente (35/35 passando)

### 🧠 crompressor-neuronio (Pesquisa Neural)
**O quê**: CromGPT — transformer treinado com CromLinear layers em vez de nn.Linear.
**Pesquisas**: pesquisa0 (5D Active Inference), pesquisa1 (CromGPT training), pesquisa2 (validação), pesquisa3 (PTQ)
**Linguagem**: Python (PyTorch)

### 🤖 crompressor-ia (Edge AI)
**O quê**: Integração com llama.cpp para rodar LLMs comprimidos em dispositivos edge.
**Abordagem**: PTY wrapper sobre llama-cli, KV cache persistente

### 🧪 crompressor-projetos (Labs)
**O quê**: Hub central de demos interativos e simuladores educacionais.
**Conteúdo**: 15 projetos web, simuladores de VQ, visualizadores de codebook

## 3.3 Como os Repos se Conectam (pkg/ wrappers)

O Go proíbe que módulos externos importem packages `internal/`. Solução: **thin public wrappers**:

```go
// crompressor/pkg/codebook/codebook.go
package codebook

import "github.com/MrJc01/crompressor/internal/codebook"

type Reader = codebook.Reader    // type alias = zero overhead
type Header = codebook.Header

func Open(path string) (*Reader, error) {
    return codebook.Open(path)
}
```

Para desenvolvimento local, cada satélite tem um `replace` no go.mod:
```
replace github.com/MrJc01/crompressor => ../crompressor
```

## 3.4 CI/CD

5 repos têm GitHub Actions configuradas:

| Repo | Trigger | Steps |
|------|---------|-------|
| crompressor | push/PR main | build → test → vet |
| crompressor-gui | push/PR main | checkout core → build → test |
| crompressor-sync | push/PR main | checkout core → build → test |
| crompressor-wasm | push/PR main | checkout core → WASM build → vet |
| crompressor-video | push/PR main | build → test → vet |

---

**Próximo**: [Parte 4 — Pesquisa Neural](./parte4_neural.md)
