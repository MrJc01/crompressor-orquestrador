#!/bin/bash

echo -e "\033[33m[!] Inicializando Ambiente CROM-LLM v3...\033[0m"

# Verificar se a matriz PCA existe
MATRIZ_PATH="pesquisa6/data_engine/matriz_pca_conversacional.json"
if [ ! -f "$MATRIZ_PATH" ]; then
    echo -e "\033[33m[*] Matriz PCA não encontrada. Treinando cérebro inicial (Simulado)...\033[0m"
    cd pesquisa6/data_engine
    python3 ingestao_chat.py
    cd ../..
else
    echo -e "\033[32m[+] Cérebro PCA Conversacional detectado.\033[0m"
fi

echo -e "\033[33m[*] Compilando Motor Go (CROM-Chat)...\033[0m"
cd pesquisa6/cmd/crom-chat-v3
go build -o crom-chat main.go

if [ $? -eq 0 ]; then
    echo -e "\033[32m[+] Compilação bem-sucedida. Iniciando UX Terminal...\033[0m"
    sleep 1
    clear
    ./crom-chat
else
    echo -e "\033[31m[-] Erro ao compilar o motor CROM-Chat.\033[0m"
    exit 1
fi
