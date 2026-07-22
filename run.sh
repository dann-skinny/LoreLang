#!/bin/bash

if [ -z "$1" ]; then
  echo "Uso: ./run.sh examples/tom_nook.lore"
  exit 1
fi

echo "Limpiando agente anterior en el runtime..."
rm -f runtime/agent.rb

echo "Compilando $1 con LoreLang..."
cd compiler
go run cmd/lorelang/main.go -o ../runtime/agent.rb ../$1
cd ..

echo "Levantando el servidor con el agente compilado..."
cd runtime
ruby server.rb