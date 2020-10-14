package plots

import (
	"github.com/hajimehoshi/ebiten"
	"image/color"
	"math"
	"visim.muon.one/internal/indicators"
	"visim.muon.one/internal/stocks"
)

func Bollinger(n int, data *stocks.MarketDay, plot *ebiten.Image) {

	for x := range data.Quotes {
		quotes := data.GetQuotesInRange(x-n, x)
		q := data.GetQuote(x)
		min, _ := data.GetRange()

		if quotes == nil || q == nil {
			return
		}

		std := indicators.StandardDeviation(quotes)
		sma := indicators.SimpleMeanAverage(quotes)

		y := int((sma - min) * 100)
		ub := int((sma + 2*std - min) * 100)
		lb := int((sma - 2*std - min) * 100)

		buy := math.Min(math.Max(q.Close-(sma+std), 0)/(2*std), 1)
		sell := math.Min(math.Max((sma-std)-q.Close, 0)/(2*std), 1)

		c := color.RGBA{
			R: 48 + uint8((235-48)*buy) + uint8((106-48)*sell),
			G: 51 + uint8((77-51)*buy) + uint8((176-51)*sell),
			B: 107 + uint8((75-107)*buy) + uint8((76-107)*sell),
			A: 255,
		}

		for i := lb; i < ub; i++ {
			plot.Set(x, i, c)
		}

		plot.Set(x, y, color.RGBA{R: 126, G: 214, B: 223, A: 255})
	}

}
