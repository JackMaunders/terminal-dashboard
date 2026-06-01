package main

import (
	"dashboard/widgets"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type model struct {
	clock   widgets.ClockWidget
	system  widgets.SystemWidget
	weather widgets.WeatherWidget
	height  int
	width   int
}

func NewModel() model {
	return model{
		clock:   widgets.NewClockWidget(),
		system:  widgets.NewSystemWidget(),
		weather: widgets.NewWeatherWidget(),
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.clock.Init(), m.system.Init(), m.weather.Init())
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
	case widgets.WeatherTickMsg:
		updatedWeather, weatherCmd := m.weather.Update(msg)
		m.weather = updatedWeather
		cmds = append(cmds, weatherCmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() tea.View {
	row1 := lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.system.View(),
		lipgloss.PlaceHorizontal(m.width-lipgloss.Width(m.system.View()), lipgloss.Right, m.weather.View()),
	)

	quitMessage := lipgloss.NewStyle().Render("\nPress q to quit")

	row2 := lipgloss.JoinHorizontal(
		lipgloss.Bottom,
		quitMessage,
		lipgloss.PlaceHorizontal(m.width-lipgloss.Width(quitMessage), lipgloss.Right, m.clock.View()),
	)

	content := lipgloss.JoinVertical(lipgloss.Left,
		row1,
		lipgloss.PlaceVertical(m.height-lipgloss.Height(row1), lipgloss.Bottom, row2),
	)

	v := tea.NewView(content)
	v.AltScreen = true
	return v
}
