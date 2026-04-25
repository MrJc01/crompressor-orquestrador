#!/bin/bash
set -e
BENCH_DIR="/home/j/Documentos/GitHub/crom/papeis/benchmarks"
CROM_DIR="/home/j/Documentos/GitHub/crom/crompressor"
CODEBOOK="$BENCH_DIR/bench.cromdb"

echo "╔═══════════════════════════════════════════════════════════╗"
echo "║  BENCHMARK COMPARATIVO: CROM vs GZIP vs ZSTD            ║"
echo "╚═══════════════════════════════════════════════════════════╝"
echo ""

for dataset in dataset_json_5k.json dataset_logs_20k.log dataset_gocode.go; do
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
    echo "  GZIP-9:  $GZ_SIZE bytes ($(python3 -c "print(f'{$GZ_SIZE/$ORIG_SIZE*100:.1f}%')")) | $((T2-T))ms"
    
    # ZSTD
    T=$(($(date +%s%N)/1000000))
    zstd -19 -q "$ORIG" -o "$BENCH_DIR/${dataset}.zst" --force
    T2=$(($(date +%s%N)/1000000))
    ZST_SIZE=$(stat -c%s "$BENCH_DIR/${dataset}.zst")
    echo "  ZSTD-19: $ZST_SIZE bytes ($(python3 -c "print(f'{$ZST_SIZE/$ORIG_SIZE*100:.1f}%')")) | $((T2-T))ms"
    
    # CROM
    T=$(($(date +%s%N)/1000000))
    cd "$CROM_DIR"
    go run ./cmd/crompressor/ pack -i "$ORIG" -o "$BENCH_DIR/${dataset}.crom" -c "$CODEBOOK" --mode vault 2>&1 | grep -E '✔|Hit Rate|Entropy' || true
    T2=$(($(date +%s%N)/1000000))
    CROM_SIZE=$(stat -c%s "$BENCH_DIR/${dataset}.crom")
    echo "  CROM:    $CROM_SIZE bytes ($(python3 -c "print(f'{$CROM_SIZE/$ORIG_SIZE*100:.1f}%')")) | $((T2-T))ms"
    
    # Verify lossless
    go run ./cmd/crompressor/ unpack -i "$BENCH_DIR/${dataset}.crom" -o "$BENCH_DIR/${dataset}.restored" -c "$CODEBOOK" 2>&1 | grep '✔' || true
    if diff -q "$ORIG" "$BENCH_DIR/${dataset}.restored" > /dev/null 2>&1; then
        echo "  ✅ LOSSLESS ROUNDTRIP VERIFICADO"
    else
        echo "  ❌ ROUNDTRIP FALHOU"
    fi
    rm -f "$BENCH_DIR/${dataset}.restored"
    echo ""
done

echo "╔═══════════════════════════════════════════════════════════╗"
echo "║  BENCHMARK CONCLUÍDO                                     ║"
echo "╚═══════════════════════════════════════════════════════════╝"
