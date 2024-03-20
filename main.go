package main

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// https://github.com/charmbracelet/bubbletea/blob/master/examples/split-editors/main.go

func main() {

}

type keymap = struct {
	next, prev, add, remove, quit key.Binding
}

type model struct {
	width  int
	height int
	keymap keymap
	inputs []textarea.Model
	focus  int
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func newTextarea() textarea.Model {
	t := textarea.New()
	t.Prompt = ""
	t.Placeholder = "Type something"
	t.ShowLineNumbers = true
	// t.Cursor.Style = cursorStyle
	// t.FocusedStyle.Placeholder = focusedPlaceholderStyle
	// t.BlurredStyle.Placeholder = placeholderStyle
	// t.FocusedStyle.CursorLine = cursorLineStyle
	// t.FocusedStyle.Base = focusedBorderStyle
	// t.BlurredStyle.Base = blurredBorderStyle
	// t.FocusedStyle.EndOfBuffer = endOfBufferStyle
	// t.BlurredStyle.EndOfBuffer = endOfBufferStyle
	t.KeyMap.DeleteWordBackward.SetEnabled(false)
	t.KeyMap.LineNext = key.NewBinding(key.WithKeys("down"))
	t.KeyMap.LinePrevious = key.NewBinding(key.WithKeys("up"))
	t.Blur()
	return t
}


func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd


	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.quit):
			for i := range m.inputs {
				m.inputs[i].Blur()
			}

			return m, tea.Quit
		case key.Matches(msg, m.keymap.next):
			m.inputs[m.focus].Blur()
			m.focus++
			if m.focus > len(m.inputs) -1 {
				m.focus = 0
			}

			cmd := m.inputs[m.focus].Focus()
			cmds = append(cmds, cmd)
		case key.Matches(msg, m.keymap.prev):
			m.inputs[m.focus].Blur()

			m.focus --

			if m.focus < 0 {
				m.focus = len(m.inputs)-1
			}

			cmd := m.inputs[m.focus].Focus()

			cmds = append(cmds, cmd)
		case key.Matches(msg, m.keymap.add):
			m.inputs = append(m.inputs, newTextarea())
		case key.Matches(msg, m.keymap.remove):
			m.inputs = m.inputs[:len(m.inputs)-1]
			if m.focus > len(m.inputs) - 1 {
				m.focus = len(m.inputs) - 1
			}
	}
case tea.WindowSizeMsg:
	m.height = msg.Height
	m.width = msg.Width
}


func newTextArea(tea) textarea.model {

}
