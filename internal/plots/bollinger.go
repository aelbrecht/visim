package plots

import (
	"image"
	"image/color"
	"math"
	"visim.muon.one/internal/indicators"
	"visim.muon.one/internal/stocks"
	"visim.muon.one/internal/view"
)

func Bollinger(n int, quotes []stocks.Quote, plot *image.RGBA, screen *view.Screen) {

	for x := screen.Camera.X; x < screen.Camera.X+screen.Window.W/3; x++ {

		if x < n || x >= len(quotes) {
			continue
		}

		q := quotes[x]

		std := indicators.StandardDeviation(quotes[x-n : x])
		sma := indicators.SimpleMeanAverage(quotes[x-n : x])

		y := (sma - screen.Camera.Bottom) * screen.Camera.ScaleY
		ub := (sma + 2*std - screen.Camera.Bottom) * screen.Camera.ScaleY
		lb := (sma - 2*std - screen.Camera.Bottom) * screen.Camera.ScaleY

		buy := math.Min(math.Max(q.Close-(sma+std), 0)/(2*std), 1)
		sell := math.Min(math.Max((sma-std)-q.Close, 0)/(2*std), 1)

		c := color.RGBA{
			R: 48 + uint8((235-48)*buy) + uint8((106-48)*sell),
			G: 51 + uint8((77-51)*buy) + uint8((176-51)*sell),
			B: 107 + uint8((75-107)*buy) + uint8((76-107)*sell),
			A: 255,
		}

		for i := lb; i < ub; i++ {
			for j := 0; j < 3; j++ {
				plot.Set((x-screen.Camera.X)*3+j, int(i), c)
			}
		}

		plot.Set((x-screen.Camera.X)*3+1, int(y), color.RGBA{126, 214, 223, 255})
	}

}
