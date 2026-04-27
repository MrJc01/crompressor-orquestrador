# 🛤️ CROM Roadmap: O Que Podemos Explorar

## O Que Já Temos (Capital Acumulado em 9 Pesquisas)

| Pesquisa | O Que Construímos | Estado |
|---|---|---|
| P1-P3 | Motor LSH base, Sigmóide calibrada, Max Pooling | ✅ Sólido |
| P4 | Feedback/Recompensa, Memória Evolutiva, Fusão Multimodal | ✅ Conceito |
| P5 | PCA real (SVD Power Iteration), Multi-Head Hashing | ✅ Produção |
| P6 | RAG completo (TF-IDF+PCA+Hamming), 100% precision, 124 intents | ✅ Produção |
| P7 | Escala massiva (SQuAD ingestado) | ✅ Dados |
| P8 | Micro-GPT conceito (causal attention) | 📝 Papel |
| P9 | **MLM Go-Native, Backprop, Cosine LR, Gradient Clip, Difusão** | ✅ Funcional |

---

## Caminho 1: 📊 Mais Dados (Impacto Imediato, Esforço Baixo)

### O Problema
O modelo treina em **50 frases** (~340 palavras). Qualquer arquitectura vai memorizar e misturar com tão pouco.

### O Que Fazer
1. **Ingerir o SQuAD-PT** — A Pesquisa 7 já tem dados do SQuAD. Extrair as respostas em PT e usar como corpus
2. **Wikipedia PT resumida** — Scrape dos primeiros parágrafos dos 1000 artigos mais populares
3. **Gerar paráfrases** — Para cada frase existente, criar 3-5 variações (sujeito diferente, ordem diferente)
4. **Meta: 2000+ frases, 5000+ palavras únicas**

### Impacto Esperado
Com 10× mais dados, o modelo passa de "fragmentos reconhecíveis" para "frases semi-fluentes". A loss deve cair de ~0.7 para <0.3.

---

## Caminho 2: 💾 Salvar/Carregar Pesos (Prático, Esforço Baixo)

### O Problema
Cada execução treina do zero (~6 minutos). Impossível iterar rápido.

### O Que Fazer
1. Serializar pesos treinados no formato `.crom` existente (já temos o leitor)
2. Na inicialização, verificar se existe `.crom` treinado → carregar → saltar treino
3. Flag `--retrain` para forçar novo treino

### Impacto Esperado
Arranque instantâneo. Permite experimentar com temperature/top-k sem esperar 6 minutos.

---

## Caminho 3: 🔀 RAG + Neural Híbrido (Melhor UX, Esforço Médio)

### O Problema
O RAG responde perfeito mas só para queries conhecidas. O Neural gera mas é incoerente. Separados, ambos falham.

### O Que Fazer
1. **Pipeline:** Input → RAG primeiro → Se confiança > 70%, retorna resposta exacta
2. **Fallback:** Se RAG falha → Neural gera completação
3. **Bonus:** Usar o documento RAG recuperado como **contexto** para o Neural (injectar na janela)

### Impacto Esperado
O utilizador vê respostas perfeitas quando existem, e geração "criativa" quando não existem. Nenhum caso fica sem resposta.

---

## Caminho 4: 🧠 Self-Attention Escalada (Alto Impacto, Esforço Alto)

### O Que Aprendemos
A attention não superou a janela fixa com 50 frases. MAS com 2000+ frases, a attention deve brilhar.

### O Que Fazer
1. Primeiro implementar Caminho 1 (mais dados)
2. Depois retomar a Self-Attention com:
   - **Multi-head** (4 heads × 16 dims = 64 total)
   - **2 camadas empilhadas** (2-layer Transformer)
   - **Layer Normalization**
   - **Residual connections**
3. Treinar com ~500K parâmetros em 2000+ frases

### Impacto Esperado
Frases gramaticalmente correctas e semanticamente coerentes. O salto qualitativo real.

---

## Caminho 5: 📐 BPE Tokenizer em Go (Fundamental, Esforço Médio)

### O Problema
O vocabulário actual é word-level (340 palavras). Palavras desconhecidas viram `[MASK]`. O modelo não pode generalizar para palavras novas.

### O Que Fazer
1. Implementar Byte-Pair Encoding (BPE) em Go
2. Vocabulário de ~2000 subwords (cobre qualquer palavra PT)
3. "computador" → ["comp", "ut", "ador"] — o modelo vê partes que partilha com "computação"

### Impacto Esperado
Elimina palavras desconhecidas. O modelo pode inferir sobre palavras nunca vistas porque partilham subwords com palavras conhecidas.

---

## Caminho 6: 🔬 Benchmarks Científicos (Validação, Esforço Baixo)

### O Que Fazer
1. Criar um **gauntlet de avaliação** com 50 prompts classificados:
   - 10 factuais ("a gravidade é") → espera-se precisão
   - 10 criativos ("era uma vez") → espera-se fluência
   - 10 identidade ("quem é você") → espera-se consistência
   - 10 OOD ("receita de bolo") → espera-se rejeição honesta
   - 10 ambíguos ("o que é isso") → espera-se generalização
2. Pontuar cada resposta (0-3): 0=lixo, 1=fragmento, 2=semi-coerente, 3=perfeito
3. Comparar versões automaticamente

### Impacto Esperado
Métrica objectiva para saber se uma mudança melhorou ou piorou. Sem isto, estamos a avaliar "a olho".

---

## Ordem Recomendada

```
[1] Salvar/Carregar Pesos  →  iteração rápida
[2] Benchmarks             →  saber o que medir
[3] Mais Dados             →  o maior ganho
[4] RAG+Neural Híbrido     →  melhor experiência
[5] BPE Tokenizer          →  generalização
[6] Self-Attention Escalada →  o salto final
```

> **Princípio:** Dados > Arquitectura. Até o modelo mais simples (janela fixa) produz resultados bons com dados suficientes. A attention só vale depois de ter os dados.
