#!/usr/bin/env bash
set -Eeuo pipefail

# Runs performance benchmarks entirely in Docker (same pipeline as build.sh).
# Usage: ./bench.sh <label>   e.g. ./bench.sh before | ./bench.sh after
#
# Outputs (per label):
#   .benchout/<label>/go-bench.txt        Go benchmarks (go test -bench -benchmem -count=6)
#   .benchout/<label>/frontend-bench.json frontend pipeline benchmark
#   .benchout/<label>/golden.json         rendered-block dump of synth-100k fixture
#
# If a previous label's golden.json exists, the new one is compared against it;
# a mismatch means rendering behaviour changed -> non-zero exit.

ROOT="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
cd "$ROOT"

LABEL="${1:?Usage: ./bench.sh <label> (e.g. before / after / stepN)}"
[[ "$LABEL" =~ ^[A-Za-z0-9_-]+$ ]] || { echo "Invalid label: $LABEL" >&2; exit 2; }

OUT="$ROOT/.benchout"
RUN_DIR="$OUT/$LABEL"
mkdir -p "$OUT"
rm -rf "$RUN_DIR"

export DOCKER_BUILDKIT=1
echo "mdmark bench: label=$LABEL (building in Docker with BuildKit...)"

docker build \
  --target benchout \
  --output "type=local,dest=$RUN_DIR" \
  .

echo
echo "=== Go benchmarks: $RUN_DIR/go-bench.txt ==="
cat "$RUN_DIR/go-bench.txt" || true
echo
echo "=== Frontend pipeline: $RUN_DIR/frontend-bench.json ==="
cat "$RUN_DIR/frontend-bench.json" || true
echo

# Behavioural equivalence: compare rendered blocks against the first recorded run.
BASE_GOLDEN="$(ls "$OUT"/*/golden.json 2>/dev/null | grep -v "^$RUN_DIR/golden.json$" | head -1 || true)"
if [[ -n "$BASE_GOLDEN" ]]; then
  if python3 - "$BASE_GOLDEN" "$RUN_DIR/golden.json" <<'PY'
import json, sys
a = json.load(open(sys.argv[1]))
b = json.load(open(sys.argv[2]))
sys.exit(0 if a == b else 1)
PY
  then
    echo "Golden check: MATCH (rendered output identical to $BASE_GOLDEN)"
  else
    echo "Golden check: MISMATCH vs $BASE_GOLDEN — rendering behaviour changed!" >&2
    exit 1
  fi
fi

echo
echo "Done. Results in $RUN_DIR"
