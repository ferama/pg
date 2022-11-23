package texteditor

import (
	"fmt"
	"io"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	minHeight     = 1
	minWidth      = 2
	defaultHeight = 6
	defaultWidth  = 40
	maxHeight     = 99
	maxWidth      = 500

	lineDecoratorWidth = 3 // a tilde and a space

	defaultSyntaxColorStyle = "monokai"
)

// KeyMap is the key bindings for different actions within the textarea.
type KeyMap struct {
	CharacterBackward       key.Binding
	CharacterForward        key.Binding
	LineNext                key.Binding
	LinePrevious            key.Binding
	DeleteCharacterBackward key.Binding
	InsertNewline           key.Binding
}

var DefaultKeyMap = KeyMap{
	CharacterForward:        key.NewBinding(key.WithKeys("right", "ctrl+f")),
	CharacterBackward:       key.NewBinding(key.WithKeys("left", "ctrl+b")),
	InsertNewline:           key.NewBinding(key.WithKeys("enter", "ctrl+m")),
	DeleteCharacterBackward: key.NewBinding(key.WithKeys("backspace", "ctrl+h")),
	LineNext:                key.NewBinding(key.WithKeys("down", "ctrl+n")),
	LinePrevious:            key.NewBinding(key.WithKeys("up", "ctrl+p")),
}

type Model struct {
	// The focus status
	focused bool

	// Underlying text value.
	value [][]rune

	// Cursor column.
	col int
	// Cursor row.
	row int

	// KeyMap encodes the keybindings recognized by the widget.
	KeyMap KeyMap

	// Last character offset, used to maintain state when the cursor is moved
	// vertically such that we can maintain the same navigating position.
	lastCharOffset int

	// viewport is the vertically-scrollable viewport of the multi-line text
	// input.
	viewport *viewport.Model

	// Styling. FocusedStyle and BlurredStyle are used to style the textarea in
	// focused and blurred states.
	FocusedStyle Style
	BlurredStyle Style

	// style is the current styling to use.
	// It is used to abstract the differences in focus state when styling the
	// model, since we can simply assign the set of styles to this variable
	// when switching focus states.
	style *Style

	chromaStyle *chroma.Style

	// How many columns should be hidden
	xOffset int

	// if set enable syntax highlight.
	// available syntaxes come from chroma https://github.com/alecthomas/chroma
	syntaxLang string

	// syntax color style
	highlighterStyle string
}

func New() Model {
	vp := viewport.New(0, 0)

	focusedStyle, blurredStyle := DefaultStyles()

	// Determine style.
	chromaStyle := styles.Get(defaultSyntaxColorStyle)
	if chromaStyle == nil {
		chromaStyle = styles.Fallback
	}

	m := Model{
		focused:  true,
		viewport: &vp,

		highlighterStyle: defaultSyntaxColorStyle,
		style:            &blurredStyle,
		chromaStyle:      chromaStyle,
		FocusedStyle:     focusedStyle,
		BlurredStyle:     blurredStyle,

		col: 0,
		row: 0,

		xOffset: 0,
		KeyMap:  DefaultKeyMap,

		value: make([][]rune, minHeight, maxHeight),
	}

	vp.Width = defaultWidth
	vp.Height = defaultHeight
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

// Activate focus on the texteditor
func (m *Model) Focus() tea.Cmd {
	m.focused = true
	m.style = &m.FocusedStyle

	m.style.LineDecorator.Width(lineDecoratorWidth)
	return nil
}

// Deactivate focus
func (m *Model) Blur() tea.Cmd {
	m.focused = false
	m.style = &m.BlurredStyle

	m.style.LineDecorator.Width(lineDecoratorWidth)
	m.style.CursorLine.Width(m.viewport.Width)
	return nil
}

// Reports focus status
func (m Model) Focused() bool {
	return m.focused
}

// SetSyntax lexer
func (m *Model) SetSyntax(syntax string) {
	m.syntaxLang = syntax
}

// Set textarea size
func (m Model) SetSize(width, height int) {
	m.SetWidth(width)
	m.SetHeight(height)
}

// Set textarea width
func (m Model) SetWidth(width int) {
	m.viewport.Width = width

	m.style.CursorLine.Width(width)
}

// Set textarea height
func (m Model) SetHeight(height int) {
	m.viewport.Height = height
}

func (m *Model) splitLine(row, col int) {
	// To perform a split, take the current line and keep the content before
	// the cursor, take the content after the cursor and make it the content of
	// the line underneath, and shift the remaining lines down by one
	head, tailSrc := m.value[row][:col], m.value[row][col:]
	tail := make([]rune, len(tailSrc))
	copy(tail, tailSrc)

	m.value = append(m.value[:row+1], m.value[row:]...)

	m.value[row] = head
	m.value[row+1] = tail

	m.col = 0
	m.row++
}

// mergeLineAbove merges the current line the cursor is on with the line above.
func (m *Model) mergeLineAbove(row int) {
	if row <= 0 {
		return
	}

	m.col = len(m.value[row-1])
	m.row = m.row - 1

	// To perform a merge, we will need to combine the two lines and then
	m.value[row-1] = append(m.value[row-1], m.value[row]...)

	// Shift all lines up by one
	for i := row; i < len(m.value)-1; i++ {
		m.value[i] = m.value[i+1]
	}

	// And, remove the last line
	if len(m.value) > 0 {
		m.value = m.value[:len(m.value)-1]
	}
}

// cursorDown moves the cursor down by one line.
// Returns whether or not the cursor blink should be reset.
func (m *Model) cursorDown() {
	if m.row < len(m.value)-1 {
		m.row++
		if len(m.value[m.row]) <= m.col {
			m.col = len(m.value[m.row])
		}
	}
}

// cursorUp moves the cursor up by one line.
func (m *Model) cursorUp() {
	if m.row > 0 {
		m.row--
		if len(m.value[m.row]) < m.col {
			m.col = len(m.value[m.row])
		}
	}
}

// Reset sets the input to its default state with no input.
func (m *Model) Reset() {
	m.value = make([][]rune, minHeight, maxHeight)
	m.col = 0
	m.row = 0
	m.viewport.GotoTop()
	m.SetCursor(0)
}

// SetValue sets the value of the text input.
func (m *Model) SetValue(s string) {
	m.Reset()
	m.InsertString(s)
}

// InsertString inserts a string at the cursor position.
func (m *Model) InsertString(s string) {
	lines := strings.Split(s, "\n")
	for l, line := range lines {
		for _, rune := range line {
			m.InsertRune(rune)
		}
		if l != len(lines)-1 {
			m.InsertRune('\n')
		}
	}
}

// InsertRune inserts a rune at the cursor position.
func (m *Model) InsertRune(r rune) {
	if r == '\n' {
		m.splitLine(m.row, m.col)
		return
	}

	m.value[m.row] = append(m.value[m.row][:m.col], append([]rune{r}, m.value[m.row][m.col:]...)...)
	m.col++
}

// cursorRight moves the cursor one character to the right.
func (m *Model) cursorRight() {
	if m.col < len(m.value[m.row]) {
		m.SetCursor(m.col + 1)
	} else {
		if m.row < len(m.value)-1 {
			m.row++
			m.CursorStart()
		}
	}
}

// cursorLeft moves the cursor one character to the left.
// If insideLine is set, the cursor is moved to the last
// character in the previous line, instead of one past that.
func (m *Model) cursorLeft(insideLine bool) {
	if m.col == 0 && m.row != 0 {
		m.row--
		m.CursorEnd()
		if !insideLine {
			return
		}
	}
	if m.col > 0 {
		m.SetCursor(m.col - 1)
	}
}

// repositionView repositions the view of the viewport based on the defined
// scrolling behavior.
func (m *Model) repositionView() {
	min := m.viewport.YOffset
	max := min + m.viewport.Height - 1

	if row := m.row; row < min {
		m.viewport.LineUp(min - row)
	} else if row > max {
		m.viewport.LineDown(row - max)
	}
}

// CursorEnd moves the cursor to the end of the input field.
func (m *Model) CursorEnd() {
	m.SetCursor(len(m.value[m.row]))
}

// CursorStart moves the cursor to the start of the input field.
func (m *Model) CursorStart() {
	m.SetCursor(0)
}

// SetCursor moves the cursor to the given position. If the position is
// out of bounds the cursor will be moved to the start or end accordingly.
func (m *Model) SetCursor(col int) {
	m.col = clamp(col, 0, len(m.value[m.row]))
	// Any time that we move the cursor horizontally we need to reset the last
	// offset so that the horizontal position when navigating is adjusted.
	m.lastCharOffset = 0
}

// updates current x offset using info from current
// cursor position and underlying viewport size
func (m *Model) updateXOffset() {
	diff := m.col - m.viewport.Width + lineDecoratorWidth + 1
	if diff >= 0 {
		m.xOffset = diff
	} else {
		m.xOffset = 0
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if !m.focused {
		return m, nil
	}

	var cmds []tea.Cmd
	// var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.CharacterForward):
			m.cursorRight()
		case key.Matches(msg, m.KeyMap.CharacterBackward):
			m.cursorLeft(false /* insideLine */)
		case key.Matches(msg, m.KeyMap.DeleteCharacterBackward):
			m.col = clamp(m.col, 0, len(m.value[m.row]))
			if m.col <= 0 {
				m.mergeLineAbove(m.row)
				break
			}
			if len(m.value[m.row]) > 0 {
				m.value[m.row] = append(m.value[m.row][:max(0, m.col-1)], m.value[m.row][m.col:]...)
				if m.col > 0 {
					m.SetCursor(m.col - 1)
				}
			}
		case key.Matches(msg, m.KeyMap.LineNext):
			m.cursorDown()
		case key.Matches(msg, m.KeyMap.LinePrevious):
			m.cursorUp()
		case key.Matches(msg, m.KeyMap.InsertNewline):
			if len(m.value) >= maxHeight {
				return m, nil
			}
			m.col = clamp(m.col, 0, len(m.value[m.row]))
			m.splitLine(m.row, m.col)
		default:
			m.col = min(m.col, len(m.value[m.row]))
			m.value[m.row] = append(m.value[m.row][:m.col], append(msg.Runes, m.value[m.row][m.col:]...)...)
			m.col = clamp(m.col+len(msg.Runes), 0, len(m.value[m.row]))
		}
	}
	m.updateXOffset()
	m.repositionView()

	vp, cmd := m.viewport.Update(msg)
	m.viewport = &vp
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// Value returns the value of the text input.
func (m Model) Value() string {
	if m.value == nil {
		return ""
	}

	var v strings.Builder
	for _, l := range m.value {
		v.WriteString(string(l))
		v.WriteByte('\n')
	}

	return strings.TrimSuffix(v.String(), "\n")
}

// renderLine renders code line applying syntax highlights and handling cursor
func (m Model) renderLine(w io.Writer, source string, hasCursor bool) error {
	lexer := m.syntaxLang

	// Determine lexer.
	l := lexers.Get(lexer)
	if l == nil {
		l = lexers.Analyse(source)
	}
	if l == nil {
		l = lexers.Fallback
	}
	l = chroma.Coalesce(l)

	it, err := l.Tokenise(nil, source)
	if err != nil {
		return err
	}

	defer func() {
		if perr := recover(); perr != nil {
			err = perr.(error)
		}
	}()

	theme := clearBackground(m.chromaStyle)
	column := 0
	doneWithCursor := false
	cursorColumn := m.col
	for token := it(); token != chroma.EOF; token = it() {
		columnNext := column + len(token.Value)

		// do not render offsetted columns
		if column < m.xOffset {
			if columnNext >= m.xOffset {
				diff := m.xOffset - column
				token.Value = token.Value[diff:]
				cursorColumn += diff
				columnNext += diff
			} else {
				column += len(token.Value)
				continue
			}
		}

		entry := theme.Get(token.Type)
		fmt.Fprint(w, m.applyTheme(entry, hasCursor))

		if hasCursor && columnNext > cursorColumn && !doneWithCursor {
			pos := cursorColumn - column
			tv := token.Value
			preCursor := tv[0:pos]
			cursor := tv[pos : pos+1]
			postCursor := tv[pos+1:]

			fmt.Fprint(w, preCursor)
			fmt.Fprint(w, m.style.Cursor.Render(cursor))

			// reapply theme resetted by cursor
			fmt.Fprint(w, m.applyTheme(entry, hasCursor))
			fmt.Fprint(w, postCursor)
			doneWithCursor = true
		} else {
			fmt.Fprint(w, token.Value)
		}

		if !entry.IsZero() {
			fmt.Fprint(w, "\033[0m")
		}

		column = columnNext
	}
	return nil
}

func (m Model) applyTheme(entry chroma.StyleEntry, hasCursor bool) string {
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
		} else {
			if hasCursor {
				r, g, b, _ := m.style.CursorLine.GetBackground().RGBA()
				out += fmt.Sprintf("\033[48;2;%d;%d;%dm", r, g, b)
			}
		}
	}
	return out
}

func (m Model) View() string {
	sb := new(strings.Builder)

	for ir, r := range m.value {
		haveCursor := false
		if ir == m.row {
			haveCursor = true
		}
		sb.WriteString(m.style.LineDecorator.Render(fmt.Sprint(ir + 1)))
		lsb := new(strings.Builder)

		for _, c := range r {
			lsb.WriteRune(c)
		}

		if m.col >= len(r) {
			lsb.WriteString(" ")
		}

		hlsbs := new(strings.Builder)
		m.renderLine(
			hlsbs,
			lsb.String(),
			haveCursor,
		)
		if haveCursor {
			sb.WriteString(m.style.CursorLine.Render(hlsbs.String()))
		} else {
			sb.WriteString(hlsbs.String())
		}
		sb.WriteString("\n")
	}

	m.viewport.SetContent(sb.String())
	return m.style.Base.Render(m.viewport.View())
}
