package parser

import (
	"fmt"

	"github.com/alecthomas/participle/v2"
)

var loreParser = participle.MustBuild[Program](
	participle.Lexer(loreLexer),
	participle.Elide("Whitespace", "Comment"),
)

// Parse convierte contenido .lore en AST.
func Parse(filename string, src []byte) (*Program, error) {
	ast, err := loreParser.ParseBytes(filename, src)
	if err != nil {
		return nil, fmt.Errorf("error de parsing: %w", err)
	}
	return ast, nil
}
