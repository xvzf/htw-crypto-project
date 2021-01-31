#!/bin/zsh -eu

wget https://maxhalford.github.io/img/blog/halftoning-1/gray.png

for i in {1..50}; do
  echo "[+] Encrypting ../assets/plain.txt ... iteration $i"
  go run ../../cmd encrypt -k gray.png ../assets/plain.txt ./$i.json
done
