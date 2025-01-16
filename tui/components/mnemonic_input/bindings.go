package mnemonic_input

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
)

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
		accept:    key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "confirm"), key.WithDisabled()),
		quit:      key.NewBinding(key.WithKeys("ctrl+c")),
	}
}

func (b bindings) ShortHelp() []key.Binding {
	return []key.Binding{b.toggle, b.next, b.prev, b.accept}
}

func (b bindings) FullHelp() [][]key.Binding {
	return [][]key.Binding{b.ShortHelp()}
}

func inputBinding(bindings textinput.KeyMap) textinput.KeyMap {
	bindings.CharacterForward.SetEnabled(false)
	bindings.CharacterBackward.SetEnabled(false)
	bindings.WordForward.SetEnabled(false)
	bindings.WordBackward.SetEnabled(false)
	bindings.DeleteWordBackward.SetEnabled(true)
	bindings.DeleteWordForward.SetEnabled(true)
	bindings.DeleteAfterCursor.SetEnabled(true)
	bindings.DeleteBeforeCursor.SetEnabled(true)
	bindings.DeleteCharacterBackward.SetEnabled(true)
	bindings.DeleteCharacterForward.SetEnabled(true)
	bindings.LineStart.SetEnabled(false)
	bindings.LineEnd.SetEnabled(false)
	bindings.Paste.SetEnabled(false)
	bindings.AcceptSuggestion.SetEnabled(false)
	bindings.NextSuggestion.SetEnabled(false)
	bindings.PrevSuggestion.SetEnabled(false)

	return bindings
}
