package inputs

import (
	"github.com/hajimehoshi/ebiten"
	"visim.muon.one/internal/view"
)

var lastX = 0
var lastY = 0

func HandleCamera(s *view.Screen) {
	x, y := ebiten.CursorPosition()
	dx, dy := lastX-x, lastY-y

	if s.HasMoved {
		s.Camera.X += dx / 3
		s.Camera.Y += dy
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		s.HasMoved = true
	} else {
		s.HasMoved = false
	}

	lastX, lastY = x, y
	s.Cursor = view.CursorPos{X: x, Y: y}
}
