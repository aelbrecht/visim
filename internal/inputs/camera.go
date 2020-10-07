package inputs

import (
	"github.com/hajimehoshi/ebiten"
	"visim.muon.one/internal/view"
)

var lastX = 0
var lastY = 0
var moving = false

func HandleCamera(s *view.Screen) {
	x, y := ebiten.CursorPosition()
	dx, dy := lastX-x, lastY-y

	if moving {
		s.Camera.X += dx/3
		s.Camera.Y += dy
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		moving = true
	} else {
		moving = false
	}

	lastX, lastY = x, y
}
