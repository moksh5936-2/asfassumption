#!/bin/bash
set -euo pipefail

# Parity Runner — compares Go native engine vs Python ASF CLI outputs
# Usage: ./scripts/run-parity.sh [--regenerate] [samples-dir]

SAMPLE_DIR="${2:-asf-tui/testdata/parity/samples}"
PYTHON="${PYTHON:-.venv/bin/python3}"
NATIVE="${NATIVE:-/tmp/native-asf}"
OUTDIR="${OUTDIR:-asf-tui/testdata/parity}"

# Supported evidence files for policy analysis
EVIDENCE_FILES=(
    "${SAMPLE_DIR}/payroll_acl.csv"
    "${SAMPLE_DIR}/mfa_status.csv"
    "${SAMPLE_DIR}/iam_export.json"
    "${SAMPLE_DIR}/network_exposure.csv"
)

build_native() {
    echo ">>> Building Go native binary..."
    export PATH="/tmp/go/bin:$PATH"
    cd asf-tui && go build -o "$NATIVE" .
}

run_python() {
    local file="$1"
    local outfile="$2"
    local extra="${3:-}"
    export PYTHONPATH="${PYTHONPATH:-$(pwd)}"
    
    local cmd=("$PYTHON" -m asf.cli.main analyze --json)
    if [ "$extra" = "with-evidence" ]; then
        for ev in "${EVIDENCE_FILES[@]}"; do
            cmd+=(-e "$ev")
        done
    fi
    cmd+=("$file")
    
    "${cmd[@]}" 2>/dev/null > "$outfile"
    echo "  Python -> $(basename "$outfile") ($(wc -c < "$outfile") bytes)"
}

run_go() {
    local file="$1"
    local outfile="$2"
    local extra="${3:-}"
    
    local cmd=("$NATIVE" analyze)
    if [ "$extra" = "with-evidence" ]; then
        for ev in "${EVIDENCE_FILES[@]}"; do
            cmd+=(-e "$ev")
        done
    fi
    cmd+=("$file")
    
    "${cmd[@]}" 2>/dev/null > "$outfile"
    echo "  Go       -> $(basename "$outfile") ($(wc -c < "$outfile") bytes)"
}

compare_field_level() {
    local go_out="$1"
    local py_out="$2"
    local label="$3"
    
    python3 -c "
import json, sys

with open('$py_out') as f:
    py = json.load(f)
with open('$go_out') as f:
    go = json.load(f)

failures = 0
checks = 0

def safe_list(d, key):
    v = d.get(key)
    return v if v is not None else []

# 1. Summary
for key in ['claims_found','assumptions','verified','partially_verified','contradicted','unknown','critical_gaps']:
    if key in py.get('summary', {}) and key in go.get('summary', {}):
        checks += 1
        if py['summary'][key] != go['summary'][key]:
            failures += 1
            print(f'  FAIL: summary.{key} PY={py[\"summary\"][key]} GO={go[\"summary\"][key]}')

# 2. Verifications (sorted)
py_ver = sorted([v['result'] for v in safe_list(py, 'verifications')])
go_ver = sorted([v.get('result', '') for v in safe_list(go, 'verifications')])
checks += 1
if py_ver != go_ver:
    failures += 1
    print(f'  FAIL: verifications mismatch (PY={len(py_ver)} GO={len(go_ver)})')

# 3. Assumptions (sorted by type + whitespace-normalized text)
py_asms = sorted([(a['assumption_type'], ' '.join(a['text'].split())) for a in safe_list(py, 'assumptions')])
go_asms = sorted([(a['assumption_type'], ' '.join(a['text'].split())) for a in safe_list(go, 'assumptions')])
checks += 1
if py_asms != go_asms:
    failures += 1
    print(f'  FAIL: assumptions mismatch')
    for i, (p, g) in enumerate(zip(py_asms, go_asms)):
        if p != g:
            print(f'    #{i}: PY={p} vs GO={g}')

# 4. Gaps (sorted)
py_gaps = sorted([(g['type'], g['severity']) for g in safe_list(py, 'gaps')])
go_gaps = sorted([(g['type'], g['severity']) for g in safe_list(go, 'gaps')])
checks += 1
if py_gaps != go_gaps:
    failures += 1
    print(f'  FAIL: gaps mismatch')

print(f'  Checks: {checks}, Failures: {failures}')
sys.exit(0 if failures == 0 else 1)
" 3>&1 1>&2 2>&3 | tail -1
}

# Main
echo "=== ASF Parity Runner ==="
echo ""

if [ "${1:-}" = "--regenerate" ]; then
    build_native
    echo ""
    
    # Run TXT samples (no evidence)
    echo ">>> TXT samples (no evidence)..."
    for f in "$SAMPLE_DIR"/*.txt; do
        base=$(basename "$f" .txt)
        run_python "$f" "$OUTDIR/python/${base}.json" "no-evidence"
        run_go "$f" "$OUTDIR/go/${base}.json" "no-evidence"
    done
    
    # Run TXT finance_policy with evidence
    echo ""
    echo ">>> TXT finance_policy with evidence..."
    run_python "$SAMPLE_DIR/finance_policy.txt" "$OUTDIR/python/finance_policy.json" "with-evidence"
    run_go "$SAMPLE_DIR/finance_policy.txt" "$OUTDIR/go/finance_policy.json" "with-evidence"
    
    # Run PDF finance_policy with evidence
    echo ""
    echo ">>> PDF finance_policy with evidence..."
    run_python "$SAMPLE_DIR/finance_policy.pdf" "$OUTDIR/python/finance_policy.pdf.json" "with-evidence"
    run_go "$SAMPLE_DIR/finance_policy.pdf" "$OUTDIR/go/finance_policy.pdf.json" "with-evidence"
fi

# Compare
echo ""
echo "=== Field-Level Comparison ==="
failures=0
for f in "$OUTDIR"/go/*.json; do
    fname=$(basename "$f")
    pyfile="$OUTDIR/python/$fname"
    [ -f "$pyfile" ] || continue
    echo "--- $fname ---"
    compare_field_level "$f" "$pyfile" "$fname" || failures=$((failures + 1))
done

echo ""
if [ "$failures" -eq 0 ]; then
    echo "ALL SAMPLES: FULL PARITY CERTIFIED"
else
    echo "$failures SAMPLES HAVE FAILURES"
    exit 1
fi
