package mnemonicInput

import "github.com/charmbracelet/bubbles/key"

type bindings struct {
	toggle    key.Binding
	prev      key.Binding
	next      key.Binding
	space     key.Binding
	backspace key.Binding
	accept    key.Binding
	quit      key.Binding
}

func newBindings() bindings {
	return bindings{
		toggle:    key.NewBinding(key.WithKeys("ctrl+t"), key.WithHelp("ctrl+t", "show/hide password")),
		next:      key.NewBinding(key.WithKeys("ctrl+n"), key.WithHelp("ctrl+n", "next word")),
		prev:      key.NewBinding(key.WithKeys("ctrl+p"), key.WithHelp("ctrl+p", "prev word")),
		space:     key.NewBinding(key.WithKeys(" ")),
		backspace: key.NewBinding(key.WithKeys("backspace")),
		accept:    key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "confirm")),
		quit:      key.NewBinding(key.WithKeys("ctrl+c")),
	}
}

func (b bindings) ShortHelp() []key.Binding {
	return []key.Binding{b.toggle, b.accept, b.next, b.prev}
}

func (b bindings) FullHelp() [][]key.Binding {
	return [][]key.Binding{b.ShortHelp()}
}
