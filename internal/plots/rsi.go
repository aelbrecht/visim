package plots

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"image/color"
	"math"
	"visim.muon.one/internal/fonts"
	"visim.muon.one/internal/indicators"
	"visim.muon.one/internal/stocks"
	"visim.muon.one/internal/view"
)

func TooltipRSI(i int, n int, m *stocks.Model, buffer *ebiten.Image, screen *view.Screen) {
	day := stocks.GetDay(i)
	quotes := m.GetQuoteDay(day).GetQuotesInRange(i-n, i)
	if quotes == nil {
		return
	}
	rsi := indicators.RelativeStrengthIndex(quotes)
	y := screen.Plot.H - int(rsi*100) + 100 + 40
	x := int(float64(i-screen.Camera.X)*screen.Camera.ScaleXF) + paddingLeft
	fonts.Background(x-3, y+3, 54, 13, color.RGBA{48, 51, 107, 200}, buffer)
	text.Draw(buffer, fmt.Sprintf("RSI: %d", int(rsi*100)), fonts.FaceNormal, x, y, color.White)
}

func RSI(n int, data *stocks.MarketDay, plot *ebiten.Image) {

	for x := range data.Quotes {

		quotes := data.GetQuotesInRange(x-n, x)
		if quotes == nil {
			continue
		}

		rsi := indicators.RelativeStrengthIndex(quotes)

		sell := math.Max(rsi-0.65, 0) / 0.35
		buy := math.Max(1-rsi-0.65, 0) / 0.35

		c := color.RGBA{
			R: 48 + uint8((235-48)*sell) + uint8((106-48)*buy),
			G: 51 + uint8((77-51)*sell) + uint8((176-51)*buy),
			B: 107 + uint8((75-107)*sell) + uint8((76-107)*buy),
			A: 255,
		}

		y := int(rsi * 100)
		for i := 0; i < y; i++ {
			plot.Set(x, i, c)
		}
	}
}
