package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/urfave/cli/v2"
	"github.com/viwet/GoDepositCLI/app"
)

type NewModel[Config app.ConfigConstraint] func(ctx *cli.Context, state *app.State[Config]) (tea.Model, tea.Cmd)

type MainModel struct {
	model       tea.Model
	initCommand tea.Cmd
	err         error
}

func (m MainModel) Err() error {
	return m.err
}

func newMainModel[Config app.ConfigConstraint](model tea.Model, initCommand tea.Cmd) *MainModel {
	return &MainModel{
		model:       model,
		initCommand: initCommand,
	}
}

func Run[Config app.ConfigConstraint](ctx *cli.Context, state *app.State[Config], newModel NewModel[Config]) error {
	var (
		model = newMainModel[Config](newModel(ctx, state))
		err   error
	)

	if _, err = tea.NewProgram(model).Run(); err != nil {
		return err
	}

	if err := model.Err(); err != nil {
		return err
	}

	return nil
}

func (m *MainModel) Init() tea.Cmd {
	return m.initCommand
}

func (m *MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case quit:
		m.err = msg.Err()
		return m, tea.Quit
	default:
		model, cmd := m.model.Update(msg)
		m.model = model
		return m, cmd
	}
}

var container = lipgloss.NewStyle().Padding(1, 2)

func (m *MainModel) View() string {
	if m.model != nil {
		return container.Render(m.model.View())
	}
	return ""
}
