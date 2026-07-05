# LoreLang

## Probar el parser

Desde la raiz del proyecto:

```bash
bash ./compile.sh
```

Ese script:

1. Ejecuta `lex` sobre `lorelang.lex`.
2. Ejecuta `yacc -d` sobre `lorelang.y`.
3. Compila el binario `lorelang`.
4. Ejecuta `./lorelang < tom_nook.lore`.

Tambien puedes correrlo manualmente:

```bash
lex lorelang.lex
yacc -d lorelang.y
cc y.tab.c lex.yy.c -lfl -o lorelang
./lorelang < tom_nook.lore
```

Para probar otro archivo:

```bash
./lorelang < mi_prueba.lore
```

## Limpiar archivos generados

Para eliminar artefactos generados por la compilacion:

```bash
bash ./compile.sh --clean
```

Tambien disponible en forma corta:

```bash
bash ./compile.sh -c
```
