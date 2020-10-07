package plots

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"image"
	"image/color"
	"visim.muon.one/internal/fonts"
	"visim.muon.one/internal/stocks"
	"visim.muon.one/internal/view"
)

func TooltipCandle(i int, quotes []stocks.Quote, buffer *ebiten.Image, screen *view.Screen) {
	q := quotes[i]
	x := (i-screen.Camera.X)*3 + paddingLeft
	y0 := screen.Window.H - int((q.High-screen.Camera.Bottom)*screen.Camera.ScaleY)
	fonts.Background(x-3, y0+16*3+4, 120, 62, color.RGBA{48, 51, 107, 200}, buffer)
	text.Draw(buffer, fmt.Sprintf("Open:  %f", q.Open), fonts.FaceNormal, x, y0, color.White)
	text.Draw(buffer, fmt.Sprintf("High:  %f", q.High), fonts.FaceNormal, x, y0+16, color.White)
	text.Draw(buffer, fmt.Sprintf("Low:   %f", q.Low), fonts.FaceNormal, x, y0+16*2, color.White)
	text.Draw(buffer, fmt.Sprintf("Close: %f", q.Close), fonts.FaceNormal, x, y0+16*3, color.White)
}

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
