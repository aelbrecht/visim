package inputs

import (
	"github.com/hajimehoshi/ebiten"
	"visim.muon.one/internal/view"
)

var (
	Pressed1 = false
	Pressed2 = false
	Pressed3 = false
	Pressed4 = false
)

func handlePlotCamera(dx int, dy int, s *view.Screen) {
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
		w1 := float64(s.Plot.W) / s.Camera.ScaleXF
		w2 := float64(s.Plot.W) / sx
		dw := w2 - w1
		s.Camera.XF += dw / 2
	}
	s.Camera.X = int(s.Camera.XF)
	s.Camera.ScaleX = int(sx)
}

func HandlePlot(options *Options) bool {

	update := false

	if ebiten.IsKeyPressed(ebiten.KeyF1) {
		Pressed1 = true
	} else {
		if Pressed1 {
			Pressed1 = false
			options.ShowQuotes = !options.ShowQuotes
			update = true
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyF2) {
		Pressed2 = true
	} else {
		if Pressed2 {
			Pressed2 = false
			options.ShowBollinger = !options.ShowBollinger
			update = true
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyF3) {
		Pressed3 = true
	} else {
		if Pressed3 {
			Pressed3 = false
			options.ShowRSI = !options.ShowRSI
			update = true
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyF4) {
		Pressed4 = true
	} else {
		if Pressed4 {
			Pressed4 = false
			options.ShowSupportResistance = !options.ShowSupportResistance
			update = true
		}
	}

	return update
}

type Options struct {
	ShowBollinger         bool
	ShowRSI               bool
	ShowQuotes            bool
	ShowSupportResistance bool
}
