package widgets

import (
	"time"

	tea "charm.land/bubbletea/v2"
)

type tickMsg time.Time

func tick(duration time.Duration) tea.Cmd {
	return tea.Tick(duration, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
