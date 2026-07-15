package parser

import (
	"fmt"
	"strconv"
)

// Validate aplica reglas semánticas básicas de coherencia FSM.
func Validate(program *Program) error {
	char := program.Character
	if char == nil {
		return fmt.Errorf("no se encontro bloque Personaje")
	}

	stateSet := map[string]struct{}{}
	for _, st := range char.States {
		if _, exists := stateSet[st.Name]; exists {
			return fmt.Errorf("estado duplicado: %q", st.Name)
		}
		stateSet[st.Name] = struct{}{}

		triggerSet := map[string]struct{}{}
		for _, tr := range st.Transitions {
			trigger, err := unquote(tr.Trigger)
			if err != nil {
				return fmt.Errorf("trigger invalido en estado %q: %w", st.Name, err)
			}
			if _, duplicated := triggerSet[trigger]; duplicated {
				return fmt.Errorf("trigger AlRecibir duplicado en estado %q: %q", st.Name, trigger)
			}
			triggerSet[trigger] = struct{}{}
		}
	}

	if _, ok := stateSet[char.InitialState]; !ok {
		return fmt.Errorf("EstadoInicial %q no existe en los bloques Estado", char.InitialState)
	}

	for _, st := range char.States {
		for _, tr := range st.Transitions {
			if _, ok := stateSet[tr.ToState]; !ok {
				trigger, err := unquote(tr.Trigger)
				if err != nil {
					return fmt.Errorf("trigger invalido en estado %q: %w", st.Name, err)
				}
				return fmt.Errorf(
					"transicion invalida en estado %q con trigger %q: destino %q no existe",
					st.Name, trigger, tr.ToState,
				)
			}
		}
	}

	attrSet := map[string]struct{}{}
	for _, attr := range char.Attributes {
		if _, exists := attrSet[attr.Key]; exists {
			return fmt.Errorf("atributo duplicado: %q", attr.Key)
		}
		attrSet[attr.Key] = struct{}{}
		if _, err := unquote(attr.Value); err != nil {
			return fmt.Errorf("valor invalido para atributo %q: %w", attr.Key, err)
		}
	}

	if _, err := unquote(char.Name); err != nil {
		return fmt.Errorf("nombre de personaje invalido: %w", err)
	}

	return nil
}

func unquote(raw string) (string, error) {
	return strconv.Unquote(raw)
}
