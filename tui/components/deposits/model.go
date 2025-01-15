package deposits

import (
	"context"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/urfave/cli/v2"
	"github.com/viwet/GoDepositCLI/app"
	"github.com/viwet/GoDepositCLI/tui"
)

type Model struct {
	ctx    *cli.Context
	cancel context.CancelFunc

	ticks     <-chan deposit
	increment float64
	quitting  bool
	dir       string

	bindings bindings
	style    style
	progress progress.Model
}

func New(ctx *cli.Context, state *app.State[app.DepositConfig]) (tea.Model, tea.Cmd) {
	engineCtx, engineCancel := context.WithCancel(ctx.Context)
	result, ticks := RunEngine(engineCtx, state)
	return &Model{
		ctx:       ctx,
		cancel:    engineCancel,
		ticks:     ticks,
		increment: 1.0 / float64(state.Config().Number),
		dir:       state.Config().Directory,

		bindings: newBindings(),
		style:    newStyle(tui.DefaultColorscheme()),
		progress: progress.New(progress.WithSolidFill(tui.DefaultColorscheme().Magenta.Light)),
	}, tea.Batch(result, WaitDeposit(ticks))
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case result:
		if msg.err != nil {
			return m, tui.QuitWithError(msg.err)
		}
		m.quitting = true
		return m, SaveResult(msg, m.dir)
	case deposit:
		return m, tea.Batch(m.progress.IncrPercent(m.increment), WaitDeposit(m.ticks))
	case tea.KeyMsg:
		if key.Matches(msg, m.bindings.quit) {
			m.cancel()
		}

		return m, nil
	default:
		prog, cmd := m.progress.Update(msg)
		m.progress = prog.(progress.Model)
		return m, cmd
	}
}

func (m *Model) View() string {
	m.progress.PercentageStyle.Foreground(m.style.colors.Black).Bold(true)
	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.style.title.Foreground(m.style.colors.White).Render("Deposits and Keystores"),
		m.style.container.Render(m.progressView()),
	)
}

func (m *Model) progressView() string {
	if m.quitting {
		return m.progress.ViewAs(1)
	}
	return m.progress.View()
}
