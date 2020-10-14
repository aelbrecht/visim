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

const MinutesInDay = 6*60 + 30

func GetDay(x int) int {
	return x / (MinutesInDay)
}

func GetQuoteIndex(x int) int {
	return x % MinutesInDay
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

func (screen *Screen) AutoYAxis(m *stocks.Model) {

	min := math.MaxFloat64
	max := 0.0

	for i := range m.Data {
		m1, m2 := m.Data[i].GetRange()
		if m1 < min {
			min = m1
		}
		if m2 > max {
			max = m2
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
