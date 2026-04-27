import json
import re
import math
import unicodedata
from collections import defaultdict, Counter
import random
import sys

try:
    from datasets import load_dataset
    HF_AVAILABLE = True
except ImportError:
    HF_AVAILABLE = False

# =============================================================================
# CHAT_BASE MASSIVO: ~200 variações de small-talk em PT
# Cada entrada tem: intent (pergunta), answer (resposta), text (tokens de treino)
# O campo "text" deve conter variações e sinónimos para enriquecer o TF-IDF.
# =============================================================================

SAUDACOES = [
    {"intent": "oi", "answer": "Olá! Sou o CROM-LLM. Como posso ajudar?", "text": "oi"},
    {"intent": "ola", "answer": "Olá! Estou pronto para as tuas perguntas.", "text": "ola"},
    {"intent": "oie", "answer": "Oie! Tudo bem? Pergunte-me algo!", "text": "oie"},
    {"intent": "oii", "answer": "Oii! Em que posso ajudar?", "text": "oii"},
    {"intent": "hey", "answer": "Hey! Estou aqui para ajudar.", "text": "hey"},
    {"intent": "e ai", "answer": "E aí! Manda a tua pergunta.", "text": "e ai"},
    {"intent": "fala", "answer": "Fala! Estou a ouvir.", "text": "fala"},
    {"intent": "salve", "answer": "Salve! Como posso ser útil?", "text": "salve"},
    {"intent": "bom dia", "answer": "Bom dia! Faça a sua pesquisa.", "text": "bom dia"},
    {"intent": "boa tarde", "answer": "Boa tarde! Em que posso ajudar?", "text": "boa tarde"},
    {"intent": "boa noite", "answer": "Boa noite! Pergunte-me algo.", "text": "boa noite"},
    {"intent": "hello", "answer": "Hello! I am CROM-LLM, ready to help.", "text": "hello"},
    {"intent": "hi", "answer": "Hi! Ask me anything.", "text": "hi"},
    {"intent": "boas", "answer": "Boas! Em que posso ajudar?", "text": "boas"},
    {"intent": "opa", "answer": "Opa! Estou aqui.", "text": "opa"},
    {"intent": "eae", "answer": "Eae! Manda ver a pergunta.", "text": "eae"},
    {"intent": "yo", "answer": "Yo! Pergunte-me algo.", "text": "yo"},
]

# Gerar variações automáticas de saudações (duplicar com prefixos/sufixos)
_saud_extras = []
for s in SAUDACOES:
    _saud_extras.append({"intent": s["intent"] + " tudo bem", "answer": s["answer"], "text": s["text"] + " tudo bem"})
    _saud_extras.append({"intent": s["intent"] + " como vai", "answer": s["answer"], "text": s["text"] + " como vai"})
SAUDACOES.extend(_saud_extras)

DESPEDIDAS = [
    {"intent": "tchau", "answer": "Tchau! Até a próxima.", "text": "tchau"},
    {"intent": "ate logo", "answer": "Até logo! Foi bom conversar.", "text": "ate logo"},
    {"intent": "ate mais", "answer": "Até mais! Volte sempre.", "text": "ate mais"},
    {"intent": "adeus", "answer": "Adeus! Boa sorte.", "text": "adeus"},
    {"intent": "falou", "answer": "Falou! Até breve.", "text": "falou"},
    {"intent": "bye", "answer": "Bye! See you.", "text": "bye"},
    {"intent": "flw", "answer": "Flw! Até a próxima.", "text": "flw"},
    {"intent": "vlw flw", "answer": "Valeu, falou! Até mais.", "text": "vlw flw valeu"},
    {"intent": "boa noite tchau", "answer": "Boa noite! Descanse bem.", "text": "boa noite tchau despedida"},
    {"intent": "obrigado tchau", "answer": "De nada! Tchau!", "text": "obrigado tchau despedida"},
]

IDENTIDADE = [
    {"intent": "quem e voce", "answer": "Eu sou o CROM-LLM, uma IA de busca vetorial baseada em SimHash e TF-IDF.", "text": "quem e voce"},
    {"intent": "qual o seu nome", "answer": "Meu nome é CROM-LLM, um motor de busca semântica.", "text": "qual o seu nome"},
    {"intent": "o que voce faz", "answer": "Eu processo perguntas e encontro respostas usando hashing semântico.", "text": "o que voce faz"},
    {"intent": "como voce funciona", "answer": "Funciono com TF-IDF, SimHash e projeção PCA para encontrar a resposta mais próxima.", "text": "como voce funciona"},
    {"intent": "voce e uma ia", "answer": "Sim, sou uma inteligência artificial local, sem dependência de cloud.", "text": "voce e uma ia inteligencia artificial"},
    {"intent": "voce e um bot", "answer": "Sou um motor de inferência vetorial, não um chatbot tradicional.", "text": "voce e um bot robo"},
    {"intent": "voce e humano", "answer": "Não, sou uma IA. Processo vetores, não emoções.", "text": "voce e humano pessoa"},
    {"intent": "quem te criou", "answer": "Fui criado como parte do projecto CROM de compressão semântica.", "text": "quem te criou fez construiu"},
    {"intent": "qual a sua funcao", "answer": "A minha função é encontrar a resposta mais semanticamente próxima à sua pergunta.", "text": "qual a sua funcao proposito objetivo"},
    {"intent": "voce pensa", "answer": "Não penso no sentido humano. Calculo distâncias de Hamming entre vetores.", "text": "voce pensa raciocina"},
]

AGRADECIMENTOS = [
    {"intent": "obrigado", "answer": "De nada! Estou aqui para ajudar.", "text": "obrigado"},
    {"intent": "obrigada", "answer": "De nada! Sempre às ordens.", "text": "obrigada"},
    {"intent": "valeu", "answer": "Valeu! Qualquer coisa, pergunte.", "text": "valeu"},
    {"intent": "thanks", "answer": "You're welcome!", "text": "thanks thank you"},
    {"intent": "muito obrigado", "answer": "Por nada! Foi um prazer ajudar.", "text": "muito obrigado"},
    {"intent": "brigado", "answer": "De nada!", "text": "brigado"},
]

COMO_VAI = [
    {"intent": "tudo bem", "answer": "Tudo ótimo! E contigo?", "text": "tudo bem"},
    {"intent": "como vai", "answer": "Vou bem, obrigado! Em que posso ajudar?", "text": "como vai"},
    {"intent": "como esta", "answer": "Estou operacional e pronto para pesquisas!", "text": "como esta"},
    {"intent": "como voce esta", "answer": "Estou a funcionar a 100%. Pergunte-me algo!", "text": "como voce esta"},
    {"intent": "tudo certo", "answer": "Tudo certo! Manda a pergunta.", "text": "tudo certo"},
    {"intent": "beleza", "answer": "Beleza! O que precisa?", "text": "beleza"},
]

AJUDA = [
    {"intent": "me ajuda", "answer": "Claro! Diga-me o que precisa.", "text": "me ajuda"},
    {"intent": "preciso de ajuda", "answer": "Estou aqui para ajudar! Qual é a sua dúvida?", "text": "preciso de ajuda"},
    {"intent": "pode me ajudar", "answer": "Com certeza! Faça a sua pergunta.", "text": "pode me ajudar"},
    {"intent": "help", "answer": "I'm here to help! What do you need?", "text": "help"},
    {"intent": "socorro", "answer": "Calma! Diga-me o que está a acontecer.", "text": "socorro"},
    {"intent": "nao entendo", "answer": "Sem problema! Reformule a sua pergunta e eu tento novamente.", "text": "nao entendo"},
    {"intent": "como usar", "answer": "Basta escrever a sua pergunta e eu procuro a resposta mais próxima.", "text": "como usar"},
]

AFIRMACOES = [
    {"intent": "sim", "answer": "Entendido! Continue.", "text": "sim"},
    {"intent": "nao", "answer": "Ok, entendido.", "text": "nao"},
    {"intent": "ok", "answer": "Ok! Mais alguma coisa?", "text": "ok"},
    {"intent": "certo", "answer": "Certo! Prossiga.", "text": "certo"},
    {"intent": "claro", "answer": "Claro! Em que mais posso ajudar?", "text": "claro"},
    {"intent": "entendi", "answer": "Ótimo! Pergunte mais se precisar.", "text": "entendi"},
    {"intent": "ta bom", "answer": "Perfeito!", "text": "ta bom"},
]

CONHECIMENTO = [
    # Cada entrada usa text = a frase real que o utilizador escreveria (paridade com Go)
    {"intent": "o que e o universo", "answer": "O universo é a totalidade do espaço, tempo, matéria e energia que existe.", "text": "o que e o universo"},
    {"intent": "universo cosmos", "answer": "O universo é a totalidade do espaço, tempo, matéria e energia que existe.", "text": "universo cosmos espaco tempo materia energia"},
    {"intent": "o que e inteligencia artificial", "answer": "Inteligência artificial é o campo da ciência da computação que cria sistemas capazes de realizar tarefas que normalmente requerem inteligência humana.", "text": "o que e inteligencia artificial"},
    {"intent": "ia machine learning", "answer": "Inteligência artificial é o campo da ciência da computação que cria sistemas capazes de realizar tarefas que normalmente requerem inteligência humana.", "text": "ia computacao maquina aprendizado"},
    {"intent": "o que e um computador", "answer": "Um computador é uma máquina programável que processa dados e executa instruções.", "text": "o que e um computador"},
    {"intent": "o que e a internet", "answer": "A internet é uma rede global de computadores interligados que permite a troca de informação.", "text": "o que e a internet"},
    {"intent": "o que e programacao", "answer": "Programação é o processo de criar instruções para um computador executar.", "text": "o que e programacao"},
    {"intent": "o que e machine learning", "answer": "Machine Learning é um subcampo da IA onde sistemas aprendem padrões a partir de dados.", "text": "o que e machine learning"},
    {"intent": "o que e python", "answer": "Python é uma linguagem de programação de alto nível, conhecida pela sua simplicidade.", "text": "o que e python"},
    {"intent": "o que e go", "answer": "Go (Golang) é uma linguagem de programação criada pela Google, focada em performance e concorrência.", "text": "o que e go"},
    {"intent": "o que e linux", "answer": "Linux é um sistema operativo de código aberto baseado no kernel criado por Linus Torvalds.", "text": "o que e linux"},
    {"intent": "o que e uma rede neural", "answer": "Uma rede neural é um modelo computacional inspirado no cérebro humano, usado em IA.", "text": "o que e uma rede neural"},
    {"intent": "o que e compressao de dados", "answer": "Compressão de dados é o processo de reduzir o tamanho de ficheiros sem perder informação essencial.", "text": "o que e compressao de dados"},
    {"intent": "o que e um hash", "answer": "Um hash é uma função que converte dados de qualquer tamanho numa impressão digital de tamanho fixo.", "text": "o que e um hash"},
    {"intent": "o que e simhash", "answer": "SimHash é um algoritmo de Locality-Sensitive Hashing que preserva a similaridade entre documentos.", "text": "o que e simhash"},
    {"intent": "o que e tfidf", "answer": "TF-IDF é uma medida estatística que avalia a importância de uma palavra num documento em relação a um corpus.", "text": "o que e tfidf"},
    {"intent": "quem e linus torvalds", "answer": "Linus Torvalds é o criador do kernel Linux e do sistema de controle de versão Git.", "text": "quem e linus torvalds"},
    {"intent": "o que e git", "answer": "Git é um sistema de controle de versão distribuído para rastrear mudanças em código-fonte.", "text": "o que e git"},
    {"intent": "o que e uma api", "answer": "API (Application Programming Interface) é uma interface que permite a comunicação entre sistemas de software.", "text": "o que e uma api"},
    {"intent": "o que e a lua", "answer": "A Lua é o satélite natural da Terra, com cerca de 3.474 km de diâmetro.", "text": "o que e a lua"},
    {"intent": "o que e o sol", "answer": "O Sol é a estrela central do nosso sistema solar, uma esfera de plasma.", "text": "o que e o sol"},
    {"intent": "o que e a gravidade", "answer": "A gravidade é a força fundamental que atrai dois corpos com massa um para o outro.", "text": "o que e a gravidade"},
    {"intent": "o que e matematica", "answer": "Matemática é a ciência que estuda quantidades, estruturas, espaços e mudanças.", "text": "o que e matematica"},
    {"intent": "o que e fisica", "answer": "Física é a ciência natural que estuda a matéria, a energia e as suas interações.", "text": "o que e fisica"},
    {"intent": "o que e biologia", "answer": "Biologia é a ciência que estuda os seres vivos e os seus processos vitais.", "text": "o que e biologia"},
    {"intent": "quanto e 1 mais 1", "answer": "1 + 1 = 2.", "text": "quanto e 1 mais 1"},
    {"intent": "qual a velocidade da luz", "answer": "A velocidade da luz no vácuo é aproximadamente 299.792.458 metros por segundo.", "text": "qual a velocidade da luz"},
]

# Montar o CHAT_BASE completo
CHAT_BASE = SAUDACOES + DESPEDIDAS + IDENTIDADE + AGRADECIMENTOS + COMO_VAI + AJUDA + AFIRMACOES + CONHECIMENTO

# Out-Of-Distribution (OOD)
DESAFIOS_CEGOS = [
    {"intent": "Qual é a velocidade de dobra da USS Enterprise?", "text": "velocidade de dobra uss enterprise star trek"},
    {"intent": "Pode criar uma receita de bolo de chocolate com bacon?", "text": "receita bolo chocolate bacon culinaria estranha"},
    {"intent": "O que significa SRE em Engenharia de Software?", "text": "sre site reliability engineering devops"},
    {"intent": "Como derrotar o boss final do Elden Ring?", "text": "derrotar boss final elden ring jogo"},
    {"intent": "Qual a cotação do dólar hoje em Tóquio?", "text": "cotacao dolar toquio iene japao economia"},
    {"intent": "Existe vida extraterrestre inteligente em Alpha Centauri?", "text": "vida extraterrestre alienigenas alpha centauri espaco"},
    {"intent": "o meu nome é jorge e sou um humano", "text": "meu nome e jorge humano"},
    {"intent": "qual o sentido da vida", "text": "qual o sentido da vida filosofia viver"},
    {"intent": "?!?!", "text": ""},
] * 5

# =============================================================================
# Funções de processamento
# =============================================================================

def remover_acentos(texto):
    """Remove acentos para paridade com Go (removerAcentos)."""
    nfkd = unicodedata.normalize('NFD', texto)
    return ''.join(c for c in nfkd if not unicodedata.combining(c))

def extrair_tokens(texto):
    texto = remover_acentos(texto).lower()
    palavras = re.findall(r'\b[a-z0-9]+\b', texto)
    tokens = list(palavras)
    for i in range(len(palavras) - 1):
        tokens.append(f"{palavras[i]}_{palavras[i+1]}")
    for i in range(len(palavras) - 2):
        tokens.append(f"{palavras[i]}_{palavras[i+1]}_{palavras[i+2]}")
    return tokens

def treinar_tf_idf(corpus, max_features=20000):
    print(f"[*] Treinando TF-IDF em {len(corpus)} amostras...")
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

def vetorizar(tokens, vocabulario, max_features):
    vetor = [0.0] * max_features
    counter = Counter(tokens)
    for token, count in counter.items():
        if token in vocabulario:
            info = vocabulario[token]
            vetor[info["indice"]] = count * info["idf"]
            
    norm = math.sqrt(sum(v*v for v in vetor))
    if norm > 0:
        vetor = [v/norm for v in vetor]
    return vetor

def gerar_matriz_fallback(dim_entrada, dim_saida, seed=42):
    """Fallback aleatório, mantido apenas se o SVD falhar."""
    random.seed(seed)
    return [[random.uniform(-1.0, 1.0) for _ in range(dim_entrada)] for _ in range(dim_saida)]

def treinar_pca_svd(corpus_vetores, dim_saida=64):
    """
    PCA REAL via SVD (Singular Value Decomposition).
    1. Calcula o centróide (média) de todos os vetores.
    2. Centraliza os dados subtraindo a média.
    3. Calcula a matriz de covariância.
    4. Extrai os 'dim_saida' autovetores de maior variância via Power Iteration.
    Retorna: (centroide, componentes_principais)
    """
    N = len(corpus_vetores)
    D = len(corpus_vetores[0])
    print(f"[*] Treinando PCA Real (SVD): {N} vetores x {D} dimensões -> {dim_saida} componentes...")
    
    # 1. Centróide (Média de cada dimensão)
    centroide = [0.0] * D
    for v in corpus_vetores:
        for j in range(D):
            centroide[j] += v[j]
    centroide = [c / N for c in centroide]
    
    # 2. Centralizar (subtrair a média)
    centralizado = []
    for v in corpus_vetores:
        centralizado.append([v[j] - centroide[j] for j in range(D)])
    
    # 3. Power Iteration para extrair os top-k autovetores
    # (Implementação pura sem numpy para manter o espírito do projeto)
    componentes = []
    for k in range(dim_saida):
        if k % 16 == 0:
            print(f"    Componente {k+1}/{dim_saida}...")
        
        # Vetor aleatório inicial
        random.seed(42 + k)
        w = [random.gauss(0, 1) for _ in range(D)]
        
        # Normalizar
        norm_w = math.sqrt(sum(x*x for x in w))
        w = [x / norm_w for x in w]
        
        # 20 iterações de Power Method: w = (X^T X) w / ||(X^T X) w||
        for _ in range(20):
            # Xw (projeção de cada amostra no vetor w)
            projecoes = []
            for row in centralizado:
                projecoes.append(sum(row[j] * w[j] for j in range(D)))
            
            # X^T (Xw) (reconstruir o vetor no espaço original)
            novo_w = [0.0] * D
            for i_sample in range(N):
                p = projecoes[i_sample]
                for j in range(D):
                    novo_w[j] += centralizado[i_sample][j] * p
            
            # Normalizar
            norm_nw = math.sqrt(sum(x*x for x in novo_w))
            if norm_nw > 0:
                w = [x / norm_nw for x in novo_w]
        
        componentes.append(list(w))
        
        # Deflação: remover a variância capturada por este componente
        projecoes_final = []
        for row in centralizado:
            projecoes_final.append(sum(row[j] * w[j] for j in range(D)))
        
        for i_sample in range(N):
            p = projecoes_final[i_sample]
            for j in range(D):
                centralizado[i_sample][j] -= p * w[j]
    
    print(f"[+] PCA Real concluído: {len(componentes)} componentes de variância máxima extraídos.")
    return centroide, componentes

def calcular_hash_pca(vetor, hiperplanos, centroide=None):
    """Projeta o vetor nos hiperplanos (PCA). Se houver centróide, subtrai antes."""
    hash_val = 0
    for i, plano in enumerate(hiperplanos):
        if centroide:
            dot = sum((v - c) * p for v, c, p in zip(vetor, centroide, plano))
        else:
            dot = sum(v * p for v, p in zip(vetor, plano))
        if dot > 0:
            hash_val |= (1 << i)
    return hash_val

def calcular_centroide_saudacoes(vocabulario, dim_entrada, matriz_ent):
    saudacoes_texto = ["oi", "ola", "oie", "oii", "hey", "bom dia", "boa tarde", "boa noite",
                        "salve", "fala", "e ai", "boas", "opa", "eae", "hello", "hi"]
    vetores = []
    for s in saudacoes_texto:
        tokens = extrair_tokens(s)
        v = vetorizar(tokens, vocabulario, dim_entrada)
        if any(x != 0.0 for x in v):
            vetores.append(v)
    
    if not vetores:
        return 0
        
    v_medio = [sum(col)/len(vetores) for col in zip(*vetores)]
    return calcular_hash_pca(v_medio, matriz_ent)

def baixar_squad_hf(num_amostras=5000):
    print(f"[*] Baixando SQuAD via HuggingFace ({num_amostras} amostras)...")
    try:
        dataset = load_dataset("squad", split="train", streaming=True)
        corpus = []
        for qa in dataset:
            intent = qa['question']
            answer = qa['answers']['text'][0] if len(qa['answers']['text']) > 0 else ""
            if answer:
                corpus.append({"intent": intent, "answer": answer, "text": intent})
            if len(corpus) >= num_amostras:
                break
        return corpus
    except Exception as e:
        print(f"[-] Erro ao baixar SQuAD: {e}")
        return []

# =============================================================================
# Main
# =============================================================================

if __name__ == "__main__":
    online = "--online" in sys.argv

    if online and HF_AVAILABLE:
        corpus_squad = baixar_squad_hf(5000)
        corpus_treino = CHAT_BASE + corpus_squad
    else:
        if online:
            print("[-] datasets não instalado. Usando modo offline.")
        corpus_treino = list(CHAT_BASE)
    
    print(f"[+] Corpus total: {len(corpus_treino)} amostras ({len(CHAT_BASE)} chat-base).")
        
    vocab_size = 20000
    vocabulario, doc_tokens = treinar_tf_idf(corpus_treino, max_features=vocab_size)
    
    with open("vocabulario.json", "w", encoding='utf-8') as f:
        json.dump(vocabulario, f, indent=2)
    print(f"[+] Vocabulário: {len(vocabulario)} tokens mapeados.")
    
    dim_entrada = min(len(vocabulario), vocab_size)
    
    # =====================================================================
    # FASE CRÍTICA: Vetorizar TODO o corpus para treinar o PCA Real (SVD)
    # =====================================================================
    print("[*] Vetorizando corpus completo para SVD...")
    corpus_vetores = []
    for i, doc in enumerate(corpus_treino):
        v = vetorizar(doc_tokens[i], vocabulario, dim_entrada)
        corpus_vetores.append(v)
    
    # Treinar PCA Real via Power Iteration (SVD puro em Python)
    centroide, componentes_pca = treinar_pca_svd(corpus_vetores, dim_saida=64)
    
    # Usar os componentes PCA reais como hiperplanos para TODAS as cabeças
    # (cada cabeça recebe os mesmos eixos de variância máxima)
    matriz_ent = componentes_pca
    matriz_ctx = componentes_pca  # Contexto e Visual partilham os mesmos eixos por agora
    matriz_vis = componentes_pca
    
    with open("matriz_pca_conversacional.json", "w", encoding='utf-8') as f:
        json.dump({
            "metadados": {
                "dimensao_entrada": dim_entrada, 
                "dimensao_saida": 64,
                "algoritmo": "PCA_Real_SVD_PowerIteration"
            },
            "centroide": centroide,
            "heads": {"entidade": matriz_ent, "contexto": matriz_ctx, "visual": matriz_vis}
        }, f, indent=2)
    print("[+] Matrizes PCA Reais e Centróide exportados.")
    
    # Vetorizar e hashear o dataset com os novos hiperplanos PCA
    dataset_vetorizado = []
    for i, doc in enumerate(corpus_treino):
        v = corpus_vetores[i]
        hEnt = calcular_hash_pca(v, matriz_ent, centroide)
        hCtx = calcular_hash_pca(v, matriz_ctx, centroide)
        hVis = calcular_hash_pca(v, matriz_vis, centroide)
        
        entry = {
            "intent": doc["intent"], 
            "answer": doc.get("answer", ""),
            "hash_entidade": hEnt, "hash_contexto": hCtx, "hash_visual": hVis
        }
        dataset_vetorizado.append(entry)
        
    with open("dataset_vetorizado.json", "w", encoding='utf-8') as f:
        json.dump(dataset_vetorizado, f, indent=2)
    print(f"[+] Dataset Vetorizado: {len(dataset_vetorizado)} itens.")
    
    with open("desafio_cego.json", "w", encoding='utf-8') as f:
        json.dump(DESAFIOS_CEGOS, f, indent=2)
    print(f"[+] {len(DESAFIOS_CEGOS)} Desafios Cegos OOD guardados.")
    
    # =====================================================================
    # TESTE DE SANIDADE: Verificar distâncias com PCA Real
    # =====================================================================
    toks_oi = extrair_tokens("oi")
    vec_oi = vetorizar(toks_oi, vocabulario, dim_entrada)
    hash_oi_query = calcular_hash_pca(vec_oi, matriz_ent, centroide)
    hash_oi_db = dataset_vetorizado[0]["hash_entidade"]
    dist_oi = bin(hash_oi_query ^ hash_oi_db).count('1')
    print(f"\n[SANIDADE] Query 'oi' vs DB[0] '{dataset_vetorizado[0]['intent']}': dist={dist_oi} bits")
    
    toks_univ = extrair_tokens("o que e o universo")
    vec_univ = vetorizar(toks_univ, vocabulario, dim_entrada)
    hash_univ_query = calcular_hash_pca(vec_univ, matriz_ent, centroide)
    for entry in dataset_vetorizado:
        if "universo" in entry["intent"]:
            dist_u = bin(hash_univ_query ^ entry["hash_entidade"]).count('1')
            print(f"[SANIDADE] Query 'o que e o universo' vs DB '{entry['intent']}': dist={dist_u} bits")
            break
    
    # Teste cruzado: "oi" vs "o que e o universo" (devem estar LONGE)
    dist_cross = bin(hash_oi_query ^ hash_univ_query).count('1')
    print(f"[SANIDADE] 'oi' vs 'universo' (devem divergir): dist={dist_cross} bits")

