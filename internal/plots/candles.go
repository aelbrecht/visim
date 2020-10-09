package plots

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"image"
	"image/color"
	"strings"
	"time"
	"visim.muon.one/internal/fonts"
	"visim.muon.one/internal/stocks"
	"visim.muon.one/internal/view"
)

func TooltipCandle(i int, quotes []stocks.Quote, buffer *ebiten.Image, screen *view.Screen) {
	q := quotes[i]
	x := (i-screen.Camera.X)*int(screen.Camera.ScaleX) + paddingLeft
	y0 := screen.Window.H - int((q.High-screen.Camera.Bottom)*screen.Camera.ScaleY)
	fonts.Background(x-3, y0+13*6+6, 120, 16*6, color.RGBA{48, 51, 107, 200}, buffer)
	date := time.Unix(q.Time, 0).In(time.FixedZone("GMT", 0))
	stringDate := strings.Split(date.Format(time.RFC3339), "T")
	text.Draw(buffer, stringDate[0]+" "+strings.Split(stringDate[1], "Z")[0], fonts.FaceNormal, x, y0, color.White)
	text.Draw(buffer, fmt.Sprintf("Open:  %f", q.Open), fonts.FaceNormal, x, y0+16*2, color.White)
	text.Draw(buffer, fmt.Sprintf("High:  %f", q.High), fonts.FaceNormal, x, y0+16*3, color.White)
	text.Draw(buffer, fmt.Sprintf("Low:   %f", q.Low), fonts.FaceNormal, x, y0+16*4, color.White)
	text.Draw(buffer, fmt.Sprintf("Close: %f", q.Close), fonts.FaceNormal, x, y0+16*5, color.White)
}

func Candles(quotes []stocks.Quote, plot *image.RGBA, screen *view.Screen) {

	for x := screen.Camera.X; x < screen.Camera.X+screen.Window.W/int(screen.Camera.ScaleX); x++ {

		if x < 0 || x >= len(quotes) {
			continue
		}

		q := quotes[x]

		lb := int((q.Low - screen.Camera.Bottom) * screen.Camera.ScaleY)
		ub := int((q.High - screen.Camera.Bottom) * screen.Camera.ScaleY)
		yo := int((q.Open - screen.Camera.Bottom) * screen.Camera.ScaleY)
		yc := int((q.Close - screen.Camera.Bottom) * screen.Camera.ScaleY)

		c := color.RGBA{235, 77, 75, 255}
		if q.Open < q.Close {
			c = color.RGBA{106, 176, 76, 255}
		} else if q.Open == q.Close {
			c = color.RGBA{249, 202, 36, 255}
		}

		for j := 0; j < int(screen.Camera.ScaleX/2); j++ {
			plot.Set((x-screen.Camera.X)*int(screen.Camera.ScaleX)+j, yo, c)
		}
		for j := int(screen.Camera.ScaleX / 2); j < int(screen.Camera.ScaleX); j++ {
			plot.Set((x-screen.Camera.X)*int(screen.Camera.ScaleX)+j, yc, c)
		}
		for y := lb; y < ub; y++ {
			plot.Set((x-screen.Camera.X)*int(screen.Camera.ScaleX)+int(screen.Camera.ScaleX/2), y, c)
		}
	}

}
