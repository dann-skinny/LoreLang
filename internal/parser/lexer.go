package parser

import "github.com/alecthomas/participle/v2/lexer"

var loreLexer = lexer.MustSimple([]lexer.SimpleRule{
	{Name: "Comment", Pattern: `//[^\n]*`},
	{Name: "Whitespace", Pattern: `[\s\t\r\n]+`},
	{Name: "Arrow", Pattern: `->`},
	{Name: "Keyword", Pattern: `\b(Personaje|Atributos|EstadoInicial|Estado|AlRecibir)\b`},
	{Name: "Ident", Pattern: `[a-zA-Z_][a-zA-Z0-9_]*`},
	{Name: "String", Pattern: `"(\\"|[^"])*"`},
	{Name: "Punct", Pattern: `[:{}();,]`},
})
