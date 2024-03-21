package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
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
	mainScreen.SetStyle(defStyle)
	mainScreen.EnableMouse()
	mainScreen.EnablePaste()
	mainScreen.Clear()

	quit := func() {
		// You have to catch panics in a defer, clean up, and
		// re-raise them - otherwise your application can
		// die without leaving any diagnostic trace.
		maybePanic := recover()
		mainScreen.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit()

	const minCursorX = 7

	cursorIndexX := minCursorX
	cursorIndexY := 0
	mainScreen.ShowCursor(cursorIndexX, cursorIndexY)
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
			_, screenHeight := mainScreen.Size()
			for i := 0; i < screenHeight; i++ {
				whiteSpace := strings.Repeat(" ", 5-len(strconv.Itoa(i)))

				screenIndexStr := fmt.Sprintf("%s%d >", whiteSpace, i)

				for ri, r := range screenIndexStr {
					mainScreen.SetContent(ri, i, r, nil, defStyle)
				}
			}

		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEscape, tcell.KeyCtrlC:
				return
			case tcell.KeyCtrlL:
				mainScreen.Sync()
			case tcell.KeyUp:
				if cursorIndexY != 0 {
					cursorIndexY--
					mainScreen.ShowCursor(cursorIndexX, cursorIndexY)
				}
			case tcell.KeyDown:
				_, screenHeight := mainScreen.Size()
				if cursorIndexY < screenHeight {
					cursorIndexY++
					mainScreen.ShowCursor(cursorIndexX, cursorIndexY)
				}
			case tcell.KeyLeft:
				if cursorIndexX > minCursorX {
					cursorIndexX--
					mainScreen.ShowCursor(cursorIndexX, cursorIndexY)
				}
			case tcell.KeyRight:
				screenWidth, _ := mainScreen.Size()
				if cursorIndexX < screenWidth {
					cursorIndexX++
					mainScreen.ShowCursor(cursorIndexX, cursorIndexY)
				}
			default:
				mainScreen.SetContent(cursorIndexX, cursorIndexY, ev.Rune(), nil, defStyle)
				cursorIndexX++
				mainScreen.ShowCursor(cursorIndexX, cursorIndexY)
			}
		case *tcell.EventMouse:
			// x, y := ev.Position()
		}
	}
}
