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
				s.Camera.XF += float64(dx) / s.Camera.ScaleXF
				s.Camera.Y -= dy
			if s.Camera.Y < 0 {
				s.Camera.Y = 0
			} else if s.Camera.Y > 1000 {
				s.Camera.Y = 1000
			}
			sx := s.Camera.ScaleXF
			s.Camera.ScaleXF = (float64(s.Camera.Y) / 50.0) + 1
			if sx != s.Camera.ScaleXF {
				w1 := float64(s.Window.W) / s.Camera.ScaleXF
				w2 := float64(s.Window.W) / sx
				dw := w2 - w1
				s.Camera.XF += dw/2
			}
			s.Camera.X = int(s.Camera.XF)
			s.Camera.ScaleX = int(sx)
		}
		s.HasMoved = true
	} else {
		s.HasMoved = false
	}

	lastX, lastY = x, y
	s.Cursor = view.CursorPos{X: x, Y: y}
}
