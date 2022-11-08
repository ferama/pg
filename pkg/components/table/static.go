package table

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/ferama/pg/pkg/conf"
)

type StaticTable struct {
	table     Model
	totalRows int
}

func NewStatic(columns []string) *StaticTable {
	t := New(columns, 0, 0)
	t.Styles.Title = lipgloss.NewStyle().
		Bold(true).
		Background(lipgloss.Color(conf.ColorTitle)).
		Foreground(lipgloss.Color("#000000"))

	st := &StaticTable{
		table: t,
	}
	return st
}

func (t *StaticTable) SetRows(rows []SimpleRow) {
	t.totalRows = len(rows)

	var rs []Row
	for _, row := range rows {
		rs = append(rs, row)
	}
	t.table.SetRows(rs)
}

func (t *StaticTable) Render() string {

	t.table.SetSize(t.table.contentWidth, t.totalRows+1) // +1 include header
	t.table.Styles.Title.
		Width(t.table.contentWidth)
	return t.table.View()
}
