# Pesquisa 9 — Papel 4: A Crise Lexical e o Paradigma do BPE (Byte-Pair Encoding)

## 1. A Autópsia da "Alucinação Estocástica"

Durante a transição do CROM-LLM para a **Arquitetura Híbrida V9** (onde o motor executa K-NN RAG para 100% de precisão e usa Redes Neurais de Difusão como fallback para OOD - Out Of Distribution), esbarramos num obstáculo comportamental crítico. 

Quando o utilizador submeteu o prompt simples `"oi"`, a camada RAG (que não tinha este cumprimento mapeado) delegou a tarefa para a Difusão Neural. O modelo neural respondeu: `"oi a a galáxia é a medida da química a a galáxia de"`.

**A Causa Raiz:** O nosso `corpusTreino` na V9 era composto por meras 300 frases focadas estritamente em ciência, tecnologia e astrofísica. O modelo nunca foi exposto à estrutura sintática de um diálogo social real ("oi", "tudo bem", "como estás"). Confrontado com a palavra inédita `"oi"` e forçado matematicamente a resolver a matriz de *Masked Language Modeling* (MLM), o modelo preencheu os `[MASK]` com as sinapses de maior peso residual que conhecia: "galáxia", "química".

## 2. A Barreira da Escala (O Paradoxo das 50.000 Mensagens)

A solução óbvia seria injetar o dataset SQuAD-PT (com mais de 50.000 mensagens reais). No entanto, o nosso Tokenizer atual era **Word-Level** (dividia o texto por espaços em branco). 

**O Problema Matemático:**
- 300 frases geraram `1.022` palavras únicas. A matriz final (`Hidden -> Vocab`) pesava apenas **1.5 MB** e treinava em 6 minutos.
- 50.000 frases gerariam mais de **60.000 a 80.000 palavras únicas**.
- Com 80.000 palavras únicas, as matrizes de projeção (Embeddings e Camadas Densas) explodiriam em tamanho (Centenas de Megabytes de memória RAM alocada em CPU). O tempo de treino passaria de minutos para **dias**, e o processo fatalmente resultaria em `OOM (Out Of Memory)` num ambiente de borda.

Para romper esta barreira e injetar conhecimento enciclopédico massivo sem destruir a máquina, identificámos duas arquiteturas de contenção.

---

## 3. As Estratégias de Contenção Lexical

### O Caminho 1: Limite Guloso Frequencista (Anotado como Fallback/Plan B)
A primeira abordagem considerada foi um limite estatístico brutal.
* **Mecanismo:** Ler as 50.000 mensagens, contar todas as palavras, mas **só guardar as 2.000 a 3.000 palavras mais usadas**.
* **Tratamento de Anomalias:** Qualquer palavra fora deste Top-K seria forçada a virar um token `<UNK>` (Unknown).
* **Vantagens:** Extremamente simples de codificar em Go. Treino ultra-rápido (manteria o tamanho da matriz controlado).
* **Desvantagens (Motivo da Rejeição):** O modelo sofre de "cegueira funcional". Se lhe pedirem para falar sobre "crompressor" e a palavra não estiver no Top 3.000, ele será incapaz de a processar ou gerar, limitando permanentemente a inteligência do LLM.

### O Caminho 2: BPE (Byte-Pair Encoding) — A Decisão Tomada
Esta é a arquitetura utilizada pela OpenAI (GPT-3/4) e LLaMA. Em vez de decorar palavras inteiras, a rede aprende sílabas e fragmentos comuns.
* **Mecanismo:** O algoritmo começa com um vocabulário base de caracteres individuais (a, b, c). Ele vasculha as 50.000 frases à procura dos pares de caracteres mais adjacentes (ex: 'e' + 's' -> 'es'). Ele funde esses pares iterativamente até atingir o limite de vocabulário desejado (ex: 2.000 *subwords*).
* **Vantagens (O Santo Graal):** O modelo **nunca perde informação**. A palavra "descomplicado" pode não existir no vocabulário, mas o modelo consegue construí-la e compreendê-la juntando `"des"` + `"comp"` + `"li"` + `"cado"`.
* **Desafio Técnico:** Exige a escrita de um algoritmo recursivo de compressão em Go puro, e um "Tokenization Trainer" independente do treino da rede neural.

---

## 4. Conclusão da Pesquisa: A Direção do CROM V9

Por decisão executiva, o projeto prosseguirá com o **Caminho 2 (BPE em Go)**. 

Ao adotarmos a tokenização sub-lexical, o CROM-LLM V9 atingirá a imunidade a palavras OOV (Out Of Vocabulary). Ele será capaz de devorar as 50.000 mensagens do SQuAD, comprimindo a riqueza do léxico português num espaço matemático minúsculo (~2.000 a 3.000 matrizes de vocabulário), pavimentando a estrada para um motor conversacional que fala de forma fluente e criativa, mantendo a inferência viável para execução estrita em CPU local.
