#!/bin/bash
set -e
BENCH_DIR="/home/j/Documentos/GitHub/crom/papeis/benchmarks"
CROM_DIR="/home/j/Documentos/GitHub/crom/crompressor"
CODEBOOK="$BENCH_DIR/big.cromdb"

echo "╔═══════════════════════════════════════════════════════════╗"
echo "║  BENCHMARK V2: DADOS GRANDES + CODEBOOK 8192             ║"
echo "╚═══════════════════════════════════════════════════════════╝"
echo ""

for dataset in big_access_logs.log big_api_responses.json big_k8s_manifests.yaml; do
    ORIG="$BENCH_DIR/$dataset"
    ORIG_SIZE=$(stat -c%s "$ORIG")

    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "Dataset: $dataset ($(numfmt --to=iec $ORIG_SIZE))"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

    # GZIP
    T=$(($(date +%s%N)/1000000))
    gzip -c -9 "$ORIG" > "$BENCH_DIR/${dataset}.gz"
    T2=$(($(date +%s%N)/1000000))
    GZ_SIZE=$(stat -c%s "$BENCH_DIR/${dataset}.gz")
    echo "  GZIP-9:  $(numfmt --to=iec $GZ_SIZE) ($(python3 -c "print(f'{$GZ_SIZE/$ORIG_SIZE*100:.1f}%')")) | $((T2-T))ms"

    # ZSTD
    T=$(($(date +%s%N)/1000000))
    zstd -19 -q "$ORIG" -o "$BENCH_DIR/${dataset}.zst" --force
    T2=$(($(date +%s%N)/1000000))
    ZST_SIZE=$(stat -c%s "$BENCH_DIR/${dataset}.zst")
    echo "  ZSTD-19: $(numfmt --to=iec $ZST_SIZE) ($(python3 -c "print(f'{$ZST_SIZE/$ORIG_SIZE*100:.1f}%')")) | $((T2-T))ms"

    # CROM
    T=$(($(date +%s%N)/1000000))
    cd "$CROM_DIR"
    PACK_OUT=$(go run ./cmd/crompressor/ pack -i "$ORIG" -o "$BENCH_DIR/${dataset}.crom" -c "$CODEBOOK" --mode vault 2>&1)
    T2=$(($(date +%s%N)/1000000))
    CROM_SIZE=$(stat -c%s "$BENCH_DIR/${dataset}.crom")
    HIT=$(echo "$PACK_OUT" | grep "Hit Rate" | head -1)
    ENT=$(echo "$PACK_OUT" | grep "Entropy" | head -1)
    echo "  CROM:    $(numfmt --to=iec $CROM_SIZE) ($(python3 -c "print(f'{$CROM_SIZE/$ORIG_SIZE*100:.1f}%')")) | $((T2-T))ms"
    echo "  $HIT"
    echo "  $ENT"

    # Verify
    go run ./cmd/crompressor/ unpack -i "$BENCH_DIR/${dataset}.crom" -o "$BENCH_DIR/${dataset}.restored" -c "$CODEBOOK" 2>&1 > /dev/null
    if diff -q "$ORIG" "$BENCH_DIR/${dataset}.restored" > /dev/null 2>&1; then
        echo "  ✅ LOSSLESS VERIFICADO"
    else
        echo "  ❌ FALHOU"
    fi
    rm -f "$BENCH_DIR/${dataset}.restored"

    # Simular cenário de rede: segundo arquivo similar (apenas deltas mudam)
    echo ""
done
