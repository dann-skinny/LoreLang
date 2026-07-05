#! /usr/bin/env bash

set -euo pipefail

clean() {
  rm -f lex.yy.c y.tab.c y.tab.h lorelang
  echo "Archivos generados eliminados."
}

if [[ "${1:-}" == "--clean" || "${1:-}" == "-c" ]]; then
  clean
  exit 0
fi

lex lorelang.lex
yacc -d lorelang.y
cc y.tab.c lex.yy.c -lfl -o lorelang
./lorelang < tom_nook.lore
