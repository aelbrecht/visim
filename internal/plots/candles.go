package plots

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"image/color"
	"math"
	"strings"
	"time"
	"visim.muon.one/internal/fonts"
	"visim.muon.one/internal/stocks"
	"visim.muon.one/internal/view"
)

var candleGreen *ebiten.Image
var candleRed *ebiten.Image
var candleYellow *ebiten.Image

func init() {
	candleGreen, _ = ebiten.NewImage(1, 1, ebiten.FilterDefault)
	candleGreen.Fill(color.RGBA{R: 235, G: 77, B: 75, A: 255})
	candleRed, _ = ebiten.NewImage(1, 1, ebiten.FilterDefault)
	candleRed.Fill(color.RGBA{R: 106, G: 176, B: 76, A: 255})
	candleYellow, _ = ebiten.NewImage(1, 1, ebiten.FilterDefault)
	candleYellow.Fill(color.RGBA{R: 249, G: 202, B: 36, A: 255})
}

func TooltipCandle(i int, m *stocks.Model, buffer *ebiten.Image, screen *view.Screen) {
	q := m.GetQuote(i)
	if q == nil {
		return
	}
	x := int(float64(i-screen.Camera.X)*screen.Camera.ScaleXF) + paddingLeft
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

		x := i * 5

		lb := (q.Low - min) * 100
		ub := (q.High - min) * 100
		yo := (q.Open - min) * 100
		yc := (q.Close - min) * 100

		var t *ebiten.Image
		if q.Open < q.Close {
			t = candleRed
		} else if q.Open > q.Close {
			t = candleGreen
		} else {
			t = candleYellow
		}

		lineHeight := lb - ub

		op := ebiten.DrawImageOptions{}
		op.GeoM.Scale(1, lineHeight)
		op.GeoM.Translate(float64(x+2), yo)
		plot.DrawImage(t, &op)

		barHeight := yc - yo
		if math.Abs(barHeight) < 1 {
			if barHeight < 0 {
				barHeight = -1
			} else {
				barHeight = 1
			}
		}

		op = ebiten.DrawImageOptions{}
		op.GeoM.Scale(3, barHeight)
		op.GeoM.Translate(float64(x+1), lb)
		plot.DrawImage(t, &op)

	}
}
