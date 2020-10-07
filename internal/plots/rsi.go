package plots

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"image"
	"image/color"
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

		y := rsi * 100
		for i := 0.0; i < y; i++ {
			for j := 0; j < 3; j++ {
				plot.Set((x-screen.Camera.X)*3+j, int(i), color.RGBA{190, 46, 221, 40})
			}
		}

		for j := 0; j < 3; j++ {
			plot.Set((x-screen.Camera.X)*3+j, 100, color.RGBA{190, 46, 221, 60})
		}
	}

}
