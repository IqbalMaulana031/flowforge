#!/usr/bin/env bash

set -euo pipefail

while IFS= read -r -d '' f; do
  awk '/^import \($/,/^\)$/{if($0=="")next}{print}' "$f" > /tmp/file
  mv /tmp/file "$f"
done < <(find . -name '*.go' -print0)

go list -f '{{.Dir}}' ./... | while IFS= read -r dir; do
  goimports -w -local flowforge-api "$dir"
done
gofmt -s -w .
