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

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if s.HasMoved {
			if !ebiten.IsKeyPressed(ebiten.KeyShift) {
				s.Camera.XF += float64(dx) / s.Camera.ScaleXF
			} else {
				s.Camera.Y -= dy
			}
			s.Camera.X = int(s.Camera.XF)
			if s.Camera.Y < 0 {
				s.Camera.Y = 0
			} else if s.Camera.Y > 500 {
				s.Camera.Y = 500
			}
			sx := s.Camera.ScaleXF
			s.Camera.ScaleXF = (float64(s.Camera.Y) / 50.0) + 3
			if sx != s.Camera.ScaleXF {
				w1 := float64(s.Window.W) / s.Camera.ScaleXF
				w2 := float64(s.Window.W) / sx
				dw := w2 - w1
				s.Camera.X -= int(dw / 2)
			}
			s.Camera.ScaleX = int(sx)
		}
		s.HasMoved = true
	} else {
		s.HasMoved = false
	}

	lastX, lastY = x, y
	s.Cursor = view.CursorPos{X: x, Y: y}
}
