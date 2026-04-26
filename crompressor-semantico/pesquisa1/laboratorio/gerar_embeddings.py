import json
import random
import math

# Simula uma transformação em componentes principais.
# Criamos um "conceito_base" (vetor de N dimensões) para representar o núcleo do objeto.
def gerar_vetor_base(semente, dimensao=384):
    random.seed(semente)
    return [random.uniform(-1.0, 1.0) for _ in range(dimensao)]

# Adiciona ruído gaussiano (simulando variância de iluminação ou ordem de palavras)
def adicionar_ruido(vetor, intensidade=0.1):
    return [v + random.gauss(0, intensidade) for v in vetor]

def gerar_dados_maduros(caminho_saida):
    print("[*] Iniciando Pipeline de Geração de Embeddings Maduros (Simulação CLIP/SBERT)...")
    
    dados_finais = {
        "textos": [],
        "imagens": []
    }

    # ==========================================
    # 1. Dados de Texto (Chat Coerente)
    # ==========================================
    # "Dólar" é o nosso vetor base. Textos sinônimos recebem pouco ruído.
    vetor_dolar = gerar_vetor_base("conceito_dolar")
    
    dados_finais["textos"].append({
        "id": "T1",
        "texto": "Qual a cotação do dólar hoje?",
        "embedding": adicionar_ruido(vetor_dolar, 0.05) # Ruído muito baixo (mesma semântica)
    })
    dados_finais["textos"].append({
        "id": "T2",
        "texto": "Quanto vale o dólar agora?",
        "embedding": adicionar_ruido(vetor_dolar, 0.08) # Semântica similar
    })
    
    # Pergunta completamente diferente
    vetor_clima = gerar_vetor_base("conceito_clima")
    dados_finais["textos"].append({
        "id": "T3",
        "texto": "Vai chover amanhã?",
        "embedding": adicionar_ruido(vetor_clima, 0.05)
    })

    # ==========================================
    # 2. Dados de Visão (Patch-Hash)
    # ==========================================
    # Vamos gerar 16 patches (4x4) para uma imagem de maçã.
    # Assumimos que o objeto (maçã) ocupa os 4 patches centrais (índices 5, 6, 9, 10).
    # O resto é fundo (folhas, mesa, etc).
    vetor_fundo = gerar_vetor_base("fundo_madeira")
    vetor_maca_vermelha = gerar_vetor_base("maca_vermelha")
    vetor_maca_verde = gerar_vetor_base("maca_verde")
    
    # Criamos os 16 embeddings da "Imagem A" (Maçã Vermelha)
    patches_imagem_a = []
    patches_imagem_b = []
    
    indices_objeto = [5, 6, 9, 10]
    
    for i in range(16):
        if i in indices_objeto:
            # É a maçã!
            patches_imagem_a.append(adicionar_ruido(vetor_maca_vermelha, 0.05))
            patches_imagem_b.append(adicionar_ruido(vetor_maca_verde, 0.05))
        else:
            # É o fundo (mesmo fundo para ambas, simulando apenas mudança de cor do objeto)
            patches_imagem_a.append(adicionar_ruido(vetor_fundo, 0.02))
            patches_imagem_b.append(adicionar_ruido(vetor_fundo, 0.02))
            
    dados_finais["imagens"].append({
        "id": "IMG_A",
        "descricao": "Maçã Vermelha na Mesa",
        "patches": patches_imagem_a
    })
    
    dados_finais["imagens"].append({
        "id": "IMG_B",
        "descricao": "Maçã Verde na Mesa",
        "patches": patches_imagem_b
    })

    # Salva o Dataset Estruturado
    with open(caminho_saida, 'w', encoding='utf-8') as f:
        json.dump(dados_finais, f, ensure_ascii=False, indent=2)

    print(f"[+] Dataset Maduro Salvo em: {caminho_saida}")
    print(f"    Total de Textos: {len(dados_finais['textos'])}")
    print(f"    Total de Imagens (16 patches/cada): {len(dados_finais['imagens'])}")

if __name__ == "__main__":
    gerar_dados_maduros("../resultados/dataset_maduro.json")
