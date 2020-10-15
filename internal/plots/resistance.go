package plots

import (
	"github.com/hajimehoshi/ebiten"
	"image/color"
	"visim.muon.one/internal/indicators"
	"visim.muon.one/internal/stocks"
)

var supportLine *ebiten.Image
var resistanceLine *ebiten.Image

func init() {
	supportLine, _ = ebiten.NewImage(1, 1, ebiten.FilterDefault)
	supportLine.Fill(color.RGBA{G: 255, A: 100})
	resistanceLine, _ = ebiten.NewImage(1, 1, ebiten.FilterDefault)
	resistanceLine.Fill(color.RGBA{R: 255, A: 100})
}

func SR(n int, m *stocks.MarketDay, plot *ebiten.Image) {

	min, _ := m.GetRange()

	for i := range m.Quotes {
		quotes1 := m.GetQuotesInRange(i-n-2, i-2)
		quotes2 := m.GetQuotesInRange(i-n-1, i-1)
		quotes3 := m.GetQuotesInRange(i-n, i)

		if quotes1 == nil || quotes2 == nil || quotes3 == nil {
			continue
		}

		avg1 := indicators.SimpleMeanAverage(quotes1)
		avg2 := indicators.SimpleMeanAverage(quotes2)
		avg3 := indicators.SimpleMeanAverage(quotes3)

		if avg1 < avg2 && avg3 < avg2 {
			qs := m.GetQuotesInRange(i-n/2-2-5, i-n/2-2+5)
			high := 0.0
			for i2 := range qs {
				if qs[i2].High > high {
					high = qs[i2].High
				}
			}
			y := (high - min) * 100
			op := ebiten.DrawImageOptions{}
			op.GeoM.Scale(60, 1)
			op.GeoM.Translate(float64(i-n/2-2), y)
			plot.DrawImage(resistanceLine, &op)
		}

		if avg1 > avg2 && avg3 > avg2 {
			qs := m.GetQuotesInRange(i-n/2-2-5, i-n/2-2+5)
			low := 99999.0
			for i2 := range qs {
				if qs[i2].Low < low {
					low = qs[i2].Low
				}
			}
			y := (low - min) * 100
			op := ebiten.DrawImageOptions{}
			op.GeoM.Scale(20, 1)
			op.GeoM.Translate(float64(i-n/2-2), y)
			plot.DrawImage(supportLine, &op)
		}

		//y := (avg1 - screen.Camera.Bottom) * screen.Camera.ScaleY
		// SetDash(i-n/2-2, int(y), 6, color.RGBA{200, 255, 223, 150}, plot, screen)

	}
}
