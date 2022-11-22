package texteditor

import "github.com/charmbracelet/lipgloss"

// Style that will be applied to the text area.
//
// Style can be applied to focused and unfocused states to change the styles
// depending on the focus state.
//
// For an introduction to styling with Lip Gloss see:
// https://github.com/charmbracelet/lipgloss
type Style struct {
	Base          lipgloss.Style
	Cursor        lipgloss.Style
	LineDecorator lipgloss.Style
	Text          lipgloss.Style
}

// DefaultStyles returns the default styles for focused and blurred states for
// the textarea.
func DefaultStyles() (Style, Style) {
	lineDecoratorBase := lipgloss.NewStyle().
		Align(lipgloss.Right).
		PaddingRight(1)

	focused := Style{
		Base:          lipgloss.NewStyle(),
		Cursor:        lipgloss.NewStyle().Background(lipgloss.Color("1")).Bold(true),
		LineDecorator: lineDecoratorBase.Copy().Foreground(lipgloss.AdaptiveColor{Light: "249", Dark: "8"}),
		Text:          lipgloss.NewStyle(),
	}
	blurred := Style{
		Base:          lipgloss.NewStyle(),
		Cursor:        lipgloss.NewStyle(),
		LineDecorator: lineDecoratorBase.Copy().Foreground(lipgloss.AdaptiveColor{Light: "249", Dark: "8"}),
		Text:          lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "245", Dark: "7"}),
	}

	return focused, blurred
}
