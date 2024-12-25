package password

import "github.com/charmbracelet/bubbles/key"

type bindings struct {
	cancel key.Binding
	toggle key.Binding
	accept key.Binding
	quit   key.Binding
}

func newBindings() bindings {
	return bindings{
		cancel: key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "reset input")),
		toggle: key.NewBinding(key.WithKeys("ctrl+t"), key.WithHelp("ctrl+t", "show/hide password")),
		accept: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "confirm")),
		quit:   key.NewBinding(key.WithKeys("ctrl+c")),
	}
}

func (b bindings) ShortHelp() []key.Binding {
	return []key.Binding{b.cancel, b.toggle, b.accept}
}

func (b bindings) FullHelp() [][]key.Binding {
	return [][]key.Binding{b.ShortHelp()}
}
