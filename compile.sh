#! /usr/bin/env bash

set -euo pipefail

clean() {
  rm -f lex.yy.c y.tab.c y.tab.h lorelang
  echo "Archivos generados eliminados."
}

use_bison=false

for arg in "$@"; do
  case "$arg" in
    --clean|-c)
      clean
      exit 0
      ;;
    --bison)
      use_bison=true
      ;;
    *)
      echo "Uso: $0 [--clean|-c] [--bison]"
      exit 1
      ;;
  esac
done

if command -v lex >/dev/null 2>&1; then
  lex lorelang.lex
elif command -v flex >/dev/null 2>&1; then
  flex lorelang.lex
else
  echo "No se encontro 'lex' ni 'flex'."
  exit 1
fi

if [[ "$use_bison" == true || "$(uname -s)" == "Darwin" ]]; then
  if command -v bison >/dev/null 2>&1; then
    bison -y -d lorelang.y
  else
    yacc -d lorelang.y
  fi
else
  yacc -d lorelang.y
fi

if cc y.tab.c lex.yy.c -lfl -o lorelang >/dev/null 2>&1; then
  :
elif cc y.tab.c lex.yy.c -ll -o lorelang >/dev/null 2>&1; then
  :
else
  echo "No se pudo enlazar con la libreria de lex (-lfl o -ll)."
  exit 1
fi

./lorelang < tom_nook.lore
