#!/bin/bash
# Add copyright header to all .go files
HEADER="// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz\n\n"

find . -name "*.go" -not -path "*/.git/*" | while read f; do
  if ! grep -q "Copyright (c) 2026 erik" "$f"; then
    content=$(cat "$f")
    printf "%b" "$HEADER" > "$f"
    echo "$content" >> "$f"
  fi
done
echo "Done."
