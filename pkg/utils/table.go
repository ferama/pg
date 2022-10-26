package utils

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/viper"
)

func GetTableWriter() table.Writer {
	t := table.NewWriter()

	colors := table.ColorOptions{
		IndexColumn:  text.Colors{text.FgHiBlue, text.BgHiBlack},
		Footer:       text.Colors{text.FgBlue, text.BgHiBlack},
		Header:       text.Colors{text.BgBlue, text.FgBlack},
		Row:          text.Colors{text.FgHiWhite, text.BgBlack},
		RowAlternate: text.Colors{text.FgWhite, text.BgBlack},
	}
	options := table.Options{
		DrawBorder:      false,
		SeparateColumns: true,
		SeparateFooter:  false,
		SeparateHeader:  false,
		SeparateRows:    false,
	}
	styleColor := table.Style{
		Name:    "Custom",
		Box:     table.StyleBoxLight,
		Color:   colors,
		Format:  table.FormatOptionsDefault,
		HTML:    table.DefaultHTMLOptions,
		Options: options,
		Title:   table.TitleOptionsBlackOnBlue,
	}
	if viper.GetBool("useColors") {
		t.SetStyle(styleColor)
	} else {
		t.SetStyle(table.StyleRounded)
	}
	t.SetOutputMirror(os.Stdout)
	return t
}
