package main

import (
	"dashboard/widgets"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type model struct {
	clock   widgets.ClockWidget
	pokemon widgets.PokemonWidget
	system  widgets.SystemWidget
	weather widgets.WeatherWidget
	height  int
	width   int
}

func NewModel() model {
	return model{
		clock:   widgets.NewClockWidget(),
		pokemon: widgets.NewPokemonWidget(),
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
	quitMessage := lipgloss.NewStyle().Render("\nPress q to quit")

	leftColumn := lipgloss.JoinVertical(lipgloss.Left, m.pokemon.View(), quitMessage)

	leftWidth := lipgloss.Width(leftColumn)
	rightWidth := m.width - leftWidth

	weatherBox := lipgloss.PlaceHorizontal(rightWidth, lipgloss.Right, m.weather.View())

	systemBox := lipgloss.PlaceHorizontal(rightWidth, lipgloss.Right, m.system.View())
	systemBox = lipgloss.PlaceVertical(
		m.height-lipgloss.Height(weatherBox)-lipgloss.Height(m.clock.View()),
		lipgloss.Bottom,
		systemBox,
	)

	clockBox := lipgloss.PlaceHorizontal(rightWidth, lipgloss.Right, m.clock.View())
	clockBox = lipgloss.PlaceVertical(
		m.height-lipgloss.Height(weatherBox)-lipgloss.Height(systemBox),
		lipgloss.Bottom,
		clockBox,
	)

	rightColumn := lipgloss.JoinVertical(lipgloss.Top, weatherBox, systemBox, clockBox)

	content := lipgloss.JoinHorizontal(lipgloss.Left,
		leftColumn,
		rightColumn,
	)

	v := tea.NewView(content)
	v.AltScreen = true
	return v
}
