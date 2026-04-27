"""
Pesquisa 7 — RAG Massivo: Ingestão Enciclopédica
Objetivo: Escalar o cérebro vetorial para 100K+ entradas reais
Datasets: SQuAD v1.1 (completo) + Wikipedia PT (Simple English como proxy)
"""
import json
import re
import math
import unicodedata
import urllib.request
import sys
import random
from collections import defaultdict, Counter

# =============================================================================
# 1. Download e Parsing dos Datasets
# =============================================================================

def baixar_squad_completo():
    """Baixa o SQuAD v1.1 COMPLETO (~87K pares Q&A)."""
    url = 'https://rajpurkar.github.io/SQuAD-explorer/dataset/train-v1.1.json'
    print(f"[*] Baixando SQuAD v1.1 completo...")
    req = urllib.request.Request(url, headers={'User-Agent': 'Mozilla/5.0'})
    with urllib.request.urlopen(req, timeout=120) as response:
        dados = json.loads(response.read().decode())
    
    corpus = []
    for artigo in dados['data']:
        for paragrafo in artigo['paragraphs']:
            contexto = paragrafo['context']
            for qa in paragrafo['qas']:
                if len(qa['answers']) > 0:
                    intent = qa['question']
                    answer = qa['answers'][0]['text']
                    # Enriquecer com contexto do parágrafo (primeiras 200 chars)
                    ctx_snippet = contexto[:200] if len(contexto) > 200 else contexto
                    corpus.append({
                        "intent": intent,
                        "answer": answer,
                        "text": intent,
                        "contexto": ctx_snippet
                    })
    print(f"[+] SQuAD Completo: {len(corpus)} pares Q&A extraídos.")
    return corpus

# =============================================================================
# 2. Chat Base (PT) — Small Talk e Conhecimento Geral
# =============================================================================

CHAT_BASE_PT = [
    # Saudações
    {"intent": "oi", "answer": "Olá! Sou o CROM-LLM. Como posso ajudar?", "text": "oi"},
    {"intent": "ola", "answer": "Olá! Estou pronto.", "text": "ola"},
    {"intent": "bom dia", "answer": "Bom dia! Faça a sua pesquisa.", "text": "bom dia"},
    {"intent": "boa tarde", "answer": "Boa tarde! Em que posso ajudar?", "text": "boa tarde"},
    {"intent": "boa noite", "answer": "Boa noite! Pergunte-me algo.", "text": "boa noite"},
    {"intent": "hey", "answer": "Hey! Estou aqui para ajudar.", "text": "hey"},
    {"intent": "hello", "answer": "Hello! I am CROM-LLM.", "text": "hello"},
    {"intent": "hi", "answer": "Hi! Ask me anything.", "text": "hi"},
    # Identidade
    {"intent": "quem e voce", "answer": "Eu sou o CROM-LLM, um motor de busca semântica baseado em TF-IDF e PCA.", "text": "quem e voce"},
    {"intent": "qual o seu nome", "answer": "Meu nome é CROM-LLM.", "text": "qual o seu nome"},
    {"intent": "o que voce faz", "answer": "Processo perguntas e encontro respostas usando hashing semântico vetorial.", "text": "o que voce faz"},
    # Conhecimento PT
    {"intent": "o que e uma galaxia", "answer": "Uma galáxia é um enorme sistema gravitacionalmente ligado de estrelas, gás interestelar, poeira cósmica e matéria escura.", "text": "o que e uma galaxia"},
    {"intent": "o que e a gravidade", "answer": "A gravidade é a força fundamental que atrai dois corpos com massa um para o outro, descrita pela lei de Newton e pela Relatividade Geral de Einstein.", "text": "o que e a gravidade"},
    {"intent": "o que e um buraco negro", "answer": "Um buraco negro é uma região do espaço-tempo onde a gravidade é tão forte que nada, nem mesmo a luz, consegue escapar.", "text": "o que e um buraco negro"},
    {"intent": "o que e o universo", "answer": "O universo é a totalidade do espaço, tempo, matéria e energia que existe.", "text": "o que e o universo"},
    {"intent": "o que e a vida", "answer": "A vida é a condição que distingue organismos de matéria inorgânica, caracterizada por crescimento, reprodução e adaptação.", "text": "o que e a vida"},
    {"intent": "o que e inteligencia artificial", "answer": "Inteligência artificial é o campo da ciência da computação que cria sistemas capazes de realizar tarefas que normalmente requerem inteligência humana.", "text": "o que e inteligencia artificial"},
    {"intent": "o que e um computador", "answer": "Um computador é uma máquina programável que processa dados e executa instruções.", "text": "o que e um computador"},
    {"intent": "o que e a internet", "answer": "A internet é uma rede global de computadores interligados que permite a troca de informação.", "text": "o que e a internet"},
    {"intent": "o que e linux", "answer": "Linux é um sistema operativo de código aberto baseado no kernel criado por Linus Torvalds.", "text": "o que e linux"},
    {"intent": "o que e python", "answer": "Python é uma linguagem de programação de alto nível, conhecida pela sua simplicidade.", "text": "o que e python"},
    {"intent": "o que e go", "answer": "Go é uma linguagem de programação criada pela Google, focada em performance e concorrência.", "text": "o que e go"},
    {"intent": "qual a velocidade da luz", "answer": "A velocidade da luz no vácuo é aproximadamente 299.792.458 metros por segundo.", "text": "qual a velocidade da luz"},
    {"intent": "o que e matematica", "answer": "Matemática é a ciência que estuda quantidades, estruturas, espaços e mudanças.", "text": "o que e matematica"},
    {"intent": "o que e fisica", "answer": "Física é a ciência natural que estuda a matéria, a energia e as suas interações.", "text": "o que e fisica"},
]

# =============================================================================
# 3. NLP Pipeline
# =============================================================================

def remover_acentos(texto):
    nfkd = unicodedata.normalize('NFD', texto)
    return ''.join(c for c in nfkd if not unicodedata.combining(c))

def extrair_tokens(texto):
    texto = remover_acentos(texto).lower()
    return re.findall(r'\b[a-z0-9]+\b', texto)

def treinar_tf_idf(corpus, max_features=50000):
    print(f"[*] Treinando TF-IDF em {len(corpus)} amostras (max {max_features} features)...")
    doc_freq = defaultdict(int)
    doc_tokens = []
    
    for doc in corpus:
        tokens = extrair_tokens(doc["text"])
        doc_tokens.append(tokens)
        for token in set(tokens):
            doc_freq[token] += 1
            
    N = len(corpus)
    idf = {}
    for token, freq in doc_freq.items():
        idf[token] = math.log((N + 1) / (freq + 1)) + 1.0
        
    top_tokens = sorted(idf.items(), key=lambda x: x[1], reverse=True)
    vocabulario = {}
    for idx, (token, peso) in enumerate(top_tokens):
        if idx >= max_features:
            break
        vocabulario[token] = {"indice": idx, "idf": peso}
        
    return vocabulario, doc_tokens

def vetorizar(tokens, vocabulario, dim):
    vetor = [0.0] * dim
    counter = Counter(tokens)
    for token, count in counter.items():
        if token in vocabulario:
            info = vocabulario[token]
            vetor[info["indice"]] = count * info["idf"]
    norm = math.sqrt(sum(v*v for v in vetor))
    if norm > 0:
        vetor = [v/norm for v in vetor]
    return vetor

# =============================================================================
# 4. PCA Real (SVD via Power Iteration)
# =============================================================================

def treinar_pca_svd(corpus_vetores, dim_saida=64):
    N = len(corpus_vetores)
    D = len(corpus_vetores[0])
    print(f"[*] PCA Real (SVD): {N} vetores x {D} dims -> {dim_saida} componentes...")
    
    centroide = [0.0] * D
    for v in corpus_vetores:
        for j in range(D):
            centroide[j] += v[j]
    centroide = [c / N for c in centroide]
    
    centralizado = []
    for v in corpus_vetores:
        centralizado.append([v[j] - centroide[j] for j in range(D)])
    
    componentes = []
    for k in range(dim_saida):
        if k % 16 == 0:
            print(f"    Componente {k+1}/{dim_saida}...")
        
        random.seed(42 + k)
        w = [random.gauss(0, 1) for _ in range(D)]
        norm_w = math.sqrt(sum(x*x for x in w))
        w = [x / norm_w for x in w]
        
        for _ in range(15):
            projecoes = [sum(row[j] * w[j] for j in range(D)) for row in centralizado]
            novo_w = [0.0] * D
            for i_sample in range(N):
                p = projecoes[i_sample]
                for j in range(D):
                    novo_w[j] += centralizado[i_sample][j] * p
            norm_nw = math.sqrt(sum(x*x for x in novo_w))
            if norm_nw > 0:
                w = [x / norm_nw for x in novo_w]
        
        componentes.append(list(w))
        
        projecoes_final = [sum(row[j] * w[j] for j in range(D)) for row in centralizado]
        for i_sample in range(N):
            p = projecoes_final[i_sample]
            for j in range(D):
                centralizado[i_sample][j] -= p * w[j]
    
    print(f"[+] PCA: {len(componentes)} componentes extraídos.")
    return centroide, componentes

def calcular_hash_pca(vetor, hiperplanos, centroide):
    hash_val = 0
    for i, plano in enumerate(hiperplanos):
        dot = sum((v - c) * p for v, c, p in zip(vetor, centroide, plano))
        if dot > 0:
            hash_val |= (1 << i)
    return hash_val

# =============================================================================
# 5. Main
# =============================================================================

if __name__ == "__main__":
    print("=" * 60)
    print("  PESQUISA 7 — RAG MASSIVO: Ingestão Enciclopédica")
    print("=" * 60)
    
    # 1. Baixar SQuAD completo
    corpus_squad = baixar_squad_completo()
    
    # 2. Merge com Chat Base PT
    corpus_total = CHAT_BASE_PT + corpus_squad
    print(f"[+] Corpus total: {len(corpus_total)} amostras ({len(CHAT_BASE_PT)} PT + {len(corpus_squad)} SQuAD)")
    
    # 3. TF-IDF
    vocab_size = 50000
    vocabulario, doc_tokens = treinar_tf_idf(corpus_total, max_features=vocab_size)
    
    dim_entrada = min(len(vocabulario), vocab_size)
    with open("vocabulario.json", "w", encoding='utf-8') as f:
        json.dump(vocabulario, f)
    print(f"[+] Vocabulário: {len(vocabulario)} tokens.")
    
    # 4. Vetorizar e PCA
    print("[*] Vetorizando corpus completo...")
    corpus_vetores = [vetorizar(doc_tokens[i], vocabulario, dim_entrada) for i in range(len(corpus_total))]
    
    centroide, componentes = treinar_pca_svd(corpus_vetores, dim_saida=64)
    
    with open("matriz_pca_conversacional.json", "w", encoding='utf-8') as f:
        json.dump({
            "metadados": {"dimensao_entrada": dim_entrada, "dimensao_saida": 64, "algoritmo": "PCA_Real_SVD"},
            "centroide": centroide,
            "heads": {"entidade": componentes, "contexto": componentes, "visual": componentes}
        }, f)
    print("[+] Matrizes PCA exportadas.")
    
    # 5. Hashear tudo
    print("[*] Hasheando dataset completo...")
    dataset_vetorizado = []
    for i, doc in enumerate(corpus_total):
        v = corpus_vetores[i]
        h = calcular_hash_pca(v, componentes, centroide)
        dataset_vetorizado.append({
            "intent": doc["intent"],
            "answer": doc["answer"],
            "hash_entidade": h, "hash_contexto": h, "hash_visual": h
        })
        if (i+1) % 10000 == 0:
            print(f"    {i+1}/{len(corpus_total)} hasheados...")
    
    with open("dataset_vetorizado.json", "w", encoding='utf-8') as f:
        json.dump(dataset_vetorizado, f)
    print(f"[+] Dataset Final: {len(dataset_vetorizado)} entradas vetorizadas.")
    
    # 6. Teste de sanidade
    def testar(query, dataset, vocab, dim, comp, centro):
        tokens = extrair_tokens(query)
        v = vetorizar(tokens, vocab, dim)
        h = calcular_hash_pca(v, comp, centro)
        melhor_d, melhor_r = 99, ""
        for d in dataset:
            dist = bin(h ^ d["hash_entidade"]).count('1')
            if dist < melhor_d:
                melhor_d, melhor_r = dist, d["answer"]
        return melhor_d, melhor_r
    
    print("\n" + "=" * 60)
    print("  TESTES DE SANIDADE")
    print("=" * 60)
    for q in ["oi", "o que e uma galaxia", "what is gravity", "who was Einstein", "o que e a vida"]:
        d, r = testar(q, dataset_vetorizado, vocabulario, dim_entrada, componentes, centroide)
        print(f"  Q: '{q}' -> dist={d} bits -> {r[:80]}...")
    
    print(f"\n[✅] Pesquisa 7 concluída: {len(dataset_vetorizado)} entradas no RAG Massivo.")
