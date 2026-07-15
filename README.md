# LoreLang

## Frontend Go (Participle)

Este repositorio ahora incluye un frontend en Go para parsear archivos `.lore` y validar coherencia semántica básica de FSM (estado inicial, estados y transiciones `AlRecibir`).

```bash
go run ./cmd/lorelang ./tom_nook_fsm.lore
```

## Compilación y uso

Usa el script portable desde la raíz para compilar el parser (autodetecta las herramientas necesarias como `lex`/`flex` y `yacc`/`bison` según el OS):

```bash
./compile.sh          # Compilar
./compile.sh --bison  # Forzar uso de Bison (Recomendado para macos)
./compile.sh -c --clean       # Limpiar archivos generados 
```

**Ejecutar script de prueba:**
```bash
./lorelang < tom_nook.lore
```

<details>
<summary>Comandos manuales para depuración</summary>

```bash
lex lorelang.lex                       # lexer
yacc -d lorelang.y                     # parser (o bison -y -d)
cc y.tab.c lex.yy.c -lfl -o lorelang   # enlazar (probar -ll si falla)
```
</details>
