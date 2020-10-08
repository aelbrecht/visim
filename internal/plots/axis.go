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

	ly := math.Floor(screen.Camera.Bottom)
	for ly < screen.Camera.Top {
		y := int((ly - screen.Camera.Bottom) * screen.Camera.ScaleY)
		for x := 0; x < screen.Window.W; x++ {
			for j := -1; j < 2; j++ {
				plot.Set(x, y+j, c)
			}
		}
		ly += 1
	}

	ly = math.Floor(screen.Camera.Bottom)
	for ly < screen.Camera.Top {
		y := int((ly - screen.Camera.Bottom) * screen.Camera.ScaleY)
		for x := 0; x < screen.Window.W; x++ {
			plot.Set(x, y, c)
		}
		ly += 0.1
	}

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
