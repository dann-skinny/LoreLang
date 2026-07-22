package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"lorelang/internal/generator"
	"lorelang/internal/parser"
)

func main() {
	var outputPath string
	flag.StringVar(&outputPath, "o", "", "Ruta del archivo Ruby generado (por defecto: junto al archivo fuente)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Uso: %s [opciones] <archivo.lore>\n\nOpciones:\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}

	file := flag.Arg(0)
	src, err := os.ReadFile(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "No se pudo leer el archivo %q: %v\n", file, err)
		os.Exit(1)
	}

	ast, err := parser.Parse(file, src)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error sintactico: %v\n", err)
		os.Exit(1)
	}

	// Si no se especifica -o, el archivo .rb se genera junto al .lore
	if outputPath == "" {
		outputPath = filepath.Join(filepath.Dir(file), generator.OutputFileName(ast))
	}

	if err := generator.GenerateRuby(ast, outputPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error generando Ruby: %v\n", err)
		os.Exit(1)
	}

	char := ast.Character
	name, err := strconv.Unquote(char.Name)
	if err != nil {
		name = char.Name
	}

	fmt.Printf(
		"OK: personaje %q parseado correctamente (%d atributos, %d conocimientos, %d restricciones, estado inicial %q, %d estados). Ruby generado en %q.\n",
		name, len(char.Attributes), len(char.KnowledgeEntries), len(char.Restrictions), char.InitialState, len(char.States), outputPath,
	)
}
