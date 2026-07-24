# LoreLang

Este repositorio consta de dos partes principales:

1. **Compilador / Parser**: Desarrollado en Go (`/compiler`), se encarga de analizar el código fuente en LoreLang y compilarlo a un archivo ejecutable en Ruby (`.rb`).
2. **Servidor Runtime**: Desarrollado en Ruby (`/runtime`), ejecuta el servidor web y la lógica para el agente interactivo.

---

## 🛠️ Requisitos y Dependencias

Asegúrate de contar con los siguientes entornos instalados:

- **[Go](https://go.dev/)** (v1.22+ recomendado)
- **[Ruby](https://www.ruby-lang.org/)** (v3.0+ recomendado)
- **Gemas de Ruby necesarias**:

  ```bash
  gem install sinatra puma rackup
  ```

- **[Ollama](https://ollama.com)**
- **[Gemma3 - via Ollama](https://ollama.com/library/gemma3)** (4b parametros recomendado)
  - Asegurate de correr el llm para probar el chatbot
  ```bash
  ollama run gemma3
  ```

---

## 🚀 Guía de Uso

### 1. Compilar el Parser (Go)

Desde la **raíz del proyecto**, ejecuta:

```bash
go build -C compiler -o ../build/lorelang ./cmd/lorelang
```

Esto generará el ejecutable `build/lorelang`.

---

### 2. Generar el código Ruby

El compilador acepta tu archivo `.lore` y produce un archivo Ruby (`.rb`). Tienes dos formas de usarlo:

Todos los comandos a continuación se ejecutan desde la **raíz del proyecto**.

#### Opción A — Output por defecto (junto al archivo fuente)

```bash
./build/lorelang examples/tom_nook.lore
```

El archivo `.rb` se generará en el mismo directorio que el `.lore`.

#### Opción B — Especificar la ruta de salida con `-o`

Usa el flag `-o` para indicar exactamente dónde quieres que se genere el archivo Ruby:

```bash
./build/lorelang -o ruta/de/salida/agente.rb examples/tu_archivo.lore
```

Por ejemplo, para generar el archivo directamente en `runtime/`:

```bash
./build/lorelang -o runtime/tom_nook_agent.rb examples/tom_nook.lore
```

> **✅ Recomendado**: Usa `-o` apuntando a `runtime/` para evitar tener que mover el archivo manualmente.

#### Ayuda del compilador

```bash
./build/lorelang -h
```

```
Uso: lorelang [opciones] <archivo.lore>

Opciones:
  -o string
        Ruta del archivo Ruby generado (por defecto: junto al archivo fuente)
```

---

### 3. Ejecutar el Servidor (Ruby)

Desde la raíz del proyecto, navega a `runtime/` y ejecuta el servidor:

```bash
cd runtime
ruby server.rb
```

El servidor quedará en ejecución escuchando las peticiones locales.

### 4. Abrir `./runtime/front/index.html` en tu navegador
