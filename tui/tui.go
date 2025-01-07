package tui

import tea "github.com/charmbracelet/bubbletea"

type Model struct {
	model tea.Model
	err   error
}

func (m *Model) Init() tea.Cmd {
	return m.model.Init()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case quit:
		m.err = msg.err
		// TODO: make more elegant solution
		return empty{}, tea.Quit
	default:
		model, cmd := m.model.Update(msg)
		m.model = model
		return m, cmd
	}
}

func (m *Model) View() string {
	return appContainerStyle.Render(m.model.View())
}

func newModel(model tea.Model) *Model {
	return &Model{model: model}
}

func Run(model tea.Model) error {
	appModel := newModel(model)
	if _, err := tea.NewProgram(appModel).Run(); err != nil {
		return err
	}

	if err := appModel.err; err != nil {
		return err
	}

	return nil
}

type empty struct{}

func (e empty) Init() tea.Cmd {
	return nil
}

func (e empty) Update(tea.Msg) (tea.Model, tea.Cmd) {
	return e, nil
}

func (e empty) View() string {
	return ""
}
