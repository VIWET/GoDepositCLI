package tui

import tea "github.com/charmbracelet/bubbletea"

type Model struct {
	model tea.Model
}

func (m Model) Init() tea.Cmd {
	return m.model.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	model, cmd := m.model.Update(msg)
	m.model = model
	return m, cmd
}

func (m Model) View() string {
	return appContainerStyle.Render(m.model.View())
}

func newModel(model tea.Model) tea.Model {
	return Model{model}
}

func Run(model tea.Model) error {
	if _, err := tea.NewProgram(newModel(model)).Run(); err != nil {
		return err
	}

	return nil
}
