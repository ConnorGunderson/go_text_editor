package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
)

const (
	MinCursorX = 7
	MinCursorY = 0
	LogMinX    = 6
)

func main() {
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)

	// Initialize screen
	mainScreen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := mainScreen.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	logger := CreateLogger(mainScreen, defStyle)

	mainScreen.SetStyle(defStyle)
	mainScreen.EnableMouse()
	mainScreen.EnablePaste()
	mainScreen.Clear()

	quit := func() {
		maybePanic := recover()
		mainScreen.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit()

	cursorPosX := MinCursorX
	cursorPosY := MinCursorY
	mainScreen.ShowCursor(cursorPosX, cursorPosY)
	// Event loop
	for {
		// Update screen
		mainScreen.Show()

		// Poll event
		ev := mainScreen.PollEvent()

		// Process event
		switch ev := ev.(type) {
		case *tcell.EventResize:
			mainScreen.Sync()
			RemapScreen(mainScreen, logger, defStyle)
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEscape, tcell.KeyCtrlC:
				return
			case tcell.KeyCtrlL:
				mainScreen.Sync()
			case tcell.KeyUp:
				if cursorPosY != MinCursorY {
					cursorPosY--
					mainScreen.ShowCursor(cursorPosX, cursorPosY)
				}
			case tcell.KeyDown:
				_, screenHeight := mainScreen.Size()
				if cursorPosY < screenHeight {
					cursorPosY++
					mainScreen.ShowCursor(cursorPosX, cursorPosY)
				}
			case tcell.KeyLeft:
				if cursorPosX > MinCursorX {
					cursorPosX--
					mainScreen.ShowCursor(cursorPosX, cursorPosY)
				}
			case tcell.KeyRight:
				screenWidth, _ := mainScreen.Size()
				if cursorPosX < screenWidth {
					cursorPosX++
					mainScreen.ShowCursor(cursorPosX, cursorPosY)
				}
			case tcell.KeyCtrlR:
				mainScreen.Clear()
				RemapScreen(mainScreen, logger, defStyle)
				cursorPosX = MinCursorX
				cursorPosY = MinCursorY
				mainScreen.ShowCursor(cursorPosX, cursorPosY)
			case tcell.KeyBackspace:
				if tcell.Key(ev.Modifiers()) == tcell.Key(tcell.ModCtrl) && cursorPosX > MinCursorX {
					hasChar := false
					for i := cursorPosX; i > MinCursorX; i-- {
						r, _, _, _ := mainScreen.GetContent(cursorPosX, cursorPosY)
						if r != ' ' {
							hasChar = true
						}
						if r == ' ' && hasChar {
							break
						}
						mainScreen.SetContent(cursorPosX, cursorPosY, ' ', nil, defStyle)
						cursorPosX--
						mainScreen.ShowCursor(cursorPosX, cursorPosY)
					}
				}

				if cursorPosX == MinCursorX {
					mainScreen.SetContent(cursorPosX, cursorPosY, ' ', nil, defStyle)
				}

				if cursorPosX > MinCursorX {
					mainScreen.SetContent(cursorPosX, cursorPosY, ' ', nil, defStyle)
					cursorPosX--
					mainScreen.ShowCursor(cursorPosX, cursorPosY)
				}
			case tcell.KeyRune:
				mainScreen.SetContent(cursorPosX, cursorPosY, ev.Rune(), nil, defStyle)
				cursorPosX++
				mainScreen.ShowCursor(cursorPosX, cursorPosY)
			}
		case *tcell.EventMouse:
			x, y := ev.Position()

			if ev.Buttons() == tcell.ButtonPrimary {
				if x > MinCursorX {
					cursorPosX = x
					cursorPosY = y
					mainScreen.ShowCursor(cursorPosX, cursorPosY)
				}
			}

		}
	}
}

type Logger struct {
	screen tcell.Screen
	style  tcell.Style
}

func (l Logger) Log(str string) {
	screenWidth, screenHeight := l.screen.Size()
	for i := range screenWidth {
		if i > LogMinX {
			l.screen.SetContent(i, screenHeight, ' ', nil, l.style)
		}
	}
	for i, r := range str {
		l.screen.SetContent(LogMinX+i, screenHeight, r, nil, l.style)
	}

}

func (l Logger) Reset() {
	logStr := "Log: "

	_, screenHeight := l.screen.Size()

	for i, r := range logStr {
		l.screen.SetContent(i, screenHeight-1, r, nil, l.style)
	}
}

func CreateLogger(screen tcell.Screen, style tcell.Style) Logger {
	logger := Logger{
		screen,
		style,
	}

	return logger
}

func RemapScreen(t tcell.Screen, logger Logger, style tcell.Style) {
	_, screenHeight := t.Size()

	logger.Reset()

	for i := 0; i < screenHeight-1; i++ {

		whiteSpace := strings.Repeat(" ", 5-len(strconv.Itoa(i)))

		screenIndexStr := fmt.Sprintf("%s%d >", whiteSpace, i)

		for ri, r := range screenIndexStr {
			t.SetContent(ri, i, r, nil, style)
		}
	}
}
