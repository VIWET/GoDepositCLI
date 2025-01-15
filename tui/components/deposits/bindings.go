package deposits

import "github.com/charmbracelet/bubbles/key"

type bindings struct {
	quit key.Binding
}

func newBindings() bindings {
	return bindings{
		quit: key.NewBinding(key.WithKeys("ctrl+c")),
	}
}
