package plots

import (
	"image"
	"image/color"
	"visim.muon.one/internal/indicators"
	"visim.muon.one/internal/stocks"
	"visim.muon.one/internal/view"
)

func Resistance(n int, m *stocks.Model, plot *image.RGBA, screen *view.Screen) {

	c1 := color.RGBA{255, 0, 0, 100}
	c2 := color.RGBA{0, 255, 0, 100}

	PlotX(func(i int) {
		quotes1 := m.GetQuoteRange(i-n-2, i-2)
		quotes2 := m.GetQuoteRange(i-n-1, i-1)
		quotes3 := m.GetQuoteRange(i-n, i)

		if quotes1 == nil || quotes2 == nil {
			return
		}

		avg1 := indicators.SimpleMeanAverage(quotes1)
		avg2 := indicators.SimpleMeanAverage(quotes2)
		avg3 := indicators.SimpleMeanAverage(quotes3)

		if avg1 < avg2 && avg3 < avg2 {
			SetPixel(i-n/2-1, 100, c2, plot, screen)

			qs := m.GetQuoteRange(i-n/2-2-5, i-n/2-2+5)
			high := 0.0
			for i2 := range qs {
				if qs[i2].High > high {
					high = qs[i2].High
				}
			}
			y := (high - screen.Camera.Bottom) * screen.Camera.ScaleY
			for z := 0; z < 60; z++ {
				SetDash(i-n/2-2+z, int(y), 5, c1, plot, screen)
			}
		}

		if avg1 > avg2 && avg3 > avg2 {
			SetPixel(i-n/2-1, 100, c1, plot, screen)
			qs := m.GetQuoteRange(i-n/2-2-5, i-n/2-2+5)
			low := 99999.0
			for i2 := range qs {
				if qs[i2].Low < low {
					low = qs[i2].Low
				}
			}
			y := (low - screen.Camera.Bottom) * screen.Camera.ScaleY
			for z := 0; z < 60; z++ {
				SetDash(i-n/2-2+z, int(y), 5, c2, plot, screen)
			}
		}

		//y := (avg1 - screen.Camera.Bottom) * screen.Camera.ScaleY
		// SetDash(i-n/2-2, int(y), 6, color.RGBA{200, 255, 223, 150}, plot, screen)

	}, screen)
}
