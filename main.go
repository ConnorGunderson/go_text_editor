package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	text_editor "text_editor/pkg"

	"github.com/gdamore/tcell/v2"
)

const (
	MinCursorX = 7
	MinCursorY = 0
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

	scrollY := 0
	var lines []string
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

			lines = RemapScreen(mainScreen, logger, scrollY, lines, defStyle)
		case *tcell.EventKey:
			logger.Log(fmt.Sprintf("%v - %d", ev.Rune(), ev.Modifiers()))

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
				lines = []string{}
				lines = RemapScreen(mainScreen, logger, scrollY, lines, defStyle)
				cursorPosX = MinCursorX
				cursorPosY = MinCursorY
				mainScreen.ShowCursor(cursorPosX, cursorPosY)
			case tcell.KeyBackspace:
				if tcell.Key(ev.Modifiers()) == tcell.Key(tcell.ModAlt) && cursorPosX > MinCursorX {
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
				lines[scrollY] += string(ev.Rune())
				cursorPosX++
				lines = RemapScreen(mainScreen, logger, scrollY, lines, defStyle)
				logger.Log(fmt.Sprintf("%d - %v - here", scrollY, lines[scrollY]))
				mainScreen.ShowCursor(cursorPosX, cursorPosY)
			}
		case *tcell.EventMouse:
			x, y := ev.Position()
			switch ev.Buttons() {
			case tcell.ButtonPrimary:
				if x > MinCursorX {
					cursorPosX = x
					cursorPosY = y
					mainScreen.ShowCursor(cursorPosX, cursorPosY)
				}
			case tcell.WheelUp:
				if scrollY > MinCursorY {
					scrollY--
					lines = RemapScreen(mainScreen, logger, scrollY, lines, defStyle)
				}
			case tcell.WheelDown:
				_, screenHeight := mainScreen.Size()
				if scrollY < screenHeight {
					scrollY++
					lines = RemapScreen(mainScreen, logger, scrollY, lines, defStyle)
				}
			}
		}
	}
}

func CreateLogger(screen tcell.Screen, style tcell.Style) text_editor.Logger {
	logger := text_editor.Logger{Screen: screen, Style: style}

	return logger
}

func RemapScreen(t tcell.Screen, logger text_editor.Logger, cursorY int, lines []string, style tcell.Style) []string {
	_, screenHeight := t.Size()

	logger.Reset()

	for i := 0; i < screenHeight-1; i++ {

		currentLineIndex := cursorY + i

		whiteSpace := strings.Repeat(" ", 5-len(strconv.Itoa(currentLineIndex)))

		if len(lines) <= currentLineIndex {
			lines = append(lines, "")
		}

		lineStr := fmt.Sprintf("%s%d >%s", whiteSpace, currentLineIndex, lines[currentLineIndex])

		for ri, r := range lineStr {
			t.SetContent(ri, i, r, nil, style)
		}
	}

	return lines
}
