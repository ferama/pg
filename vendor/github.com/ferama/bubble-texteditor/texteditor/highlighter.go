package texteditor

import (
	"fmt"
	"io"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

// Clear the background colour.
func clearBackground(style *chroma.Style) *chroma.Style {
	builder := style.Builder()
	bg := builder.Get(chroma.Background)
	bg.Background = 0
	bg.NoInherit = true
	builder.AddEntry(chroma.Background, bg)
	style, _ = builder.Build()
	return style
}

func applyTheme(entry chroma.StyleEntry) string {
	out := ""

	if !entry.IsZero() {
		if entry.Bold == chroma.Yes {
			out += "\033[1m"
		}
		if entry.Underline == chroma.Yes {
			out += "\033[4m"
		}
		if entry.Italic == chroma.Yes {
			out += "\033[3m"
		}
		if entry.Colour.IsSet() {
			out += fmt.Sprintf("\033[38;2;%d;%d;%dm", entry.Colour.Red(), entry.Colour.Green(), entry.Colour.Blue())
		}
		if entry.Background.IsSet() {
			out += fmt.Sprintf("\033[48;2;%d;%d;%dm", entry.Background.Red(), entry.Background.Green(), entry.Background.Blue())
		}
	}
	return out
}

func format(w io.Writer,
	theme *chroma.Style,
	it chroma.Iterator,
	hasCursor bool,
	cursorColumn int,
	style *Style,
	xOffset int) error {

	theme = clearBackground(theme)

	column := 0
	doneWithCursor := false
	for token := it(); token != chroma.EOF; token = it() {
		columnNext := column + len(token.Value)

		// do not render offsetted columns
		if column < xOffset {
			if columnNext >= xOffset {
				diff := xOffset - column
				token.Value = token.Value[diff:]
				cursorColumn += diff
				columnNext += diff
			} else {
				column += len(token.Value)
				continue
			}
		}

		entry := theme.Get(token.Type)
		fmt.Fprint(w, applyTheme(entry))

		if hasCursor && columnNext > cursorColumn && !doneWithCursor {
			pos := cursorColumn - column
			tv := token.Value
			preCursor := tv[0:pos]
			cursor := tv[pos : pos+1]
			postCursor := tv[pos+1:]

			fmt.Fprint(w, preCursor)
			fmt.Fprint(w, style.Cursor.Render(cursor))

			// reapply theme resetted by cursor
			fmt.Fprint(w, applyTheme(entry))
			fmt.Fprint(w, postCursor)
			doneWithCursor = true
		} else {
			// for _, c := range token.Value {
			// 	fmt.Fprint(w, applyTheme(entry))
			// 	fmt.Fprintf(w, "%c", c)
			// }
			fmt.Fprint(w, token.Value)
		}

		if !entry.IsZero() {
			fmt.Fprint(w, "\033[0m")
		}

		column = columnNext
	}
	return nil
}

// renderLine renders code line applying syntax highlights and handling cursor
func renderLine(
	w io.Writer,
	source,
	lexer,
	theme string,
	hasCursor bool,
	cursorColumn int,
	style *Style,
	xOffset int) (err error) {

	// Determine lexer.
	l := lexers.Get(lexer)
	if l == nil {
		l = lexers.Analyse(source)
	}
	if l == nil {
		l = lexers.Fallback
	}
	l = chroma.Coalesce(l)

	// Determine style.
	s := styles.Get(theme)
	if s == nil {
		s = styles.Fallback
	}

	it, err := l.Tokenise(nil, source)
	if err != nil {
		return err
	}

	defer func() {
		if perr := recover(); perr != nil {
			err = perr.(error)
		}
	}()
	return format(w, s, it, hasCursor, cursorColumn, style, xOffset)
}
