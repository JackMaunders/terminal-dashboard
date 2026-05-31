package widgets

import (
	"fmt"
	"time"

	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
)

type SystemTickMsg time.Time

type memoryUsage struct {
	Used  uint64
	Total uint64
}

type diskUsage struct {
	Used  uint64
	Total uint64
}

type SystemWidget struct {
	cpuUsage float64
	memory   memoryUsage
	disk     diskUsage
}

func getCPUInfo() float64 {
	c, err := cpu.Percent(0, false)
	if err != nil {
		panic(err)
	}
	return c[0]
}

func getMemoryInfo() memoryUsage {
	v, err := mem.VirtualMemory()
	if err != nil {
		panic(err)
	}
	return memoryUsage{
		Used:  v.Used,
		Total: v.Total,
	}
}

func getDiskInfo() diskUsage {
	d, err := disk.Usage("/")
	if err != nil {
		panic(err)
	}
	return diskUsage{
		Used:  d.Used,
		Total: d.Total,
	}
}

func NewSystemWidget() SystemWidget {
	return SystemWidget{
		cpuUsage: getCPUInfo(),
		memory:   getMemoryInfo(),
		disk:     getDiskInfo(),
	}
}

func (s SystemWidget) Init() tea.Cmd {
	return systemTick(time.Second * 3)
}

func (s SystemWidget) Update(msg tea.Msg) (SystemWidget, tea.Cmd) {
	switch msg.(type) {
	case SystemTickMsg:
		s.cpuUsage = getCPUInfo()
		s.memory = getMemoryInfo()
		s.disk = getDiskInfo()

		return s, systemTick(time.Second * 2)
	}
	return s, nil
}

func (s SystemWidget) View() string {
	const GiB = 1 << 30

	columns := []table.Column{
		{Title: "CPU Usage", Width: 20},
		{Title: "Memory Usage", Width: 20},
		{Title: "Disk Usage", Width: 20},
	}
	rows := []table.Row{
		{
			fmt.Sprintf("%.2f%%", s.cpuUsage),
			fmt.Sprintf("%.1f/%.1fGB", float64(s.memory.Used)/float64(GiB), float64(s.memory.Total)/float64(GiB)),
			fmt.Sprintf("%dGB / %dGB", s.disk.Used/1e9, s.disk.Total/1e9),
		},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithHeight(7),
		table.WithWidth(100),
	)

	ts := table.DefaultStyles()
	ts.Header = ts.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)

	t.SetStyles(ts)
	return t.View()
}

func systemTick(duration time.Duration) tea.Cmd {
	return tea.Tick(duration, func(t time.Time) tea.Msg {
		return SystemTickMsg(t)
	})
}
