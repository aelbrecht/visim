package plots

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"image/color"
	"strings"
	"time"
	"visim.muon.one/internal/fonts"
	"visim.muon.one/internal/stocks"
	"visim.muon.one/internal/view"
)

func TooltipCandle(i int, m *stocks.Model, buffer *ebiten.Image, screen *view.Screen) {
	q := m.GetQuote(i)
	if q == nil {
		return
	}
	x := (i-screen.Camera.X)*screen.Camera.ScaleX + paddingLeft
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

func Candles(data *stocks.MarketDay, plot *ebiten.Image) {

	min, _ := data.GetRange()

	for i, q := range data.Quotes {

		x := i * 3

		lb := int((q.Low - min) * 100)
		ub := int((q.High - min) * 100)
		yo := int((q.Open - min) * 100)
		yc := int((q.Close - min) * 100)

		c := color.RGBA{R: 249, G: 202, B: 36, A: 255}
		if q.Open < q.Close {
			c = color.RGBA{R: 106, G: 176, B: 76, A: 255}
		} else if q.Open > q.Close {
			c = color.RGBA{R: 235, G: 77, B: 75, A: 255}
		}

		for j := 0; j < 2; j++ {
			plot.Set(x+j, yo, c)
		}
		for j := 1; j < 3; j++ {
			plot.Set(x+j, yc, c)
		}
		for j := lb; j < ub; j++ {
			plot.Set(x+1, j, c)
		}

	}
}
