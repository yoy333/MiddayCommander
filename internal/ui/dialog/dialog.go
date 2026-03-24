package dialog

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/kooler/MiddayCommander/internal/ui/theme"
)

// Kind identifies the dialog type.
type Kind int

const (
	KindConfirm Kind = iota
	KindInput
	KindProgress
	KindError
)

// Result is sent when a dialog closes.
type Result struct {
	Kind      Kind
	Confirmed bool   // for confirm dialogs
	Text      string // for input dialogs
	Tag       string // caller-defined tag to identify which operation triggered the dialog
}

// Model represents a modal dialog overlaid on the panels.
type Model struct {
	kind    Kind
	title   string
	message string
	tag     string // passed back in Result

	// Input dialog
	input    string
	inputPos int

	// Progress dialog
	progress float64
	current  string

	// State
	done   bool
	result Result

	width int
}

// NewConfirm creates a Yes/No confirmation dialog.
func NewConfirm(title, message, tag string) Model {
	return Model{
		kind:    KindConfirm,
		title:   title,
		message: message,
		tag:     tag,
		width:   50,
	}
}

// NewInput creates a text input dialog.
func NewInput(title, message, defaultValue, tag string) Model {
	return Model{
		kind:     KindInput,
		title:    title,
		message:  message,
		tag:      tag,
		input:    defaultValue,
		inputPos: len(defaultValue),
		width:    50,
	}
}

// NewError creates an error display dialog.
func NewError(title, message string) Model {
	return Model{
		kind:    KindError,
		title:   title,
		message: message,
		width:   50,
	}
}

// NewProgress creates a progress dialog.
func NewProgress(title, tag string) Model {
	return Model{
		kind:  KindProgress,
		title: title,
		tag:   tag,
		width: 50,
	}
}

// Done returns true when the dialog has been dismissed.
func (m Model) Done() bool {
	return m.done
}

// GetResult returns the dialog result.
func (m Model) GetResult() Result {
	return m.result
}

// SetProgress updates the progress dialog state.
func (m *Model) SetProgress(progress float64, current string) {
	m.progress = progress
	m.current = current
}

// Update handles key events for the dialog.
func (m *Model) Update(msg tea.KeyMsg) tea.Cmd {
	switch m.kind {
	case KindConfirm:
		return m.updateConfirm(msg)
	case KindInput:
		return m.updateInput(msg)
	case KindError:
		return m.updateError(msg)
	case KindProgress:
		// Progress dialogs can't be dismissed by keyboard (they close when the operation ends)
		return nil
	}
	return nil
}

func (m *Model) updateConfirm(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "y", "Y", "enter":
		m.done = true
		m.result = Result{Kind: KindConfirm, Confirmed: true, Tag: m.tag}
	case "n", "N", "esc":
		m.done = true
		m.result = Result{Kind: KindConfirm, Confirmed: false, Tag: m.tag}
	}
	return nil
}

func (m *Model) updateInput(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "enter":
		m.done = true
		m.result = Result{Kind: KindInput, Confirmed: true, Text: m.input, Tag: m.tag}
	case "esc":
		m.done = true
		m.result = Result{Kind: KindInput, Confirmed: false, Tag: m.tag}
	case "backspace":
		if m.inputPos > 0 {
			m.input = m.input[:m.inputPos-1] + m.input[m.inputPos:]
			m.inputPos--
		}
	case "delete":
		if m.inputPos < len(m.input) {
			m.input = m.input[:m.inputPos] + m.input[m.inputPos+1:]
		}
	case "left":
		if m.inputPos > 0 {
			m.inputPos--
		}
	case "right":
		if m.inputPos < len(m.input) {
			m.inputPos++
		}
	case "home":
		m.inputPos = 0
	case "end":
		m.inputPos = len(m.input)
	default:
		if len(msg.String()) == 1 && msg.String()[0] >= 32 {
			m.input = m.input[:m.inputPos] + msg.String() + m.input[m.inputPos:]
			m.inputPos++
		}
	}
	return nil
}

func (m *Model) updateError(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "enter", "esc", "q":
		m.done = true
		m.result = Result{Kind: KindError}
	}
	return nil
}

// Close dismisses the dialog (used for progress dialogs when the operation ends).
func (m *Model) Close() {
	m.done = true
	m.result = Result{Kind: m.kind, Tag: m.tag}
}

// View renders the dialog as a centered box.
func (m Model) View(th theme.Theme, screenWidth, screenHeight int) string {
	boxWidth := m.width
	if boxWidth > screenWidth-4 {
		boxWidth = screenWidth - 4
	}
	innerWidth := boxWidth - 4 // padding

	borderStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")).
		Background(lipgloss.Color("0"))
	contentStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")).
		Background(lipgloss.Color("0"))
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("14")).
		Background(lipgloss.Color("0")).
		Bold(true)

	var lines []string

	// Top border with title
	titleStr := " " + m.title + " "
	padLen := innerWidth - len(titleStr)
	if padLen < 0 {
		padLen = 0
	}
	topBorder := "┌" + titleStr + strings.Repeat("─", padLen) + "┐"
	lines = append(lines, borderStyle.Render(topBorder))

	// Empty line
	lines = append(lines, borderStyle.Render("│")+contentStyle.Render(strings.Repeat(" ", innerWidth+2))+borderStyle.Render("│"))

	// Message
	for _, msgLine := range wrapText(m.message, innerWidth) {
		padded := " " + padRight(msgLine, innerWidth) + " "
		lines = append(lines, borderStyle.Render("│")+contentStyle.Render(padded)+borderStyle.Render("│"))
	}

	// Empty line
	lines = append(lines, borderStyle.Render("│")+contentStyle.Render(strings.Repeat(" ", innerWidth+2))+borderStyle.Render("│"))

	// Kind-specific content
	switch m.kind {
	case KindConfirm:
		hint := " [Y]es  [N]o "
		padded := " " + padRight(hint, innerWidth) + " "
		lines = append(lines, borderStyle.Render("│")+titleStyle.Render(padded)+borderStyle.Render("│"))

	case KindInput:
		// Input field
		inputDisplay := m.input
		if len(inputDisplay) > innerWidth-2 {
			inputDisplay = inputDisplay[len(inputDisplay)-innerWidth+2:]
		}
		inputLine := " [" + padRight(inputDisplay, innerWidth-2) + "]"
		if len(inputLine) > innerWidth+2 {
			inputLine = inputLine[:innerWidth+2]
		}
		inputStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("11")).
			Background(lipgloss.Color("0"))
		lines = append(lines, borderStyle.Render("│")+inputStyle.Render(inputLine)+borderStyle.Render("│"))

	case KindProgress:
		// Progress bar
		barWidth := innerWidth - 2
		filled := int(m.progress * float64(barWidth))
		if filled > barWidth {
			filled = barWidth
		}
		bar := " " + strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled) + " "
		progressStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("10")).
			Background(lipgloss.Color("0"))
		lines = append(lines, borderStyle.Render("│")+progressStyle.Render(bar)+borderStyle.Render("│"))

		if m.current != "" {
			curLine := " " + padRight(m.current, innerWidth) + " "
			lines = append(lines, borderStyle.Render("│")+contentStyle.Render(curLine)+borderStyle.Render("│"))
		}

	case KindError:
		hint := " Press Enter to close "
		padded := " " + padRight(hint, innerWidth) + " "
		lines = append(lines, borderStyle.Render("│")+titleStyle.Render(padded)+borderStyle.Render("│"))
	}

	// Bottom border
	lines = append(lines, borderStyle.Render("└"+strings.Repeat("─", innerWidth+2)+"┘"))

	box := strings.Join(lines, "\n")

	// Center the box on screen
	return lipgloss.Place(screenWidth, screenHeight, lipgloss.Center, lipgloss.Center, box)
}

func padRight(s string, width int) string {
	if len(s) >= width {
		return s[:width]
	}
	return s + strings.Repeat(" ", width-len(s))
}

func wrapText(text string, width int) []string {
	if len(text) <= width {
		return []string{text}
	}
	var lines []string
	for len(text) > width {
		// Find last space before width
		cut := width
		for cut > 0 && text[cut] != ' ' {
			cut--
		}
		if cut == 0 {
			cut = width
		}
		lines = append(lines, text[:cut])
		text = strings.TrimLeft(text[cut:], " ")
	}
	if text != "" {
		lines = append(lines, text)
	}
	return lines
}
