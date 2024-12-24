package mnemonic

import "github.com/charmbracelet/bubbles/key"

type bindings struct {
	toggle key.Binding
	accept key.Binding
	quit   key.Binding
}

func newBindings() bindings {
	return bindings{
		toggle: key.NewBinding(key.WithKeys(" "), key.WithHelp("space", "show/hide mnemonic")),
		accept: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "accept")),
		quit:   key.NewBinding(key.WithKeys("q", "esc", "ctrl+c"), key.WithHelp("q", "quit")),
	}
}

func (b bindings) ShortHelp() []key.Binding {
	return []key.Binding{b.toggle, b.accept, b.quit}
}

func (b bindings) FullHelp() [][]key.Binding {
	return [][]key.Binding{b.ShortHelp()}
}
