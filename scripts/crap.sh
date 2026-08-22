#!/bin/sh

set -eu

repo_dir=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
cd "$repo_dir"

workers=${CRAP_WORKERS:-2}
report=target/quality/crap.txt

mkdir -p "$(dirname "$report")"

go run github.com/unclebob/crap4go/cmd/crap4go@v0.0.0-20260521190544-bee16dbdadb4 \
	--max-workers "$workers" >"$report"

cat "$report"
echo "CRAP report: $report"
