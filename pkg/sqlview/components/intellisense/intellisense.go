package intellisense

import (
	"fmt"
	"strings"
	"unicode"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ferama/pg/pkg/conf"
	"github.com/ferama/pg/pkg/db"
	"github.com/ferama/pg/pkg/utils"
)

var sqlKeywords = []string{
	"ALL", "ALTER", "AND", "AS", "ASC", "AVG",
	"BETWEEN", "BY",
	"CASE", "CAST", "COALESCE", "COUNT", "CREATE", "CROSS",
	"DELETE", "DESC", "DISTINCT", "DROP",
	"ELSE", "END", "EXISTS",
	"FULL", "FROM",
	"GROUP",
	"HAVING",
	"ILIKE", "IN", "INNER", "INSERT", "INTO", "IS",
	"JOIN",
	"LEFT", "LIKE", "LIMIT",
	"MAX", "MIN",
	"NOT", "NULL", "NULLIF",
	"OFFSET", "ON", "OR", "ORDER", "OUTER",
	"RETURNING", "RIGHT",
	"SELECT", "SET", "SUM",
	"TABLE", "THEN",
	"UNION", "UPDATE",
	"VALUES", "VIEW",
	"WHEN", "WHERE", "WITH",
}

const maxVisible = 8

var (
	popupStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color(conf.ColorFocus))

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#000000")).
			Background(lipgloss.Color(conf.ColorFocus))

	normalStyle = lipgloss.NewStyle()

	hintStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888"))
)

// SchemaLoadedMsg carries table and column data loaded from the database.
type SchemaLoadedMsg struct {
	Tables  []string
	Columns []string
}

type Model struct {
	path    *utils.PathParts
	visible bool
	cursor  int
	offset  int

	items []string

	tables  []string
	columns []string
	loaded  bool

	// width is the outer width the popup must match (= editor outer width).
	width int
}

func New(path *utils.PathParts) *Model {
	return &Model{
		path: path,
	}
}

// LoadSchema returns a Cmd that queries information_schema and returns the
// results inside SchemaLoadedMsg. Works even when SchemaName is empty
// (loads tables from all non-system schemas).
func (m *Model) LoadSchema() tea.Cmd {
	if m.path == nil || (m.path.ConfigConnection == "" && m.path.DatabaseName == "") {
		return nil
	}
	path := m.path
	return func() tea.Msg {
		var tables, columns []string

		schemaFilter := buildSchemaFilter(path.SchemaName)

		tablesQuery := fmt.Sprintf(
			`SELECT table_name FROM information_schema.tables
			 WHERE table_type = 'BASE TABLE' %s
			 ORDER BY table_name`,
			schemaFilter)
		if r, err := db.Query(path.ConfigConnection, path.DatabaseName, path.SchemaName, tablesQuery); err == nil {
			for _, row := range r.Rows {
				if len(row) > 0 {
					tables = append(tables, row[0])
				}
			}
		}

		colQuery := fmt.Sprintf(
			`SELECT DISTINCT column_name FROM information_schema.columns
			 WHERE TRUE %s
			 ORDER BY column_name`,
			schemaFilter)
		if r, err := db.Query(path.ConfigConnection, path.DatabaseName, path.SchemaName, colQuery); err == nil {
			for _, row := range r.Rows {
				if len(row) > 0 {
					columns = append(columns, row[0])
				}
			}
		}

		return SchemaLoadedMsg{Tables: tables, Columns: columns}
	}
}

func buildSchemaFilter(schemaName string) string {
	if schemaName != "" {
		return fmt.Sprintf("AND table_schema = '%s'", schemaName)
	}
	return "AND table_schema NOT IN ('pg_catalog', 'information_schema', 'pg_toast')"
}

// SetSchemaData is called from the editor's Update when SchemaLoadedMsg arrives.
func (m *Model) SetSchemaData(tables, columns []string) {
	m.tables = tables
	m.columns = columns
	m.loaded = true
}

func (m *Model) Show(prefix string) {
	filtered := m.filter(prefix)
	if len(filtered) == 0 {
		m.visible = false
		return
	}
	m.items = filtered
	m.visible = true
	m.cursor = 0
	m.offset = 0
}

func (m *Model) Hide() {
	m.visible = false
	m.items = nil
	m.cursor = 0
	m.offset = 0
}

func (m *Model) IsVisible() bool {
	return m.visible && len(m.items) > 0
}

// Height returns the number of terminal rows occupied by the popup (0 when hidden).
func (m *Model) Height() int {
	if !m.IsVisible() {
		return 0
	}
	rows := len(m.items)
	if rows > maxVisible {
		rows = maxVisible
	}
	return rows + 2
}

func (m *Model) MoveUp() {
	if m.cursor > 0 {
		m.cursor--
		if m.cursor < m.offset {
			m.offset--
		}
	}
}

func (m *Model) MoveDown() {
	if m.cursor < len(m.items)-1 {
		m.cursor++
		if m.cursor >= m.offset+maxVisible {
			m.offset++
		}
	}
}

func (m *Model) Selected() string {
	if !m.IsVisible() || m.cursor >= len(m.items) {
		return ""
	}
	return m.items[m.cursor]
}

// SetWidth receives the OUTER width of the editor so the popup can match it.
func (m *Model) SetWidth(w int) {
	m.width = w
}

func (m *Model) filter(prefix string) []string {
	upper := strings.ToUpper(prefix)
	var results []string

	for _, kw := range sqlKeywords {
		if strings.HasPrefix(kw, upper) {
			results = append(results, matchCase(prefix, kw))
		}
	}

	tableSet := make(map[string]bool, len(m.tables))
	for _, t := range m.tables {
		tableSet[strings.ToUpper(t)] = true
		if strings.HasPrefix(strings.ToUpper(t), upper) {
			results = append(results, t)
		}
	}

	for _, c := range m.columns {
		if strings.HasPrefix(strings.ToUpper(c), upper) && !tableSet[strings.ToUpper(c)] {
			results = append(results, c)
		}
	}

	return results
}

func matchCase(prefix, keyword string) string {
	if prefix == "" {
		return keyword
	}
	if unicode.IsLower(rune(prefix[0])) {
		return strings.ToLower(keyword)
	}
	return keyword
}

func (m *Model) View() string {
	if !m.IsVisible() {
		return ""
	}

	end := m.offset + maxVisible
	if end > len(m.items) {
		end = len(m.items)
	}
	visible := m.items[m.offset:end]

	contentW := m.width - 2
	if contentW < 10 {
		contentW = 40
	}

	var lines []string
	for i, item := range visible {
		absIdx := m.offset + i
		padded := fmt.Sprintf(" %-*s", contentW-1, item)
		if absIdx == m.cursor {
			lines = append(lines, selectedStyle.Render(padded))
		} else {
			lines = append(lines, normalStyle.Render(padded))
		}
	}

	hint := ""
	if len(m.items) > maxVisible {
		hint = hintStyle.Render(fmt.Sprintf(" (%d/%d)", m.cursor+1, len(m.items)))
	}

	content := strings.Join(lines, "\n")
	if hint != "" {
		content += "\n" + hint
	}

	return popupStyle.Width(contentW).Render(content)
}
