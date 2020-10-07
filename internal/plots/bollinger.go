package plots

import (
	"image"
	"image/color"
	"visim.muon.one/internal/indicators"
	"visim.muon.one/internal/stocks"
	"visim.muon.one/internal/view"
)

func Bollinger(n int, quotes []stocks.Quote, plot *image.RGBA, screen *view.Screen) {

	for x := screen.Camera.X; x < screen.Camera.X+screen.Window.W/3; x++ {

		if x < n || x >= len(quotes) {
			continue
		}

		std := indicators.StandardDeviation(quotes[x-n : x])
		sma := indicators.SimpleMeanAverage(quotes[x-n : x])

		y := (sma - screen.Camera.Bottom) * screen.Camera.ScaleY
		ub := (sma + 2*std - screen.Camera.Bottom) * screen.Camera.ScaleY
		lb := (sma - 2*std - screen.Camera.Bottom) * screen.Camera.ScaleY

		for i := lb; i < ub; i++ {
			plot.Set((x-screen.Camera.X)*3-1, int(i), color.RGBA{34, 166, 179, 40})
			plot.Set((x-screen.Camera.X)*3, int(i), color.RGBA{34, 166, 179, 40})
			plot.Set((x-screen.Camera.X)*3+1, int(i), color.RGBA{34, 166, 179, 40})
		}

		plot.Set((x-screen.Camera.X)*3+1, int(y), color.RGBA{126, 214, 223, 255})
	}

}
