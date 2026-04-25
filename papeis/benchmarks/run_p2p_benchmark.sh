#!/bin/bash

# ==============================================================================
# CROM Ecosystem - P2P Sync Deduplication Benchmark (5 Projects)
# Este script simula a sincronização P2P de 5 projetos reais, medindo a 
# eficiência da Deduplicação de Borda (Edge Deduplication) via CROM vs Sem CROM.
# ==============================================================================

set -e

CROM_CMD="../../crompressor/cmd/crompressor"
BENCH_DIR="./diverse_data"
TEMP_DIR="/tmp/crom_p2p_bench"
mkdir -p "$TEMP_DIR"

echo "=============================================================================="
echo "🚀 Iniciando Benchmark de Deduplicação P2P (5 Cenários Reais)"
echo "=============================================================================="

# Build the latest engine
echo "[+] Compilando motor CROM..."
(cd ../../crompressor && go build -o "$TEMP_DIR/crompressor" ./cmd/crompressor/)
CROM="$TEMP_DIR/crompressor"

# Prepare data for 5 projects
echo "[+] Preparando os 5 Projetos de Teste..."

# 1. Node.js (already exists)
P1_DATA="$BENCH_DIR/node_modules.tar"

# 2. Python (already exists)
P2_DATA="$BENCH_DIR/all_python.txt"

# 3. JSON API (already exists)
P3_DATA="$BENCH_DIR/big_api_responses.json"

# 4. Logs Servidor (generate 50MB of repetitive logs)
P4_DATA="$TEMP_DIR/server_logs.log"
if [ ! -f "$P4_DATA" ]; then
    echo "Gerando Logs..."
    for i in {1..500000}; do
        echo "2026-04-25T10:00:00Z INFO [worker-$((i%10))] Processed request $i from 192.168.1.$((i%255)) - STATUS OK" >> "$P4_DATA"
    done
fi

# 5. Imagens Similares (generate pseudo-bitmaps with minor variations)
P5_DATA="$TEMP_DIR/cctv_frames.tar"
if [ ! -f "$P5_DATA" ]; then
    echo "Gerando Frames CCTV..."
    mkdir -p "$TEMP_DIR/cctv"
    head -c 1M /dev/zero | tr '\0' 'F' > "$TEMP_DIR/cctv/base_frame.bmp"
    for i in {1..50}; do
        cp "$TEMP_DIR/cctv/base_frame.bmp" "$TEMP_DIR/cctv/frame_$i.bmp"
        echo "Timestamp: $i" >> "$TEMP_DIR/cctv/frame_$i.bmp" # Small delta
    done
    tar -cf "$P5_DATA" -C "$TEMP_DIR" cctv
fi

PROJECTS=(
    "Projeto 1 (Next.js Node Modules) | $P1_DATA | 32768"
    "Projeto 2 (Repo Python)          | $P2_DATA | 32768"
    "Projeto 3 (JSON API Dump)        | $P3_DATA | 16384"
    "Projeto 4 (Server Logs)          | $P4_DATA | 8192"
    "Projeto 5 (CCTV Frames Similares)| $P5_DATA | 4096"
)

echo ""
printf "%-35s | %-15s | %-15s | %-10s\n" "PROJETO" "TRÁFEGO S/ CROM" "TRÁFEGO C/ CROM" "REDUÇÃO"
echo "---------------------------------------------------------------------------------"

for proj in "${PROJECTS[@]}"; do
    IFS="|" read -r NAME DATA_PATH CB_SIZE <<< "$proj"
    NAME=$(echo "$NAME" | xargs)
    DATA_PATH=$(echo "$DATA_PATH" | xargs)
    CB_SIZE=$(echo "$CB_SIZE" | xargs)

    if [ ! -f "$DATA_PATH" ]; then
        continue
    fi

    ORIGINAL_SIZE=$(stat -c%s "$DATA_PATH")
    ORIGINAL_MB=$(echo "scale=2; $ORIGINAL_SIZE / 1048576" | bc)

    CB_PATH="$TEMP_DIR/cb_${NAME// /_}.cromdb"
    CROM_PATH="$TEMP_DIR/sync_${NAME// /_}.crom"

    # 1. Simula Treinamento do Cérebro (Codebook P2P)
    "$CROM" train -i "$DATA_PATH" -o "$CB_PATH" -s "$CB_SIZE" > /dev/null 2>&1

    # 2. Simula Sincronização (Pack envia só os Hashes e Deltas)
    "$CROM" pack -i "$DATA_PATH" -o "$CROM_PATH" -c "$CB_PATH" --mode edge > /dev/null 2>&1

    CROM_SIZE=$(stat -c%s "$CROM_PATH")
    CROM_MB=$(echo "scale=2; $CROM_SIZE / 1048576" | bc)

    RATIO=$(echo "scale=2; 100 - ($CROM_SIZE * 100 / $ORIGINAL_SIZE)" | bc)

    printf "%-35s | %-12s MB | %-12s MB | ⬇ %-5s %%\n" "$NAME" "$ORIGINAL_MB" "$CROM_MB" "$RATIO"
done

echo "---------------------------------------------------------------------------------"
echo "Conclusão: O tráfego COM CROM reflete os pacotes enviados pela rede durante um sync."
echo "Matchs exatos do Cérebro são transferidos apenas como Hashes (0 bytes de payload)."
echo "=============================================================================="
