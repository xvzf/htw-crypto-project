#!/bin/zsh -eu

for i in {1..50}; do
  echo "[+] Encrypting plain.txt ... iteration $i"
  go run ./cmd encrypt -k key.png plain.txt ./solver/ciphertexts/$i.json
done
