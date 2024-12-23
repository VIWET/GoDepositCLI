package menu

import "github.com/charmbracelet/bubbles/key"

type bindings struct {
	up     key.Binding
	down   key.Binding
	accept key.Binding
	quit   key.Binding
}

func newBindings() bindings {
	return bindings{
		up:     key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "prev")),
		down:   key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "next")),
		accept: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "accept")),
		quit:   key.NewBinding(key.WithKeys("q", "esc", "ctrl+c"), key.WithHelp("q", "quit")),
	}
}

func (b bindings) ShortHelp() []key.Binding {
	return []key.Binding{b.up, b.down, b.accept, b.quit}
}

func (b bindings) FullHelp() [][]key.Binding {
	return [][]key.Binding{b.ShortHelp()}
}
