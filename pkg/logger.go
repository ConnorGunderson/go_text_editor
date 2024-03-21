package text_editor

import "github.com/gdamore/tcell/v2"

const (
	LogMinX = 6
)

type Logger struct {
	Screen tcell.Screen
	Style  tcell.Style
}

func (l Logger) Log(str string) {
	screenWidth, screenHeight := l.Screen.Size()
	for i := range screenWidth {
		if i > LogMinX {
			l.Screen.SetContent(i, screenHeight, ' ', nil, l.Style)
		}
	}
	for i, r := range str {
		l.Screen.SetContent(LogMinX+i, screenHeight-1, r, nil, l.Style)
	}

}

func (l Logger) Reset() {
	logStr := "Log: "

	screenWidth, screenHeight := l.Screen.Size()

	for i, r := range logStr {
		l.Screen.SetContent(i, screenHeight-1, r, nil, l.Style)
	}

	for i := len(logStr); i < screenWidth; i++ {
		l.Screen.SetContent(i, screenHeight-1, ' ', nil, l.Style)
	}
}
