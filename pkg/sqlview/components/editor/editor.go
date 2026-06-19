package editor

import (
	"time"
	"unicode"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ferama/bubble-texteditor/texteditor"
	"github.com/ferama/pg/pkg/conf"
	"github.com/ferama/pg/pkg/db"
	"github.com/ferama/pg/pkg/history"
	"github.com/ferama/pg/pkg/sqlview/components/hbrowser"
	"github.com/ferama/pg/pkg/sqlview/components/intellisense"
	"github.com/ferama/pg/pkg/utils"
)

var (
	style = lipgloss.NewStyle().
		Border(lipgloss.ThickBorder(), false, true, false, true).
		BorderForeground(lipgloss.Color(conf.ColorFocus))
)

type QueryStatusMsg struct {
	Content string
	Elapsed time.Duration
}

type QueryResultsMsg struct {
	Results *db.QueryResults
}

// IntelliSenseHeightMsg is emitted when the intellisense popup appears or
// disappears so the parent view can adjust the results area height.
type IntelliSenseHeightMsg struct {
	Height int
}

type Model struct {
	path       *utils.PathParts
	texteditor texteditor.Model

	history *history.History
	err     error

	intellisense    *intellisense.Model
	currentWord     string // word at cursor, recomputed from Value() on every char key
	prevPopupHeight int
}

func New(path *utils.PathParts) *Model {
	te := texteditor.New()

	te.SetSyntax("sql")
	te.SetWidth(10)
	te.SetHeight(conf.SqlTextareaHeight)
	te.Focus()

	return &Model{
		path:         path,
		texteditor:   te,
		history:      history.GetInstance(),
		intellisense: intellisense.New(path),
		err:          nil,
	}
}

func (m *Model) Focus() tea.Cmd {
	style.
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color(conf.ColorFocus))
	return m.texteditor.Focus()
}

func (m *Model) Blur() {
	style.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(conf.ColorBlur))
	m.texteditor.Blur()
}

func (m *Model) SetValue(value string) {
	m.texteditor.SetValue(value)
}

func (m *Model) value() string {
	return m.texteditor.Value()
}

func (m *Model) IntelliSenseVisible() bool {
	return m.intellisense.IsVisible()
}

func (m *Model) sqlExecute(connString, dbName, schema, query string) (*db.QueryResults, error) {
	results, err := db.Query(connString, dbName, schema, query)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (m *Model) doQuery() tea.Cmd {
	return func() tea.Msg {
		query := m.value()

		m.history.Append(query)

		results, err := m.sqlExecute(
			m.path.ConfigConnection,
			m.path.DatabaseName,
			m.path.SchemaName,
			query,
		)
		if err != nil {
			return QueryStatusMsg{Content: err.Error()}
		}
		if len(results.Columns) == 0 {
			return QueryStatusMsg{Content: "done", Elapsed: results.Elapsed}
		}
		return QueryResultsMsg{Results: results}
	}
}

func (m *Model) emitHeightChange() tea.Cmd {
	h := m.intellisense.Height()
	if h == m.prevPopupHeight {
		return nil
	}
	m.prevPopupHeight = h
	return func() tea.Msg { return IntelliSenseHeightMsg{Height: h} }
}

// refreshCompletion reads Value() after the texteditor has processed a key,
// extracts the word ending at the cursor, and shows/hides the popup.
// This is called only for character-input keys (not navigation).
func (m *Model) refreshCompletion() tea.Cmd {
	m.currentWord = wordFromValueTail(m.texteditor.Value())
	if len(m.currentWord) >= 2 {
		m.intellisense.Show(m.currentWord)
	} else {
		m.intellisense.Hide()
	}
	return m.emitHeightChange()
}

func (m *Model) hideCompletion() tea.Cmd {
	m.currentWord = ""
	m.intellisense.Hide()
	return m.emitHeightChange()
}

// acceptCompletion replaces the typed prefix with the selected suggestion.
func (m *Model) acceptCompletion(selected string) {
	value := m.texteditor.Value()
	runes := []rune(value)
	wordRunes := []rune(m.currentWord)
	if len(wordRunes) <= len(runes) {
		newVal := string(runes[:len(runes)-len(wordRunes)]) + selected
		m.texteditor.SetValue(newVal)
	}
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		textarea.Blink,
		m.intellisense.LoadSchema(),
	)
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case intellisense.SchemaLoadedMsg:
		m.intellisense.SetSchemaData(msg.Tables, msg.Columns)

	case hbrowser.HBrowserSelectedMsg:
		q, err := m.history.GetAtIdx(msg.Idx)
		if err == nil {
			m.texteditor.SetValue(q)
		}
		cmds = append(cmds, m.hideCompletion())

	case tea.WindowSizeMsg:
		m.texteditor.SetWidth(msg.Width - 2)
		style.Width(msg.Width - 2)
		// Pass full terminal width; popup subtracts its own border internally.
		m.intellisense.SetWidth(msg.Width)

	case tea.KeyMsg:
		// When popup is visible, intercept navigation and acceptance keys.
		if m.intellisense.IsVisible() {
			switch msg.Type {
			case tea.KeyUp:
				m.intellisense.MoveUp()
				return m, nil
			case tea.KeyDown:
				m.intellisense.MoveDown()
				return m, nil
			case tea.KeyTab, tea.KeyEnter:
				selected := m.intellisense.Selected()
				if selected != "" {
					m.acceptCompletion(selected)
				}
				m.currentWord = ""
				m.intellisense.Hide()
				cmds = append(cmds, m.emitHeightChange())
				return m, tea.Batch(cmds...)
			case tea.KeyEsc:
				cmds = append(cmds, m.hideCompletion())
				// Consume Esc so sqlview doesn't quit.
				return m, tea.Batch(cmds...)
			}
		}

		switch msg.Type {
		case tea.KeyShiftDown:
			q, err := m.history.GoNext()
			if err == nil {
				m.texteditor.SetValue(q)
			}
			cmds = append(cmds, m.hideCompletion())

		case tea.KeyShiftUp:
			q, err := m.history.GoPrev()
			if err == nil {
				m.texteditor.SetValue(q)
			}
			cmds = append(cmds, m.hideCompletion())

		case tea.KeyCtrlD:
			m.SetValue("")
			cmds = append(cmds, m.hideCompletion())

		case tea.KeyCtrlX:
			cmds = append(cmds,
				func() tea.Msg { return QueryStatusMsg{Content: "running query..."} },
				m.doQuery(),
			)

		case tea.KeyNull:
			// Ctrl+Space: manual trigger using current word from value.
			m.currentWord = wordFromValueTail(m.texteditor.Value())
			if len(m.currentWord) >= 1 {
				m.intellisense.Show(m.currentWord)
				cmds = append(cmds, m.emitHeightChange())
			}
			// Do not forward NUL to the texteditor.
			return m, tea.Batch(cmds...)

		case tea.KeyLeft, tea.KeyRight, tea.KeyUp, tea.KeyDown,
			tea.KeyHome, tea.KeyEnd, tea.KeyPgUp, tea.KeyPgDown:
			// Cursor moved: we no longer know where the word is. Hide popup.
			cmds = append(cmds, m.hideCompletion())

		default:
			// All character input (printable chars, backspace, space, enter).
			// Forward to texteditor FIRST so Value() is up-to-date, then
			// recompute the word at the cursor from the new value.
			m.texteditor, cmd = m.texteditor.Update(msg)
			cmds = append(cmds, cmd)
			cmds = append(cmds, m.refreshCompletion())
			return m, tea.Batch(cmds...)
		}

	case error:
		m.err = msg
		return m, nil
	}

	m.texteditor, cmd = m.texteditor.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	editorView := style.Render(m.texteditor.View())
	if m.intellisense.IsVisible() {
		return lipgloss.JoinVertical(lipgloss.Left,
			editorView,
			m.intellisense.View(),
		)
	}
	return editorView
}

// wordFromValueTail extracts the contiguous sequence of word-runes that ends
// at the tail of s. This gives the word currently being typed when the cursor
// is at the end of the text (the common case while completing).
func wordFromValueTail(s string) string {
	runes := []rune(s)
	end := len(runes)
	start := end
	for start > 0 && isWordRune(runes[start-1]) {
		start--
	}
	return string(runes[start:end])
}

func isWordRune(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}
