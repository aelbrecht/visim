package plots

import (
	"github.com/hajimehoshi/ebiten"
	"image/color"
	"math"
	"visim.muon.one/internal/stocks"
	"visim.muon.one/internal/view"
)

var lineX *ebiten.Image
var lineY *ebiten.Image

func Axis(plot *ebiten.Image, screen *view.Screen) {

	if lineX == nil || lineY == nil {
		c := color.RGBA{R: 104, G: 109, B: 224, A: 25}
		lineX, _ = ebiten.NewImage(stocks.MinutesInDay, 1, ebiten.FilterDefault)
		lineX.Fill(c)
		lineY, _ = ebiten.NewImage(1, screen.Window.H, ebiten.FilterDefault)
		lineY.Fill(c)
	}

	ly := math.Floor(screen.Camera.Bottom)
	for ly < screen.Camera.Top {
		y := (ly - screen.Camera.Bottom) * screen.Camera.ScaleY
		op := ebiten.DrawImageOptions{}
		op.GeoM.Translate(0, y)
		op.GeoM.Scale(3, 1)
		_ = plot.DrawImage(lineX, &op)
		ly += 1
	}

	ly = math.Floor(screen.Camera.Bottom)
	for ly < screen.Camera.Top {
		y := (ly - screen.Camera.Bottom) * screen.Camera.ScaleY
		op := ebiten.DrawImageOptions{}
		op.GeoM.Translate(0, y)
		op.GeoM.Scale(3, 1)
		op.ColorM.Scale(1, 1, 1, 0.5)
		_ = plot.DrawImage(lineX, &op)
		ly += 0.05
	}

	lx := 0
	for lx < stocks.MinutesInDay*screen.Camera.ScaleX {
		op := ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(lx), 0)
		op.ColorM.Scale(1, 1, 1, 0.5)
		_ = plot.DrawImage(lineY, &op)
		lx += 10 * screen.Camera.ScaleX
	}
}
