package plots

import (
	"image"
	"image/color"
	"visim.muon.one/internal/stocks"
	"visim.muon.one/internal/view"
)

func Candles(quotes []stocks.Quote, plot *image.RGBA, screen *view.Screen) {

	for x := screen.Camera.X; x < screen.Camera.X+screen.Window.W/3; x++ {

		if x < 0 || x >= len(quotes) {
			continue
		}

		q := quotes[x]

		lb := int((q.Low - screen.Camera.Bottom) * screen.Camera.ScaleY)
		ub := int((q.High - screen.Camera.Bottom) * screen.Camera.ScaleY)

		c := color.RGBA{235, 77, 75, 255}
		if q.Open <= q.Close {
			c = color.RGBA{106, 176, 76, 255}
		}

		for y := lb; y < ub; y++ {
			plot.Set((x-screen.Camera.X)*3+1, y, c)
		}
	}

}
