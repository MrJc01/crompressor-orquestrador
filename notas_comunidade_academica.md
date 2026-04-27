# 📝 Notas da Comunidade — Feedback para Paper e Direção Acadêmica

> Arquivo vivo para anotar ideias, recomendações e insights recebidos da comunidade (TabNews, GitHub, etc.)
> que impactam diretamente a estratégia de publicação e evolução do Crompressor.

---

## 2026-04-26 — @clacerda (TabNews, Papel0)

**Contexto:** Comentário no artigo "A Ilusão da Compressão" sobre como entrar na conversa acadêmica certa.

### Recomendações Recebidas

#### 1. Mapear as "Casas Naturais" do Crompressor
- [ ] Identificar 5-10 venues/conferências/workshops que discutem a intersecção entre **compressão, deduplicação e sistemas P2P**
- [ ] Ler introdução, metodologia e avaliação de papers próximos para entender:
  - Como a comunidade fala (idioma acadêmico)
  - Quais **benchmarks** usam
  - Quais **datasets** são padrão
  - Quais **métricas** são esperadas
  - Qual é o **formato dos papers**
  - Quais trabalhos são **citação obrigatória**

#### 2. Venues Potenciais para Investigar
- [ ] **IEEE** — Conferências de storage e sistemas distribuídos
- [ ] **ACM** — SIGCOMM, SIGMOD, SOSP (sistemas operacionais)
- [ ] **USENIX** — FAST (File and Storage Technologies), ATC, OSDI
- [ ] **Springer** — Journals de compressão e data management
- [ ] **Elsevier** — Information Sciences, Journal of Systems and Software
- [ ] **ArXiv** — Categorias: cs.DC (Distributed Computing), cs.IR (Information Retrieval), cs.DS (Data Structures)

#### 3. Transformar o Projeto em Pesquisa Formal
- **DE:** "Eu tenho um motor que funciona"
- **PARA:** "Eu tenho uma questão de pesquisa, uma hipótese e um experimento reproduzível"

**Questões de pesquisa candidatas:**
- [ ] "Codebook-based deduplication pode substituir compressão clássica (LZ77) em cenários de sincronização P2P com redundância massiva?" (Core engine)
- [ ] "Vector Quantization de pesos neurais via Codebook compartilhado pode atingir accuracy comparável com redução de parâmetros > 40x?" (Neural)
- [ ] "LSH-guided RAG pode servir como governor sintático para LLMs generativos rodando em CPU-only?" (Semântico)

#### 4. Acesso a Bases Acadêmicas (Sem Vínculo Universitário)
- [ ] Acessar máquinas de **bibliotecas de universidades públicas** (acesso local a IEEE, ACM, Springer, Elsevier)
- [ ] Usar preprints do **ArXiv** e **Zenodo** (abertos)
- [ ] **E-mail universitário NÃO é pré-requisito** para publicar — o que pesa é:
  - Trabalho bem escrito
  - Metodologia reproduzível
  - Claims calibradas (sem hype)

#### 5. Estratégia de Abordagem a Pesquisadores
- **NÃO** chegar pedindo nada logo de primeira
- **SIM** mandar e-mail simples:
  > "Estou trabalhando nisso, acho que conversa com tal trabalho seu. Se tiver interesse de dar uma olhada."
- Depois de mapear venues e papers, os pesquisadores certos vão aparecer naturalmente
- O endorsement vem do trabalho, não do pedido

---

## Ideias Futuras da Comunidade

> Espaço para anotar novas ideias conforme chegarem.

### @ExtraMobs — IPFS + Crompressor (2026-04-26)
- Gateway IPFS com Codebook CROM pinado como CID
- `.cromdb` distribuído descentralizado via DHT do próprio IPFS
- Bitswap negociando ponteiros de 24 bytes em vez de blocos crus
- **Status:** Ideia documentada, bridge não implementada ainda

### @bugzoid — LLMs sem GPU via Crompressor (2026-04-26)
- Codebook como espaço de aprendizado substituto de tensores float
- Mudar paradigma de acesso à memória (não apenas armazenamento)
- Referência direta aos repos: crompressor-neuronio, crompressor-ia, crompressor-sinapse
- **Status:** Pesquisa ativa (Pesquisa 9 / CROM-Chat)

---

## Checklist Geral de Ações

- [ ] Mapear 5-10 venues/papers próximos ao Crompressor
- [ ] Ler 3 papers completos de cada venue (intro + metodologia + avaliação)
- [ ] Adaptar linguagem dos papéis (pasta `papeis/`) para idioma acadêmico
- [ ] Formular 3 questões de pesquisa formais
- [ ] Montar experimento reproduzível com dataset público e métricas padrão
- [ ] Identificar 3-5 pesquisadores cujos trabalhos conversam com o CROM
- [ ] Redigir e-mail modelo para abordagem inicial
- [ ] Visitar biblioteca de universidade pública para acessar bases pagas
