package view

import (
	"math"
	"visim.muon.one/internal/stocks"
)

type Camera struct {
	X, Y     int
	XF       float64
	ScaleX   int
	ScaleXF  float64
	ScaleY   float64
	Top      float64
	Bottom   float64
	GridSize int
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
	Plot     Window
	Program  Window
}

func (screen *Screen) VisibleDays() (int, int) {
	x0 := screen.Camera.X
	x1 := screen.Camera.X + screen.Plot.W/screen.Camera.ScaleX
	return stocks.GetDay(x0), stocks.GetDay(x1)
}

func (screen *Screen) AutoYAxis(m *stocks.Model) {

	x0 := screen.Camera.X
	x1 := screen.Camera.X + int(float64(screen.Plot.W)/screen.Camera.ScaleXF)

	min := math.MaxFloat64
	max := 0.0

	for x := x0; x < x1; x++ {
		q := m.GetQuote(x)
		if q == nil {
			continue
		}
		if q.Low < min {
			min = q.Low
		}
		if q.High > max {
			max = q.High
		}
	}

	// add some padding
	minMaxDelta := max - min
	min -= minMaxDelta / 10
	max += minMaxDelta / 10
	minMaxDelta = max - min

	if min > 999999999 {
		min = 0
		max = 10
	}

	for minMaxDelta < 2 {
		min -= 0.1
		max += 0.1
		minMaxDelta = max - min
	}

	screen.Camera.ScaleY = float64(screen.Plot.H) / minMaxDelta
	screen.Camera.Top = max
	screen.Camera.Bottom = min
}
