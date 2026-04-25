#!/bin/bash
set -e
BENCH="/home/j/Documentos/GitHub/crom/papeis/benchmarks"
CROM_DIR="/home/j/Documentos/GitHub/crom/crompressor"
CB="$BENCH/large.cromdb"
RESULTS=""

echo "╔═══════════════════════════════════════════════════════════╗"
echo "║  BENCHMARK V4: ARQUIVOS GRANDES REAIS (200MB-500MB)      ║"
echo "╚═══════════════════════════════════════════════════════════╝"
echo ""

# Test files: real binary data from the system
declare -A FILES
FILES["GGUF-Model-469MB"]="/home/j/Documentos/GitHub/crom/crompressor-neuronio/pesquisas/modelos/qwen2.5-0.5b-instruct-q4_k_m.gguf"
FILES["LibTorch-431MB"]="/home/j/Documentos/GitHub/crom/crompressor-neuronio/pesquisa1/exemplos/.venv/lib/python3.12/site-packages/torch/lib/libtorch_cpu.so"
FILES["GGUF-261MB"]="/home/j/Documentos/GitHub/crom/crompressor-neuronio/pesquisas/modelos/qwen2.5-0.5b-q4_k_m.gguf"
FILES["CUDA-517MB"]="/home/j/Documentos/GitHub/crom/crompressor-neuronio/pesquisa1/exemplos/.venv/lib/python3.12/site-packages/nvidia/cu13/lib/libcublasLt.so.13"

for label in "GGUF-Model-469MB" "GGUF-261MB" "LibTorch-431MB" "CUDA-517MB"; do
    ORIG="${FILES[$label]}"
    [ -f "$ORIG" ] || { echo "SKIP: $label not found"; continue; }
    SZ=$(stat -c%s "$ORIG")
    SZ_MB=$((SZ/1024/1024))

    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "  $label (${SZ_MB}MB real)"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

    # GZIP
    echo -n "  GZIP-9:  "
    T=$(($(date +%s%N)/1000000))
    gzip -c -9 "$ORIG" > /tmp/bench_gz 2>&1
    T2=$(($(date +%s%N)/1000000))
    GZ=$(stat -c%s /tmp/bench_gz)
    GZ_PCT=$(python3 -c "print(f'{$GZ/$SZ*100:.1f}')")
    echo "$(numfmt --to=iec $GZ) (${GZ_PCT}%) | $((T2-T))ms"

    # ZSTD (level 3 for speed, not 19)
    echo -n "  ZSTD-3:  "
    T=$(($(date +%s%N)/1000000))
    zstd -3 -q "$ORIG" -o /tmp/bench_zst --force 2>&1
    T2=$(($(date +%s%N)/1000000))
    ZST=$(stat -c%s /tmp/bench_zst)
    ZST_PCT=$(python3 -c "print(f'{$ZST/$SZ*100:.1f}')")
    echo "$(numfmt --to=iec $ZST) (${ZST_PCT}%) | $((T2-T))ms"

    # ZSTD-19
    echo -n "  ZSTD-19: "
    T=$(($(date +%s%N)/1000000))
    zstd -19 -q "$ORIG" -o /tmp/bench_zst19 --force 2>&1
    T2=$(($(date +%s%N)/1000000))
    ZST19=$(stat -c%s /tmp/bench_zst19)
    ZST19_PCT=$(python3 -c "print(f'{$ZST19/$SZ*100:.1f}')")
    echo "$(numfmt --to=iec $ZST19) (${ZST19_PCT}%) | $((T2-T))ms"

    # CROM
    echo -n "  CROM:    "
    T=$(($(date +%s%N)/1000000))
    cd "$CROM_DIR"
    PACK_OUT=$(go run ./cmd/crompressor/ pack -i "$ORIG" -o /tmp/bench_crom -c "$CB" --mode vault --concurrency 8 2>&1)
    T2=$(($(date +%s%N)/1000000))
    CROM_SZ=$(stat -c%s /tmp/bench_crom)
    CROM_PCT=$(python3 -c "print(f'{$CROM_SZ/$SZ*100:.1f}')")
    echo "$(numfmt --to=iec $CROM_SZ) (${CROM_PCT}%) | $((T2-T))ms"
    echo "$PACK_OUT" | grep -E 'Hit Rate|Entropy' | sed 's/^/  /'

    # Verify roundtrip
    go run ./cmd/crompressor/ unpack -i /tmp/bench_crom -o /tmp/bench_restored -c "$CB" 2>&1 > /dev/null
    if diff -q "$ORIG" /tmp/bench_restored > /dev/null 2>&1; then
        echo "  ✅ LOSSLESS ROUNDTRIP OK"
    else
        echo "  ❌ ROUNDTRIP FAILED"
    fi

    rm -f /tmp/bench_gz /tmp/bench_zst /tmp/bench_zst19 /tmp/bench_crom /tmp/bench_restored
    echo ""
done

echo "╔═══════════════════════════════════════════════════════════╗"
echo "║  BENCHMARK CONCLUÍDO                                     ║"
echo "╚═══════════════════════════════════════════════════════════╝"
