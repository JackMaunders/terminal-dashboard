package main

import (
	"dashboard/widgets"
	"log"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type model struct {
	clock widgets.ClockWidget
}

type tickMsg time.Time

func NewModel() model {
	return model{
		clock: widgets.NewClockWidget(),
	}
}

func (m model) Init() tea.Cmd {
	return m.clock.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		log.Printf("Key pressed: %s", msg.String())
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}

	updatedClock, clockCmd := m.clock.Update(msg)
	m.clock = updatedClock

	return m, clockCmd
}

func (m model) View() tea.View {
	content := lipgloss.JoinVertical(lipgloss.Left,
		m.clock.View(),
		"\n  Press q to quit",
	)

	v := tea.NewView(content)
	v.AltScreen = true
	return v
}
