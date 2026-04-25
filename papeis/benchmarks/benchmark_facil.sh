#!/bin/bash
# ╔═════════════════════════════════════════════════════════╗
# ║  CROM Benchmark — Copie e rode você mesmo!             ║
# ║  Compara Crompressor vs GZIP vs ZSTD no SEU dataset    ║
# ╚═════════════════════════════════════════════════════════╝
#
# Uso:
#   chmod +x benchmark_facil.sh
#   ./benchmark_facil.sh /caminho/do/seu/arquivo
#
# Pré-requisitos: go 1.21+, gzip, zstd, git
set -e

if [ -z "$1" ]; then
    echo "Uso: $0 <arquivo_para_testar>"
    echo ""
    echo "Exemplos:"
    echo "  $0 /var/log/syslog"
    echo "  $0 ~/dados/backup.tar"
    echo "  $0 ~/modelo.gguf"
    exit 1
fi

INPUT="$1"
if [ ! -f "$INPUT" ]; then
    echo "❌ Arquivo não encontrado: $INPUT"
    exit 1
fi

ORIG_SIZE=$(stat -c%s "$INPUT")
ORIG_MB=$((ORIG_SIZE/1024/1024))
WORKDIR=$(mktemp -d)

echo "╔═══════════════════════════════════════════════════╗"
echo "║  CROM Benchmark                                  ║"
echo "╠═══════════════════════════════════════════════════╣"
echo "║  Arquivo: $(basename "$INPUT")"
echo "║  Tamanho: ${ORIG_MB}MB ($ORIG_SIZE bytes)"
echo "╚═══════════════════════════════════════════════════╝"
echo ""

# 1. Clonar e buildar CROM (se não existe)
CROM_BIN="$WORKDIR/crompressor"
if ! command -v crompressor &>/dev/null; then
    echo "⏳ Buildando Crompressor..."
    git clone --depth 1 https://github.com/MrJc01/crompressor "$WORKDIR/crom-src" 2>/dev/null
    cd "$WORKDIR/crom-src"
    go build -o "$CROM_BIN" ./cmd/crompressor/ 2>/dev/null
    echo "✅ Build OK"
else
    CROM_BIN=$(which crompressor)
    echo "✅ Crompressor já instalado: $CROM_BIN"
fi

# 2. Treinar codebook no arquivo
echo ""
echo "⏳ Treinando codebook..."
TRAIN_DIR="$WORKDIR/train"
mkdir -p "$TRAIN_DIR"
cp "$INPUT" "$TRAIN_DIR/"
"$CROM_BIN" train -i "$TRAIN_DIR" -o "$WORKDIR/bench.cromdb" -s 8192 2>&1 | grep -E '✔|Elite'

# 3. Benchmark
echo ""
echo "━━━ Resultados ━━━"
echo ""

# GZIP
T=$(($(date +%s%N)/1000000))
gzip -c -9 "$INPUT" > "$WORKDIR/out.gz"
T2=$(($(date +%s%N)/1000000))
GZ=$(stat -c%s "$WORKDIR/out.gz")
GZ_PCT=$(python3 -c "print(f'{$GZ/$ORIG_SIZE*100:.1f}')" 2>/dev/null || echo "?")
echo "  GZIP-9:  $(numfmt --to=iec $GZ) (${GZ_PCT}%) | $((T2-T))ms"

# ZSTD
if command -v zstd &>/dev/null; then
    T=$(($(date +%s%N)/1000000))
    zstd -3 -q "$INPUT" -o "$WORKDIR/out.zst" --force
    T2=$(($(date +%s%N)/1000000))
    ZST=$(stat -c%s "$WORKDIR/out.zst")
    ZST_PCT=$(python3 -c "print(f'{$ZST/$ORIG_SIZE*100:.1f}')" 2>/dev/null || echo "?")
    echo "  ZSTD-3:  $(numfmt --to=iec $ZST) (${ZST_PCT}%) | $((T2-T))ms"
else
    echo "  ZSTD:    (não instalado — apt install zstd)"
fi

# CROM
T=$(($(date +%s%N)/1000000))
PACK_OUT=$("$CROM_BIN" pack -i "$INPUT" -o "$WORKDIR/out.crom" -c "$WORKDIR/bench.cromdb" --mode vault 2>&1)
T2=$(($(date +%s%N)/1000000))
CROM_SZ=$(stat -c%s "$WORKDIR/out.crom")
CROM_PCT=$(python3 -c "print(f'{$CROM_SZ/$ORIG_SIZE*100:.1f}')" 2>/dev/null || echo "?")
HIT=$(echo "$PACK_OUT" | grep "Hit Rate" | head -1 | xargs)
ENT=$(echo "$PACK_OUT" | grep "Entropy" | head -1 | xargs)
echo "  CROM:    $(numfmt --to=iec $CROM_SZ) (${CROM_PCT}%) | $((T2-T))ms"
echo "  $HIT"
echo "  $ENT"

# Verificação lossless
echo ""
echo "━━━ Verificação Lossless ━━━"
"$CROM_BIN" unpack -i "$WORKDIR/out.crom" -o "$WORKDIR/restored" -c "$WORKDIR/bench.cromdb" 2>&1 | grep '✔' || true
if diff -q "$INPUT" "$WORKDIR/restored" > /dev/null 2>&1; then
    echo "  ✅ ROUNDTRIP LOSSLESS — arquivo restaurado é idêntico ao original"
else
    echo "  ❌ FALHA — arquivos diferentes!"
fi

# Limpeza
rm -rf "$WORKDIR"

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Compartilhe seus resultados:"
echo "  github.com/MrJc01/crompressor/issues"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
