package plots

import (
	"image"
	"image/color"
	"log"
	"math"
	"strconv"
	"time"
	"visim.muon.one/internal/stocks"
	"visim.muon.one/internal/view"
)

func Axis(m *stocks.Model, plot *image.RGBA, screen *view.Screen) {
	c := color.RGBA{104, 109, 224, 25}

	PlotY(func(y int, v float64) {
		sv := v-math.Round(v)
		if !(math.Abs(sv) < 0.005) {
			return
		}
		for x := 0; x < screen.Window.W; x++ {
			plot.Set(x, y, c)
		}
	}, screen)

	PlotX(func(i int) {
		q := m.GetQuote(i)
		if q == nil {
			return
		}
		t := time.Unix(q.Time, 0)
		minute, err := strconv.Atoi(t.Format("04"))
		if err != nil {
			log.Fatal(err)
		}
		if minute%10 == 0 {
			for y := 0; y < screen.Window.H; y++ {
				SetPixel(i, y, c, plot, screen)
			}
		}
	}, screen)
}
