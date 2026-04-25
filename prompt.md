# CROM — Prompt de Orquestração do Ecossistema

> Este diretório (`crom/`) é o workspace central que contém **todos** os repositórios do ecossistema Crompressor.
> **NÃO** é um repositório Git — é uma pasta orquestradora para gerenciar tudo de um só lugar via Antigravity.

---

## Repositórios e Responsabilidades

| # | Repositório | Linguagem | Papel | Status |
|---|-------------|-----------|-------|--------|
| 1 | **crompressor** | Go | Motor core — CLI + biblioteca (Pack/Unpack/Verify/Train) | ✅ Produção |
| 2 | **crompressor-gui** | Go + React | Interface gráfica nativa (Lorca + REST API + Vite) | ✅ Estruturado |
| 3 | **crompressor-matematica** | Go + Markdown | Estudo matemático — provas, validações, Rate-Distortion | ✅ Pesquisa ativa |
| 4 | **crompressor-neuronio** | Python + Go | Pesquisa neural — CromGPT, PTQ, treinamento LLMs | ✅ Pesquisa ativa |
| 5 | **crompressor-ia** | Python | Inteligência artificial e ML | ✅ Pesquisa |
| 6 | **crompressor-security** | Go | Camada de segurança — AES-GCM, Ed25519, Kill-Switch | ✅ Estruturado |
| 7 | **crompressor-sinapse** | Go | Protocolo de transporte P2P | ✅ Estruturado |
| 8 | **crompressor-sync** | Go | Sincronização P2P — Bitswap, GossipSub, Kademlia | 🆕 Recém-criado |
| 9 | **crompressor-wasm** | Go (WASM) | Motor compilado para WebAssembly (browser) | 🆕 Recém-criado |
| 10 | **crompressor-video** | Go | Codec de vídeo CROM | ✅ Estruturado |
| 11 | **crompressor-projetos** | Misto | Projetos e aplicações usando o ecossistema | ✅ Portfólio |

---

## Arquitetura de Dependências

```
                    crompressor (CORE)
                    ┌─────┴─────┐
                    │           │
              pkg/cromlib    pkg/format
                    │
        ┌───────────┼───────────────┐──────────┐
        │           │               │          │
  crompressor-gui  crompressor-sync  crompressor-wasm  crompressor-security
   (importa core)  (importa core)   (compila core)    (importa core)
        │
        │
  crompressor-sinapse ←── crompressor-sync (protocolo de transporte)
```

```
  crompressor-matematica ←── Fundamento teórico (não depende de código)
  crompressor-neuronio   ←── Usa core via CLI/SDK para PTQ de LLMs
  crompressor-ia         ←── Pesquisa ML usando codebooks
  crompressor-video      ←── Codec especializado usando primitivas CROM
  crompressor-projetos   ←── Aplicações que usam o ecossistema
```

---

## Comandos Úteis (rodar de dentro de crom/)

```bash
# Status de todos os repos
for d in crompressor*/; do echo "=== $d ===" && cd "$d" && git status -sb && cd ..; done

# Pull de todos
for d in crompressor*/; do echo "=== $d ===" && cd "$d" && git pull 2>/dev/null && cd ..; done

# Build do core
cd crompressor && make build && cd ..

# Build da GUI
cd crompressor-gui && make build-ui && make build && cd ..

# Build WASM
cd crompressor-wasm && make build && cd ..
```

---

## Prioridades de Trabalho

### 🔴 Urgente (Fazendo agora)
- [x] Reestruturar `crompressor` main para padrão profissional
- [x] Criar `crompressor-gui` como repo separado
- [x] Criar `crompressor-sync` com código P2P extraído
- [x] Criar `crompressor-wasm` com entrypoint WASM

### 🟡 Próximo
- [x] Wiring: conectar `crompressor-wasm` ao `pkg/cromlib` real (PackBytes in-memory)
- [x] Wiring: conectar `crompressor-sync` independente do core (go.mod replace + pkg wrappers)
- [x] Wiring: conectar `crompressor-gui` independente do core (go.mod replace + pkg wrappers)
- [x] Atualizar descriptions de todos os 11 repos no GitHub
- [x] Garantir que cada repo compila independentemente (7/7 build matrix green)
- [x] Fix `crompressor-security` go.mod (module path + Go version + crommobile stub)
- [x] Criar pkg/ wrappers públicos no core (codebook, crypto, delta, entropy, vfs, trainer, network)

### 🟠 Próximo (Novo)
- [ ] Implementar `UnpackBytes` in-memory no WASM (roundtrip completo)
- [ ] Criar README padrão para repos minimalistas (wasm, sync)
- [ ] Criar `crompressor-sync/cmd/sync/main.go` CLI mínimo
- [ ] Limpar red team exploits no crompressor-security (múltiplos main no mesmo pkg)

### 🟢 Futuro
- [ ] CI/CD em todos os repos
- [ ] Release tags versionadas (v1.0.0)
- [ ] Documentação cruzada (links entre repos)
- [ ] Benchmark automatizado cross-repo
- [ ] Site unificado (crom-site)

---

## Regras de Organização

1. **Cada repo é independente** — deve compilar sozinho com `go build`
2. **Dependências via go.mod** — use `require github.com/MrJc01/crompressor` (não paths locais)
3. **Durante desenvolvimento local** — use `replace` directives no go.mod
4. **README em PT-BR** como principal, `README_en.md` opcional
5. **Licença MIT** em todos
6. **Branch `main`** = produção, `dev` = laboratório (apenas no core)
7. **Nunca commitar** binários, node_modules, datasets, modelos treinados

---

## Referência Rápida

- **GitHub:** https://github.com/MrJc01
- **Workspace local:** ~/Documentos/GitHub/crom/
- **Matemática:** crompressor-matematica/papeis/
- **Engine core:** crompressor/pkg/cromlib/compiler.go
- **Bifurcação Shannon:** docs/modes.md (Edge vs Archive)
- **Plano de migração:** crompressor/docs/MIGRATION_PLAN.md
