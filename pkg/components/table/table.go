package table

import (
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ferama/pg/pkg/conf"
	"github.com/juju/ansiterm/tabwriter"
)

// Row renderer.
type Row interface {
	// Render the row into the given tabwriter.
	// To render correctly, join each cell by a tab character '\t'.
	// Use `m.Cursor() == index` to determine if the row is selected.
	// Take a look at the `SimpleRow` implementation for an example.
	Render(w io.Writer, model Model, index int)
}

// SimpleRow is a set of cells that can be rendered into a table.
// It supports row highlight if selected.
type SimpleRow []any

// Render a simple row.
func (row SimpleRow) Render(w io.Writer, model Model, index int) {
	cells := make([]string, 0)
	for i, v := range row {
		if i < model.colCursor {
			continue
		}
		cells = append(cells, model.Styles.Cell.Render(fmt.Sprintf("%v", v)))
	}
	s := strings.Join(cells, "\t")

	if index == model.Cursor() && model.focused {
		s = model.Styles.SelectedRow.Render(s)
	} else {
		s = model.Styles.UnselectedRow.Render(s)
	}
	fmt.Fprintln(w, s)
}

// New model.
func New(cols []string, width, height int) Model {
	vp := viewport.New(width, maxInt(height-1, 0))
	tw := &tabwriter.Writer{}
	return Model{
		KeyMap:    DefaultKeyMap(),
		Styles:    DefaultStyles(),
		focused:   false,
		cols:      cols,
		header:    strings.Join(cols, " "), // simple initial header view without tabwriter.
		colCursor: 0,
		viewPort:  vp,
		tabWriter: tw,
	}
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Model of a table component.
type Model struct {
	KeyMap KeyMap
	Styles Styles

	focused      bool
	cols         []string
	rows         []Row
	header       string
	viewPort     viewport.Model
	tabWriter    *tabwriter.Writer
	cursor       int
	colCursor    int
	contentWidth int
}

// KeyMap holds the key bindings for the table.
type KeyMap struct {
	End      key.Binding
	Home     key.Binding
	PageDown key.Binding
	PageUp   key.Binding
	Down     key.Binding
	Up       key.Binding
	Right    key.Binding
	Left     key.Binding
}

// DefaultKeyMap used by the `New` constructor.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		End: key.NewBinding(
			key.WithKeys("end"),
			key.WithHelp("end", "bottom"),
		),
		Home: key.NewBinding(
			key.WithKeys("home"),
			key.WithHelp("home", "top"),
		),
		PageDown: key.NewBinding(
			key.WithKeys("pgdown", "ctrl+down"),
			key.WithHelp("pgdown", "page down"),
		),
		PageUp: key.NewBinding(
			key.WithKeys("pgup", "ctrl+up"),
			key.WithHelp("pgup", "page up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("↓", "down"),
		),
		Up: key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("↑", "up"),
		),
		Right: key.NewBinding(
			key.WithKeys("right"),
			key.WithHelp("→", "right"),
		),
		Left: key.NewBinding(
			key.WithKeys("left"),
			key.WithHelp("←", "left"),
		),
	}
}

// Styles holds the styling for the table.
type Styles struct {
	Header        lipgloss.Style
	Cell          lipgloss.Style
	HeaderCell    lipgloss.Style
	UnselectedRow lipgloss.Style
	SelectedRow   lipgloss.Style
}

// DefaultStyles used by the `New` constructor.
func DefaultStyles() Styles {
	return Styles{
		Header: lipgloss.NewStyle().
			Bold(true),
		Cell: lipgloss.NewStyle().
			Padding(0, 1, 0, 1),
		HeaderCell: lipgloss.NewStyle().
			Padding(0, 1, 0, 1),
		UnselectedRow: lipgloss.NewStyle(),
		SelectedRow: lipgloss.NewStyle().
			Bold(true).
			Background(lipgloss.Color(conf.ColorTableRowFocus)).
			Foreground(lipgloss.Color("#000000")),
	}
}

// SetSize of the table and makes sure to update the view
// and the selected row does not go out of bounds.
func (m *Model) SetSize(width, height int) {
	m.viewPort.Width = width
	m.viewPort.Height = height - 1

	// the sum here is not exactly right. it is an almost an upper bound
	m.Styles.Header.Width(m.contentWidth + m.viewPort.Width)
	m.Styles.SelectedRow.Width(m.contentWidth + m.viewPort.Width)
	m.Styles.UnselectedRow.Width(m.contentWidth + m.viewPort.Width)

	if m.cursor > m.viewPort.YOffset+m.viewPort.Height-1 {
		m.cursor = m.viewPort.YOffset + m.viewPort.Height - 1
		m.updateView()
	}
}

func (m *Model) GetViewportWidth() int {
	return m.viewPort.Width
}

func (m *Model) GetViewportHeight() int {
	return m.viewPort.Height
}

// Cursor returns the index of the selected row.
func (m Model) Cursor() int {
	return m.cursor
}

// SelectedRow returns the selected row.
// You can cast it to your own implementation.
func (m Model) SelectedRow() Row {
	return m.rows[m.cursor]
}

// SetRows of the table and makes sure to update the view
// and the selected row does not go out of bounds.
func (m *Model) SetRows(rows []Row) {
	m.rows = rows
	m.updateView()
}

func (m *Model) updateView() {
	var b strings.Builder
	m.tabWriter.Init(&b, 0, 4, 1, ' ', 0)

	cols := make([]string, 0)
	for i, c := range m.cols {
		if i < m.colCursor {
			continue
		}
		cols = append(cols, m.Styles.HeaderCell.Render(c))
	}
	// rendering the header.
	s := strings.Join(cols, "\t")
	s = m.Styles.Header.Render(s)
	fmt.Fprintln(m.tabWriter, s)

	// rendering the rows.
	for i, row := range m.rows {
		row.Render(m.tabWriter, *m, i)
	}

	m.tabWriter.Flush()

	content := b.String()
	m.contentWidth = lipgloss.Width(content)

	// split table at first line-break to take header and rows apart.
	parts := strings.SplitN(content, "\n", 2)
	if len(parts) != 0 {
		m.header = parts[0]
		if len(parts) == 2 {
			m.viewPort.SetContent(strings.TrimRightFunc(parts[1], unicode.IsSpace))
		}
	}
}

// CursorIsAtTop of the table.
func (m Model) CursorIsAtTop() bool {
	return m.cursor == 0
}

// CursorIsAtBottom of the table.
func (m Model) CursorIsAtBottom() bool {
	return m.cursor == len(m.rows)-1
}

// CursorIsPastBottom of the table.
func (m Model) CursorIsPastBottom() bool {
	return m.cursor > len(m.rows)-1
}

// GoUp moves the selection to the previous row.
// It can not go above the first row.
func (m *Model) GoUp() {
	if m.CursorIsAtTop() {
		return
	}

	m.cursor--
	m.updateView()

	if m.cursor < m.viewPort.YOffset {
		m.viewPort.LineUp(1)
	}
}

// GoDown moves the selection to the next row.
// It can not go below the last row.
func (m *Model) GoDown() {
	if m.CursorIsAtBottom() {
		return
	}

	m.cursor++
	m.updateView()

	if m.cursor > m.viewPort.YOffset+m.viewPort.Height-1 {
		m.viewPort.LineDown(1)
	}
}

// GoPageUp moves the selection one page up.
// It can not go above the first row.
func (m *Model) GoPageUp() {
	if m.CursorIsAtTop() {
		return
	}

	m.cursor -= m.viewPort.Height
	if m.cursor < 0 {
		m.cursor = 0
	}

	m.updateView()

	m.viewPort.ViewUp()
}

// GoPageDown moves the selection one page down.
// It can not go below the last row.
func (m *Model) GoPageDown() {
	if m.CursorIsAtBottom() {
		return
	}

	m.cursor += m.viewPort.Height
	if m.CursorIsPastBottom() {
		m.cursor = len(m.rows) - 1
	}

	m.updateView()

	m.viewPort.ViewDown()
}

// GoTop moves the selection to the first row.
func (m *Model) GoTop() {
	if m.CursorIsAtTop() {
		return
	}

	m.cursor = 0
	m.updateView()
	m.viewPort.GotoTop()
}

func (m *Model) Focus() {
	m.focused = true
	m.updateView()
}

func (m *Model) Blur() {
	m.focused = false
	m.updateView()
}

// GoBottom moves the selection to the last row.
func (m *Model) GoBottom() {
	if m.CursorIsAtBottom() {
		return
	}

	m.cursor = len(m.rows) - 1
	m.updateView()
	m.viewPort.GotoBottom()
}

func (m *Model) GoRight() {
	if m.colCursor+1 >= len(m.cols) {
		return
	}
	m.colCursor++
	m.updateView()
}

func (m *Model) GoLeft() {
	if m.colCursor-1 < 0 {
		return
	}
	m.colCursor--
	m.updateView()
}

// Update tea.Model implementor.
// It handles the key events.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !m.focused {
			break
		}
		switch {
		case key.Matches(msg, m.KeyMap.Up):
			m.GoUp()
		case key.Matches(msg, m.KeyMap.Down):
			m.GoDown()
		case key.Matches(msg, m.KeyMap.PageUp):
			m.GoPageUp()
		case key.Matches(msg, m.KeyMap.PageDown):
			m.GoPageDown()
		case key.Matches(msg, m.KeyMap.Home):
			m.GoTop()
		case key.Matches(msg, m.KeyMap.End):
			m.GoBottom()
		case key.Matches(msg, m.KeyMap.Right):
			m.GoRight()
		case key.Matches(msg, m.KeyMap.Left):
			m.GoLeft()
		}
	}

	return m, nil
}

// View tea.Model implementors.
// It renders the table inside a viewport.
func (m Model) View() string {
	return lipgloss.NewStyle().MaxWidth(m.viewPort.Width).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			m.header,
			m.viewPort.View(),
		),
	)
}
