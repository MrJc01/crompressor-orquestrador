# Publicando pesquisa open-source: Zenodo + ArXiv para devs brasileiros — guia prático

**TL;DR**: Estou organizando o Crompressor para publicação acadêmica formal. Aqui está o que aprendi sobre Zenodo, ArXiv e como qualquer dev brasileiro pode publicar pesquisa sem precisar de universidade.

---

## Por que publicar?

Código no GitHub é acessível, mas não é **citável**. Se alguém citar seu repositório, a URL pode mudar, o repo pode ser deletado, o commit pode ser reescrito.

**Zenodo** e **ArXiv** resolvem isso:
- **Zenodo**: dá um DOI (Digital Object Identifier) permanente ao seu código/dados
- **ArXiv**: dá um identificador permanente ao seu paper/preprint

Com DOI, seu trabalho vira uma **referência bibliográfica formal**. Universidades e pesquisadores podem citar como qualquer outro paper.

---

## 1. Zenodo — Publicar código e dados

### O que é
Repositório de dados científicos mantido pelo **CERN** (sim, os mesmos do acelerador de partículas). Gratuito, permanente, DOI automático.

### Como usar com GitHub

1. Criar conta em [zenodo.org](https://zenodo.org) (login via GitHub)
2. Ir em Settings → GitHub → Conectar repositório
3. Criar uma Release no GitHub (ex: `v0.1.0`)
4. Zenodo detecta automaticamente e gera DOI

```bash
# No GitHub:
git tag v0.1.0
git push --tags
# → Zenodo gera: DOI 10.5281/zenodo.XXXXXXX
```

### O que incluir na release
- README completo com instruções de build
- Codebook de exemplo (.cromdb) para reprodução
- Script de benchmark (`run_benchmark.sh`)
- Dataset de teste (ou link para download)
- Licença (MIT, Apache, etc.)

### Resultado
Seu código agora é citável:

```bibtex
@software{crompressor2026,
  author    = {MrJ},
  title     = {Crompressor: Codebook-Based Compression Engine},
  year      = {2026},
  publisher = {Zenodo},
  doi       = {10.5281/zenodo.XXXXXXX},
  url       = {https://doi.org/10.5281/zenodo.XXXXXXX}
}
```

---

## 2. ArXiv — Publicar o paper

### O que é
Servidor de preprints mais usado em computação (CS), física e matemática. Mantido pela **Cornell University**. Gratuito.

### Requisitos
1. **Paper em LaTeX** (não aceita Word/Markdown)
2. **Endorsement**: primeira vez precisa de aval de alguém que já publicou na mesma categoria
3. **Categoria correta**: para compressão → `cs.IT` (Information Theory); para IA → `cs.LG` (Machine Learning)

### Estrutura do paper

```latex
\title{CROM: A Codebook-Based Compression Engine
       with Sovereign Synchronization}
\author{MrJ}

\begin{abstract}
We present CROM, a lossless compression system that uses
pre-trained codebooks for pattern matching instead of
traditional LZ77-based methods. While CROM does not outperform
state-of-the-art compressors in raw ratio, it enables
random access, P2P synchronization, and convergent encryption
— capabilities absent in conventional tools. We provide
benchmarks against gzip and zstd on real-world datasets.
\end{abstract}

\section{Introduction}
\section{Related Work}
\section{Method}
\section{Experiments}
\section{Results}
\section{Discussion}
\section{Conclusion}
```

### Como conseguir endorsement
- Pedir em fóruns acadêmicos (r/MachineLearning, r/compsci)
- Contatar professores de universidades brasileiras
- Participar de conferências (SBC, SBRC, WebMedia)
- Contribuir para projetos open-source de pesquisadores que publicam

---

## 3. O que estou preparando

### Paper 1: Motor de compressão
- Foco: arquitetura do codebook, pipeline Pack/Unpack, comparação justa com zstd
- Estado: rascunho em `crompressor-matematica/papeis/paper.tex`
- Falta: benchmarks reproduzíveis com datasets padronizados

### Paper 2: CromGPT
- Foco: Product Quantization para transformers, CromLinear layer
- Estado: resultados do 125M documentados
- Falta: comparação formal com AQLM/GPTQ

### Vídeo demonstrativo
- Demo WASM (compressão no browser)
- Demo CLI (train → pack → unpack → verify)
- Demo VFS (mount e acesso aleatório)
- Planejando gravação em breve

---

## Para a comunidade

Meu pedido é simples:

1. **Se você tem experiência acadêmica**: me ajude a revisar os papers
2. **Se você é dev**: rode os benchmarks no seu dataset e me mande os resultados
3. **Se você é curioso**: clone o repo e me diga o que quebrou
4. **Se você é crítico**: me aponte onde estou errado (como fizeram nos artigos anteriores — e eu agradeço)

Repositório: https://github.com/MrJc01/crompressor
Contato: crom.run

---

*MrJ — [crom.run](https://crom.run) — Abril 2026*
