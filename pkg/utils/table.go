package utils

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

func GetTableWriter() table.Writer {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	return t
}
