package inputs

import (
	"github.com/hajimehoshi/ebiten"
	"visim.muon.one/internal/view"
)

var lastX = 0
var lastY = 0

func HandleCamera(s *view.Screen) {
	x, y := ebiten.CursorPosition()
	dx, _ := lastX-x, lastY-y

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if s.HasMoved {
			s.Camera.X += dx
			//s.Camera.Y -= dy
			if s.Camera.Y < 0 {
				//s.Camera.Y = 0
			} else if s.Camera.Y > 500 {
				//s.Camera.Y = 500
			}
			s.Camera.ScaleXF = (float64(s.Camera.Y) / 50.0) + 3
			sx := int(s.Camera.ScaleXF)
			if sx != s.Camera.ScaleX {
				w1 := s.Window.W / s.Camera.ScaleX
				w2 := s.Window.W / sx
				dw := w2 - w1
				s.Camera.X -= dw / 2
			}
			s.Camera.ScaleX = sx
		}
		s.HasMoved = true
	} else {
		s.HasMoved = false
	}

	lastX, lastY = x, y
	s.Cursor = view.CursorPos{X: x, Y: y}
}
