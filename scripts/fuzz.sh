#!/bin/sh

set -eu

repo_dir=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
cd "$repo_dir"

fuzz_time=${FUZZ_TIME:-3s}

go test -list '^Fuzz' | while IFS= read -r target; do
	case "$target" in
	Fuzz*)
		echo "==> $target ($fuzz_time)"
		go test -run '^$' -fuzz "^${target}$" -fuzztime "$fuzz_time"
		;;
	esac
done
