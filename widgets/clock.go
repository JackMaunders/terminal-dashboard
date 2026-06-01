package widgets

import (
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type ClockTickMsg time.Time

// ClockWidget displays the current time and date
type ClockWidget struct {
	time time.Time
}

func NewClockWidget() ClockWidget {
	return ClockWidget{time: time.Now()}
}

// Init returns the first tick command
func (c ClockWidget) Init() tea.Cmd {
	return clockTick(time.Second)
}

// Update handles tick messages to refresh the time
func (c ClockWidget) Update(msg tea.Msg) (ClockWidget, tea.Cmd) {
	switch msg.(type) {
	case ClockTickMsg:
		c.time = time.Now()
		return c, clockTick(time.Second)
	}
	return c, nil
}

// View renders the clock
func (c ClockWidget) View() string {
	timeStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.BrightCyan)

	dateStyle := lipgloss.NewStyle().
		Foreground(lipgloss.BrightWhite)

	timStr := timeStyle.Render(c.time.Format("15:04:05"))
	dateStr := dateStyle.Render(c.time.Format("Monday, 2 January 2006"))

	return lipgloss.JoinHorizontal(lipgloss.Center, dateStr, " | ", timStr)
}

func clockTick(duration time.Duration) tea.Cmd {
	return tea.Tick(duration, func(t time.Time) tea.Msg {
		return ClockTickMsg(t)
	})
}
