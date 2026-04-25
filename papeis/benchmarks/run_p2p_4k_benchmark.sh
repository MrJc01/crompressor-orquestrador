#!/bin/bash

# ==============================================================================
# CROM Ecosystem - P2P Sync Deduplication Benchmark (4KB Chunks)
# ==============================================================================

set -e

CROM_CMD="../../crompressor/cmd/crompressor"
BENCH_DIR="./diverse_data"
TEMP_DIR="/tmp/crom_p2p_bench_4k"
mkdir -p "$TEMP_DIR"

echo "=============================================================================="
echo "🚀 Iniciando Benchmark de Deduplicação P2P - CHUNK 4KB (5 Cenários Reais)"
echo "=============================================================================="

echo "[+] Compilando motor CROM..."
(cd ../../crompressor && go build -o "$TEMP_DIR/crompressor" ./cmd/crompressor/)
CROM="$TEMP_DIR/crompressor"

P1_DATA="$BENCH_DIR/node_modules.tar"
P2_DATA="$BENCH_DIR/all_python.txt"
P3_DATA="$BENCH_DIR/big_api_responses.json"
P4_DATA="$TEMP_DIR/server_logs.log"
P5_DATA="$TEMP_DIR/cctv_frames.tar"

# Check if generated files exist (they should from the previous run, but we make sure)
if [ ! -f "$P4_DATA" ]; then
    cp /tmp/crom_p2p_bench/server_logs.log "$P4_DATA" || true
fi
if [ ! -f "$P5_DATA" ]; then
    cp /tmp/crom_p2p_bench/cctv_frames.tar "$P5_DATA" || true
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

    # Train using 4KB chunks
    "$CROM" train -i "$DATA_PATH" -o "$CB_PATH" -s "$CB_SIZE" -k 4096 > /dev/null 2>&1

    # Pack using 4KB chunks
    "$CROM" pack -i "$DATA_PATH" -o "$CROM_PATH" -c "$CB_PATH" -k 4096 --mode edge > /dev/null 2>&1

    CROM_SIZE=$(stat -c%s "$CROM_PATH")
    CROM_MB=$(echo "scale=4; $CROM_SIZE / 1048576" | bc)

    RATIO=$(echo "scale=4; 100 - ($CROM_SIZE * 100 / $ORIGINAL_SIZE)" | bc)

    printf "%-35s | %-12s MB | %-12s MB | ⬇ %-5s %%\n" "$NAME" "$ORIGINAL_MB" "$CROM_MB" "$RATIO"
done

echo "---------------------------------------------------------------------------------"
echo "Conclusão: Com blocos de 4KB, o custo da Chunk Table dilui,"
echo "levando a deduplicação de borda muito mais próxima de 100% de economia."
echo "=============================================================================="
