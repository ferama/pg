package statusbar

import (
	"context"
	"fmt"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ferama/pg/pkg/conf"
	"github.com/ferama/pg/pkg/db"
	"github.com/ferama/pg/pkg/utils"
)

var (
	shortcutStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#cccccc")).
			Foreground(lipgloss.Color("#000000")).
			Align(lipgloss.Right).
			Border(lipgloss.ThickBorder(), false, true, false, false).
			BorderForeground(lipgloss.Color(conf.ColorBlur))

	infoStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#bbbbbb")).
			Foreground(lipgloss.Color("#000000")).
			Align(lipgloss.Left).
			Border(lipgloss.ThickBorder(), false, false, false, true).
			BorderForeground(lipgloss.Color(conf.ColorBlur))
)

type connStatusUpdateMsg struct {
}

type Model struct {
	path      *utils.PathParts
	connected bool
	lock      sync.Mutex
}

func New(path *utils.PathParts) *Model {
	m := &Model{
		path:      path,
		connected: false,
	}
	return m
}

func (m *Model) Init() tea.Cmd {
	m.updateConnectionStatus()
	return func() tea.Msg {
		return connStatusUpdateMsg{}
	}
}

func (m *Model) updateConnectionStatus() {
	m.lock.Lock()
	defer m.lock.Unlock()

	conn, err := db.GetDBFromConf(m.path.ConfigConnection, "")
	if err != nil {
		m.connected = false
		return
	}
	defer conn.Close()

	ctx, cancelFn := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelFn()

	if err := conn.PingContext(ctx); err != nil {
		m.connected = false
	} else {
		m.connected = true
	}
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		shortcutStyle = shortcutStyle.Width(msg.Width/2 - 2)
		infoStyle = infoStyle.Width(msg.Width - shortcutStyle.GetWidth() - 2)
	case connStatusUpdateMsg:
		cmd = tea.Tick(3*time.Second, func(t time.Time) tea.Msg {
			m.updateConnectionStatus()
			return connStatusUpdateMsg{}
		})
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	out := "<Ctrl+x> → execute・<ESC> → exit・"

	status := "connected"
	connStatusStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(conf.ColorHeader)).
		Bold(true)

	if !m.connected {
		connStatusStyle.
			Background(lipgloss.Color(conf.ColorError))
		status = "disconneted"
	}
	path := fmt.Sprintf("・%s・%s/%s/%s・",
		status,
		m.path.ConfigConnection,
		m.path.DatabaseName,
		m.path.SchemaName,
	)
	r := lipgloss.JoinHorizontal(lipgloss.Left,
		infoStyle.Render(connStatusStyle.Render(path)),
		shortcutStyle.Render(out))
	return r
}
