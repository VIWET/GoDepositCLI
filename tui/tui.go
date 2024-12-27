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
	model, cmd := m.model.Update(msg)
	m.model = model
	return m, cmd
}

func (m *Model) View() string {
	return appContainerStyle.Render(m.model.View())
}

func newModel(model tea.Model) *Model {
	return &Model{model: model}
}

func (m *Model) filter(_ tea.Model, msg tea.Msg) tea.Msg {
	switch msg := msg.(type) {
	case quit:
		m.err = msg.err
		return tea.Quit()
	default:
		return msg
	}
}

func Run(model tea.Model) error {
	appModel := newModel(model)
	if _, err := tea.NewProgram(
		appModel,
		tea.WithAltScreen(),
		tea.WithFilter(appModel.filter),
	).Run(); err != nil {
		return err
	}

	if err := appModel.err; err != nil {
		return err
	}

	return nil
}
