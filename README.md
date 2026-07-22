# LoreLang

Este repositorio consta de dos partes principales:

1. **Compilador / Parser**: Desarrollado en Go (`/compiler`), se encarga de analizar el código fuente en LoreLang y compilarlo a un archivo ejecutable en Ruby (`.rb`).
2. **Servidor Runtime**: Desarrollado en Ruby (`/runtime`), ejecuta el servidor web y la lógica para el agente interactivo.

---

## 🛠️ Requisitos y Dependencias

Asegúrate de contar con los siguientes entornos instalados:

- **[Go](https://go.dev/)** (v1.18+ recomendado)
- **[Ruby](https://www.ruby-lang.org/)** (v3.0+ recomendado)
- **Gemas de Ruby necesarias**:

  ```bash
  gem install sinatra puma
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

Navega a la carpeta `compiler` y compila el binario ejecutable de LoreLang:

```bash
cd compiler
go build -o lorelang cmd/lorelang/main.go
```

Esto generará el ejecutable `lorelang` dentro del directorio `compiler/`.

---

### 2. Generar el código Ruby y moverlo a `runtime/`

Ejecuta el compilador pasándole tu archivo de código fuente en LoreLang (por ejemplo, `script.lore`):

```bash
./lorelang ruta/a/tu_archivo.lore
```

> **⚠️ Importante**: El compilador generará un archivo de código fuente en Ruby (`.rb`) en la ubicación donde ejecutes el comando. **Debes mover el archivo generado a la carpeta `runtime/`** antes de iniciar el servidor.

Ejemplo:

```bash
mv agente_generado.rb ../runtime/
```

---

### 3. Ejecutar el Servidor (Ruby)

Navega al directorio `runtime/` y ejecuta el servidor:

```bash
cd ../runtime
ruby server.rb
```

El servidor quedará en ejecución escuchando las peticiones locales.

### 4. Abrir './runtime/front/index.html' en tu navegador
