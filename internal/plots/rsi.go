package plots

import (
	"image"
	"image/color"
	"visim.muon.one/internal/indicators"
	"visim.muon.one/internal/stocks"
	"visim.muon.one/internal/view"
)

func RSI(n int, quotes []stocks.Quote, plot *image.RGBA, screen *view.Screen) {

	for x := screen.Camera.X; x < screen.Camera.X+screen.Window.W/3; x++ {

		if x < n || x >= len(quotes) {
			continue
		}

		rsi := indicators.RelativeStrengthIndex(quotes[x-n : x])

		y := rsi * 100
		for i := 0.0; i < y; i++ {
			for j := -1; j < 2; j++ {
				plot.Set((x-screen.Camera.X)*3+j, int(i), color.RGBA{190, 46, 221, 40})
			}
		}

		for j := -1; j < 2; j++ {
			plot.Set((x-screen.Camera.X)*3+j, 100, color.RGBA{190, 46, 221, 60})
		}
	}

}
