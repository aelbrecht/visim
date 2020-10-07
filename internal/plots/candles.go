package plots

import (
	"image"
	"image/color"
	"math"
	"visim.muon.one/internal/stocks"
	"visim.muon.one/internal/view"
)

func PlotCandles(quotes []stocks.Quote, plot *image.RGBA, screen *view.Screen) {

	min := math.MaxFloat64
	max := 0.0
	for x := screen.Camera.X; x < screen.Camera.X+screen.Window.W; x++ {
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
	min -= minMaxDelta / 10
	max += minMaxDelta / 10
	minMaxDelta = max - min

	scale := 800.0 / minMaxDelta

	for x := screen.Camera.X; x < screen.Camera.X+screen.Window.W/3; x++ {

		if x < 0 || x >= len(quotes) {
			continue
		}

		q := quotes[x]

		lb := int((q.Low - min) * scale)
		ub := int((q.High - min) * scale)

		c := color.RGBA{235, 77, 75, 255}
		if q.Open > q.Close {
			c = color.RGBA{106, 176, 76, 255}
		}

		for y := lb; y < ub; y++ {
			plot.Set((x-screen.Camera.X)*3+1, y, c)
		}
	}

}
