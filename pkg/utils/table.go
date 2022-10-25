package utils

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

func GetTableWriter() table.Writer {
	t := table.NewWriter()
	t.SetStyle(table.StyleRounded)

	// colors := table.ColorOptions{
	// 	IndexColumn:  text.Colors{text.FgHiBlue, text.BgHiBlack},
	// 	Footer:       text.Colors{text.FgBlue, text.BgHiBlack},
	// 	Header:       text.Colors{text.BgBlue, text.FgBlack},
	// 	Row:          text.Colors{text.FgHiWhite, text.BgBlack},
	// 	RowAlternate: text.Colors{text.FgWhite, text.BgBlack},
	// }
	// style := table.Style{
	// 	Name:    "Custom",
	// 	Box:     table.StyleBoxDefault,
	// 	Color:   colors,
	// 	Format:  table.FormatOptionsDefault,
	// 	HTML:    table.DefaultHTMLOptions,
	// 	Options: table.OptionsNoBordersAndSeparators,
	// 	Title:   table.TitleOptionsBlackOnBlue,
	// }
	// t.SetStyle(style)
	t.SetOutputMirror(os.Stdout)
	return t
}
