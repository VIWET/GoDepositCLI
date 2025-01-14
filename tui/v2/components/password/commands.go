package password

import tea "github.com/charmbracelet/bubbletea"

type confirm struct {
	validation   error
	confirmation error
}

func Confirm(password, confirmation string) tea.Cmd {
	return func() tea.Msg {
		return confirm{
			validatePassword(password),
			validateConfirmation(password, confirmation),
		}
	}
}
