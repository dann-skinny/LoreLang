package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"unicode"

	"lorelang/internal/parser"
)

type rubyTemplateData struct {
	ClassName          string
	InitialState       string
	Attributes         []rubyHashEntry
	KnowledgeEntries   []rubyHashEntry
	Restrictions       []rubyHashEntry
	GlobalRules        []rubyRule
	States             []rubyState
	UnknownFallback    rubyOutcome
	HasUnknownFallback bool
}

type rubyHashEntry struct {
	Key   string
	Value string
}

type rubyState struct {
	Name  string
	Rules []rubyRule
}

type rubyRule struct {
	Triggers []string
	Outcome  rubyOutcome
}

type rubyOutcome struct {
	Action      string
	Context     string
	NewState    string
	HasNewState bool
}

var rubyAgentTemplate = template.Must(template.New("ruby_agent").Parse(`class {{.ClassName}}
  ATRIBUTOS = {
{{- range .Attributes}}
    {{.Key}}: {{.Value}},
{{- end}}
  }.freeze

  CONOCIMIENTOS = {
{{- range .KnowledgeEntries}}
    {{.Key}}: {{.Value}},
{{- end}}
  }.freeze

  RESTRICCIONES = {
{{- range .Restrictions}}
    {{.Key}}: {{.Value}},
{{- end}}
  }.freeze

  def initialize
    @estado_actual = {{.InitialState}}
  end

  def evaluar_mensaje(input)
    texto = input.to_s.downcase

    case @estado_actual
{{- range .States}}
    when {{.Name}}
{{- if .Rules}}
{{- range $i, $rule := .Rules}}
      {{if eq $i 0}}if{{else}}elsif{{end}} {{range $j, $t := $rule.Triggers}}{{if gt $j 0}} || {{end}}texto.include?({{$t}}){{end}}
{{- if $rule.Outcome.HasNewState}}
        @estado_actual = {{$rule.Outcome.NewState}}
{{- end}}
        return { accion: {{$rule.Outcome.Action}}, contexto: {{$rule.Outcome.Context}}, nuevo_estado: {{if $rule.Outcome.HasNewState}}@estado_actual{{else}}nil{{end}} }
{{- end}}
      end
{{- end}}
{{- end}}
    end

{{- if .GlobalRules}}
    # ComportamientoGlobal
{{- range $i, $rule := .GlobalRules}}
    {{if eq $i 0}}if{{else}}elsif{{end}} {{range $j, $t := $rule.Triggers}}{{if gt $j 0}} || {{end}}texto.include?({{$t}}){{end}}
{{- if $rule.Outcome.HasNewState}}
      @estado_actual = {{$rule.Outcome.NewState}}
{{- end}}
      return { accion: {{$rule.Outcome.Action}}, contexto: {{$rule.Outcome.Context}}, nuevo_estado: {{if $rule.Outcome.HasNewState}}@estado_actual{{else}}nil{{end}} }
{{- end}}
    end
{{- end}}

{{- if .HasUnknownFallback}}
{{- if .UnknownFallback.HasNewState}}
    @estado_actual = {{.UnknownFallback.NewState}}
{{- end}}
    return { accion: {{.UnknownFallback.Action}}, contexto: {{.UnknownFallback.Context}}, nuevo_estado: {{if .UnknownFallback.HasNewState}}@estado_actual{{else}}nil{{end}} }
{{- else}}
    return { accion: :llm, contexto: nil, nuevo_estado: nil }
{{- end}}
  end
end
`))

// GenerateRuby genera un agente Ruby a partir del AST parseado.
func GenerateRuby(ast *parser.Program, outputPath string) error {
	if ast == nil || ast.Character == nil {
		return fmt.Errorf("AST invalido: personaje no definido")
	}
	if outputPath == "" {
		return fmt.Errorf("ruta de salida vacia")
	}

	data, err := buildTemplateData(ast)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(outputPath), 0o755); err != nil {
		return fmt.Errorf("no se pudo crear directorio de salida: %w", err)
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("no se pudo crear archivo Ruby %q: %w", outputPath, err)
	}
	defer file.Close()

	if err := rubyAgentTemplate.Execute(file, data); err != nil {
		return fmt.Errorf("no se pudo renderizar plantilla Ruby: %w", err)
	}

	return nil
}

func buildTemplateData(ast *parser.Program) (rubyTemplateData, error) {
	char := ast.Character

	data := rubyTemplateData{
		ClassName:        "Agent",
		InitialState:     toRubySymbol(char.InitialState),
		Attributes:       make([]rubyHashEntry, 0, len(char.Attributes)),
		KnowledgeEntries: make([]rubyHashEntry, 0, len(char.KnowledgeEntries)),
		Restrictions:     make([]rubyHashEntry, 0, len(char.Restrictions)),
		GlobalRules:      []rubyRule{},
		States:           make([]rubyState, 0, len(char.States)),
	}

	for _, attr := range char.Attributes {
		data.Attributes = append(data.Attributes, rubyHashEntry{
			Key:   toSnakeCase(attr.Key),
			Value: rubyStringLiteral(unquoteOrRaw(attr.Value)),
		})
	}

	for _, entry := range char.KnowledgeEntries {
		data.KnowledgeEntries = append(data.KnowledgeEntries, rubyHashEntry{
			Key:   toSnakeCase(entry.Key),
			Value: entryValueToRuby(entry.Value),
		})
	}

	for _, entry := range char.Restrictions {
		data.Restrictions = append(data.Restrictions, rubyHashEntry{
			Key:   toSnakeCase(entry.Key),
			Value: entryValueToRuby(entry.Value),
		})
	}

	for _, rule := range char.GlobalBehavior.Rules {
		if rule.ReceiveRule != nil {
			var triggers []string
			for _, t := range rule.ReceiveRule.Triggers {
				triggers = append(triggers, rubyStringLiteral(strings.ToLower(unquoteOrRaw(t))))
			}
			data.GlobalRules = append(data.GlobalRules, rubyRule{
				Triggers: triggers,
				Outcome:  actionsToOutcome(rule.ReceiveRule.Actions),
			})
		}
		if rule.UnknownRule != nil {
			data.HasUnknownFallback = true
			data.UnknownFallback = actionsToOutcome(rule.UnknownRule.Actions)
		}
	}

	for _, st := range char.States {
		state := rubyState{
			Name:  toRubySymbol(st.Name),
			Rules: make([]rubyRule, 0, len(st.Handlers)),
		}
		for _, handler := range st.Handlers {
			var triggers []string
			for _, t := range handler.Triggers {
				triggers = append(triggers, rubyStringLiteral(strings.ToLower(unquoteOrRaw(t))))
			}
			state.Rules = append(state.Rules, rubyRule{
				Triggers: triggers,
				Outcome:  actionsToOutcome(handler.Actions),
			})
		}
		data.States = append(data.States, state)
	}

	return data, nil
}

func entryValueToRuby(value parser.EntryValue) string {
	if value.String != nil {
		return rubyStringLiteral(unquoteOrRaw(*value.String))
	}

	items := make([]string, 0, len(value.List))
	for _, item := range value.List {
		items = append(items, rubyStringLiteral(unquoteOrRaw(item)))
	}
	return "[" + strings.Join(items, ", ") + "]"
}

func actionsToOutcome(actions []parser.Action) rubyOutcome {
	outcome := rubyOutcome{
		Action:  ":llm",
		Context: "nil",
	}

	for _, action := range actions {
		if action.Transition != nil {
			outcome.HasNewState = true
			outcome.NewState = toRubySymbol(*action.Transition)
		}
		if action.InjectContext != nil {
			outcome.Action = ":llm"
			outcome.Context = rubyStringLiteral(unquoteOrRaw(*action.InjectContext))
		}
		if action.DirectResponse != nil {
			outcome.Action = ":responder_directo"
			outcome.Context = rubyStringLiteral(unquoteOrRaw(*action.DirectResponse))
		}
	}

	return outcome
}

func unquoteOrRaw(value string) string {
	unquoted, err := strconv.Unquote(value)
	if err != nil {
		return value
	}
	return unquoted
}

func rubyStringLiteral(value string) string {
	return strconv.Quote(value)
}

func toRubySymbol(value string) string {
	return ":" + toSnakeCase(value)
}

func toSnakeCase(value string) string {
	var b strings.Builder
	prevUnderscore := false
	prevLowerOrDigit := false

	for _, r := range value {
		switch {
		case r == '_' || r == '-' || unicode.IsSpace(r):
			if b.Len() > 0 && !prevUnderscore {
				b.WriteRune('_')
				prevUnderscore = true
			}
			prevLowerOrDigit = false
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			if unicode.IsUpper(r) && prevLowerOrDigit && !prevUnderscore {
				b.WriteRune('_')
			}
			lower := unicode.ToLower(r)
			b.WriteRune(lower)
			prevUnderscore = false
			prevLowerOrDigit = unicode.IsLower(lower) || unicode.IsDigit(lower)
		default:
			if b.Len() > 0 && !prevUnderscore {
				b.WriteRune('_')
				prevUnderscore = true
			}
			prevLowerOrDigit = false
		}
	}

	result := strings.Trim(b.String(), "_")
	if result == "" {
		return "valor"
	}
	return result
}

func toPascalCase(value string) string {
	parts := splitWords(value)
	if len(parts) == 0 {
		return "Character"
	}

	var b strings.Builder
	for _, part := range parts {
		lower := strings.ToLower(part)
		runes := []rune(lower)
		if len(runes) == 0 {
			continue
		}
		runes[0] = unicode.ToUpper(runes[0])
		b.WriteString(string(runes))
	}
	result := b.String()
	if result == "" {
		return "Character"
	}
	return result
}

func splitWords(value string) []string {
	var words []string
	var current []rune

	flush := func() {
		if len(current) > 0 {
			words = append(words, string(current))
			current = nil
		}
	}

	for _, r := range value {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			current = append(current, r)
		} else {
			flush()
		}
	}
	flush()

	return words
}

// OutputFileName retorna siempre "agent.rb" (contrato fijo con el servidor).
func OutputFileName(_ *parser.Program) string {
	return "agent.rb"
}
