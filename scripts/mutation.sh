#!/bin/sh

set -eu

repo_dir=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
cd "$repo_dir"

workers=${MUTATION_WORKERS:-2}
test_cpu=${MUTATION_TEST_CPU:-0}
timeout_coefficient=${MUTATION_TIMEOUT_COEFFICIENT:-20}
report=target/quality/mutation.json

mkdir -p "$(dirname "$report")"

if [ "$test_cpu" -gt 0 ]; then
	set -- --test-cpu "$test_cpu"
else
	set --
fi

go run github.com/go-gremlins/gremlins/cmd/gremlins@v0.6.0 unleash . \
	--integration \
	--workers "$workers" \
	--timeout-coefficient "$timeout_coefficient" \
	--output "$report" \
	"$@"

echo "Mutation report: $report"
