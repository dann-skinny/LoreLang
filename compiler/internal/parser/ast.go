package parser

// Program es la raíz del AST.
type Program struct {
	Character *Character `@@ EOF`
}

// Character representa la definición principal del personaje.
type Character struct {
	Name             string          `"Personaje" @String "{"`
	Attributes       []Attribute     `"Atributos" "{" @@* "}"`
	KnowledgeEntries []Entry         `"Conocimientos" "{" @@* "}"`
	Restrictions     []Entry         `"Restricciones" "{" @@* "}"`
	GlobalBehavior   *GlobalBehavior `"ComportamientoGlobal" "{" @@ "}"`
	InitialState     string          `"EstadoInicial" ":" @Ident ";"`
	States           []State         `@@+ "}"`
}

// Attribute almacena propiedades estáticas del personaje.
type Attribute struct {
	Key   string `@Ident ":"`
	Value string `@String ";"`
}

// Entry representa una propiedad con valor string o lista de strings.
type Entry struct {
	Key   string     `@Ident ":"`
	Value EntryValue `@@ ";"`
}

// EntryValue representa el tipo de valor permitido para una propiedad.
type EntryValue struct {
	String *string  `  @String`
	List   []string `| "[" @String ( "," @String )* "]"`
}

// GlobalBehavior agrupa reglas que aplican siempre.
type GlobalBehavior struct {
	Rules []GlobalRule `@@+`
}

// GlobalRule puede ser un trigger AlRecibir o fallback AlDesconocer.
type GlobalRule struct {
	ReceiveRule *ReceiveRule `  @@`
	UnknownRule *UnknownRule `| @@`
}

// ReceiveRule representa un bloque AlRecibir.
type ReceiveRule struct {
	Triggers []string `"AlRecibir" "(" @String ( "," @String )* ")" "{"`
	Actions  []Action `@@+ "}"`
}

// UnknownRule representa el fallback AlDesconocer.
type UnknownRule struct {
	Actions []Action `"AlDesconocer" "{" @@+ "}"`
}

// State representa un estado de la FSM del NPC.
type State struct {
	Name     string        `"Estado" @Ident "{"`
	Handlers []ReceiveRule `@@+ "}"`
}

// Action representa una instrucción dentro de un bloque AlRecibir/AlDesconocer.
type Action struct {
	Transition     *string `  "Transicion" ":" @Ident ";"`
	InjectContext  *string `| "InyectarContexto" ":" @String ";"`
	DirectResponse *string `| "ResponderDirecto" ":" @String ";"`
}
