package plots

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"image"
	"image/color"
	"math"
	"visim.muon.one/internal/fonts"
	"visim.muon.one/internal/indicators"
	"visim.muon.one/internal/stocks"
	"visim.muon.one/internal/view"
)

func TooltipRSI(i int, n int, quotes []stocks.Quote, buffer *ebiten.Image, screen *view.Screen) {
	rsi := indicators.RelativeStrengthIndex(quotes[i-n : i])
	y := screen.Window.H - int(rsi*100)
	x := (i-screen.Camera.X)*3 + paddingLeft
	fonts.Background(x-3, y+3, 54, 13, color.RGBA{48, 51, 107, 200}, buffer)
	text.Draw(buffer, fmt.Sprintf("RSI: %d", int(rsi*100)), fonts.FaceNormal, x, y, color.White)
}

func RSI(n int, quotes []stocks.Quote, plot *image.RGBA, screen *view.Screen) {

	for x := screen.Camera.X; x < screen.Camera.X+screen.Window.W/3; x++ {

		if x < n || x >= len(quotes) {
			continue
		}

		rsi := indicators.RelativeStrengthIndex(quotes[x-n : x])

		buy := math.Max(rsi-0.65, 0) / 0.35
		sell := math.Max(1-rsi-0.65, 0) / 0.35

		c := color.RGBA{
			R: 48 + uint8((235-48)*buy) + uint8((106-48)*sell),
			G: 51 + uint8((77-51)*buy) + uint8((176-51)*sell),
			B: 107 + uint8((75-107)*buy) + uint8((76-107)*sell),
			A: 150,
		}

		y := rsi * 100
		for i := 0.0; i < y; i++ {
			for j := 0; j < 3; j++ {
				plot.Set((x-screen.Camera.X)*3+j, int(i), c)
			}
		}
	}

}
