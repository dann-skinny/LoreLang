# Copilot instructions for LoreLang

## Project Vision

- LoreLang is an Agent Orchestrator DSL.
- The ultimate goal is not just syntax checking, but generating structured output (like JSON blocks or configuration objects) that will be consumed by a Ruby/Python backend.
- The parser should be "context-aware" regarding the Persona: it will eventually map specific 'AlRecibir' event blocks to logic that manages API calls for an LLM (Ollama).
- When suggesting C code within `lorelang.y` semantic actions, keep in mind that these will eventually need to build a context tree or a dictionary of states, not just print to stdout.

## Build and run

- Build the parser:
  ```bash
  lex lorelang.lex
  yacc -d lorelang.y
  cc y.tab.c lex.yy.c -lfl -o lorelang
  ```
- Run the sample program / smoke test:
  ```bash
  ./lorelang < tom_nook.lore
  ```
- `compile.sh` wraps the same sequence.

## Architecture

- This repo is a small language frontend built with Lex/Flex and Yacc/Bison.
- `lorelang.lex` is the lexer: it maps Spanish keywords to tokens, skips whitespace and `//` comments, and returns `ID` and `STRING` tokens.
- `lorelang.y` is the grammar: it accepts one top-level `Personaje` block containing `Atributos`, `Restricciones`, zero or more `AlRecibir` events, and one `Funcion` block.
- Semantic actions only print validation messages; the parser is primarily a syntax checker.
- `tom_nook.lore` is the canonical end-to-end example for the language shape.

## Conventions

- Keywords are Spanish and case sensitive: `Personaje`, `Atributos`, `Restricciones`, `AlRecibir`, `InyectarContexto`, `Si`, `Sino`, `Responder`, `GenerarRespuesta`, `Funcion`.
- The grammar is strict about punctuation: braces, parentheses, brackets, commas, colons, and semicolons all matter.
- Strings must use double quotes.
- The parser expects `Personaje STRING { ... }` as the outer structure.
- Attribute declarations use `ID: STRING;`.
- Restriction declarations use either `ID: STRING;` or `ID: [STRING, STRING, ...];`.
- Event blocks use `AlRecibir("...") { InyectarContexto: "..."; }`.
- The function block expects `Funcion name(param) { ... }` with a `Si` branch and optional `Sino` branch.
- Keep lexer and grammar changes in sync; adding a token in one file usually requires updates in the other.
