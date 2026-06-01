package widgets

import (
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type ClockTickMsg time.Time

type ClockWidget struct {
	time time.Time
}

func NewClockWidget() ClockWidget {
	return ClockWidget{time: time.Now()}
}

func (c ClockWidget) Init() tea.Cmd {
	return clockTick(time.Second)
}

func (c ClockWidget) Update(msg tea.Msg) (ClockWidget, tea.Cmd) {
	switch msg.(type) {
	case ClockTickMsg:
		c.time = time.Now()
		return c, clockTick(time.Second)
	}
	return c, nil
}

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
