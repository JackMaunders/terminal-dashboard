package main

import (
	"dashboard/widgets"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type model struct {
	clock  widgets.ClockWidget
	system widgets.SystemWidget
	height int
	width  int
}

func NewModel() model {
	return model{
		clock:  widgets.NewClockWidget(),
		system: widgets.NewSystemWidget(),
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.clock.Init(), m.system.Init())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case widgets.ClockTickMsg:
		updatedClock, clockCmd := m.clock.Update(msg)
		m.clock = updatedClock
		cmds = append(cmds, clockCmd)
	case widgets.SystemTickMsg:
		updatedSystem, systemCmd := m.system.Update(msg)
		m.system = updatedSystem
		cmds = append(cmds, systemCmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() tea.View {
	row1 := lipgloss.JoinHorizontal(
		lipgloss.Top, m.system.View(),
		lipgloss.PlaceHorizontal(m.width-lipgloss.Width(m.system.View()), lipgloss.Right, m.clock.View()),
	)

	content := lipgloss.JoinVertical(lipgloss.Left,
		row1,
		lipgloss.PlaceVertical(m.height-lipgloss.Height(row1), lipgloss.Bottom, lipgloss.NewStyle().Render("\nPress q to quit")),
	)

	v := tea.NewView(content)
	v.AltScreen = true
	return v
}
