package mnemonic

import "github.com/charmbracelet/bubbles/key"

type bindings struct {
	toggle key.Binding
	accept key.Binding
	quit   key.Binding
}

func newBindings() bindings {
	return bindings{
		toggle: key.NewBinding(key.WithKeys("ctrl+t"), key.WithHelp("ctrl+t", "show/hide mnemonic")),
		accept: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "accept")),
		quit:   key.NewBinding(key.WithKeys("ctrl+c")),
	}
}

func (b bindings) ShortHelp() []key.Binding {
	return []key.Binding{b.toggle, b.accept}
}

func (b bindings) FullHelp() [][]key.Binding {
	return [][]key.Binding{b.ShortHelp()}
}
