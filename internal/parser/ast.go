package parser

// Program es la raíz del AST.
type Program struct {
	Character *Character `@@ EOF`
}

// Character representa la definición principal del personaje.
type Character struct {
	Name         string      `"Personaje" @String "{"`
	Attributes   []Attribute `"Atributos" "{" @@* "}"`
	InitialState string      `"EstadoInicial" ":" @Ident ";"`
	States       []State     `@@+ "}"`
}

// Attribute almacena propiedades estáticas del personaje.
type Attribute struct {
	Key   string `@Ident ":"`
	Value string `@String ";"`
}

// State representa un estado de la FSM del NPC.
type State struct {
	Name        string       `"Estado" @Ident "{"`
	Transitions []Transition `@@* "}"`
}

// Transition representa un trigger AlRecibir y su salto de estado.
type Transition struct {
	Trigger string `"AlRecibir" "(" @String ")" "->"`
	ToState string `@Ident ";"`
}
