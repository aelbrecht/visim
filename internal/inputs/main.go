package inputs

import (
	"github.com/hajimehoshi/ebiten"
	"visim.muon.one/internal/layout"
	"visim.muon.one/internal/view"
)

var leftDown = false
var dragging = false
var canDrag = false
var lastX = 0
var lastY = 0

func inPlotBounds(x int, y int, s *view.Screen) bool {
	return y > 40 && y < s.Program.H-16
}

func HandleMouseLeft(s *view.Screen, buttons []*layout.Button) {

	x, y := ebiten.CursorPosition()
	dx, dy := lastX-x, lastY-y

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if leftDown {
			// left down
			if (dx != 0 || dy != 0) && canDrag {
				dragging = true
			}
		} else {
			// left pressed
			leftDown = true
			// check if can drag here
			if inPlotBounds(x, y, s) {
				canDrag = true
			}
		}
	} else {
		if leftDown {

			if !dragging {
				for _, button := range buttons {
					button.MouseClick()
				}
			}

			// left released
			leftDown = false
			canDrag = false
			dragging = false
		}
	}

	if dragging {
		handlePlotCamera(dx, dy, s)
	}

	for _, button := range buttons {
		if dragging {
			button.MouseMove(-1, -1)
		} else {
			button.MouseMove(x, y)
		}
	}

	lastX, lastY = x, y
	s.Cursor = view.CursorPos{X: x, Y: y}
}
