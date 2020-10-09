package view

import (
	"math"
	"visim.muon.one/internal/stocks"
)

type Camera struct {
	X, Y    int
	ScaleX  int
	ScaleXF float64
	ScaleY  float64
	Top     float64
	Bottom  float64
}

type Window struct {
	W, H int
}

type CursorPos struct {
	X, Y int
}

type Screen struct {
	Cursor   CursorPos
	Camera   *Camera
	Window   Window
	HasMoved bool
}

func (screen *Screen) AutoYAxis(quotes []stocks.Quote) {

	min := math.MaxFloat64
	max := 0.0
	for x := screen.Camera.X; x < screen.Camera.X+screen.Window.W/int(screen.Camera.ScaleX); x++ {
		if x < 0 || x >= len(quotes) {
			continue
		}
		q := quotes[x]
		if q.Low < min {
			min = q.Low
		}
		if q.High > max {
			max = q.High
		}
	}

	// add some padding
	minMaxDelta := max - min
	min -= minMaxDelta / 2
	max += minMaxDelta / 10
	minMaxDelta = max - min

	if minMaxDelta < 0.01 {
		minMaxDelta = 0.01
	}
	screen.Camera.ScaleY = float64(screen.Window.H) / minMaxDelta
	screen.Camera.Top = max
	screen.Camera.Bottom = min
}
