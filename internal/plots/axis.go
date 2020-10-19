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

	gs := float64(screen.Camera.GridSize)

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
		op.GeoM.Scale(gs, 1)
		_ = plot.DrawImage(lineX, &op)
		ly += 1
	}

	ly = math.Floor(screen.Camera.Bottom)
	for ly < screen.Camera.Top {
		y := (ly - screen.Camera.Bottom) * screen.Camera.ScaleY
		op := ebiten.DrawImageOptions{}
		op.GeoM.Translate(0, y)
		op.GeoM.Scale(gs, 1)
		op.ColorM.Scale(1, 1, 1, 0.5)
		_ = plot.DrawImage(lineX, &op)
		ly += 0.1
	}

	lx := 0.0
	for lx < float64(stocks.MinutesInDay)*gs {
		op := ebiten.DrawImageOptions{}
		op.GeoM.Translate(lx, 0)
		op.ColorM.Scale(1, 1, 1, 0.5)
		_ = plot.DrawImage(lineY, &op)
		lx += 10 * float64(screen.Camera.GridSize)
	}

	lx = 30.0 * float64(screen.Camera.GridSize)
	for lx < float64(stocks.MinutesInDay)*gs {
		op := ebiten.DrawImageOptions{}
		op.GeoM.Scale(3, 1)
		op.GeoM.Translate(lx-1, 0)
		op.ColorM.Scale(1, 1, 1, 0.5)
		_ = plot.DrawImage(lineY, &op)
		lx += 60 * float64(screen.Camera.GridSize)
	}

	op := ebiten.DrawImageOptions{}
	op.GeoM.Scale(10, 1)
	op.GeoM.Translate(6.5*60.0*float64(screen.Camera.GridSize)-10, 0)
	op.ColorM.Scale(1, 1, 1, 0.5)
	_ = plot.DrawImage(lineY, &op)
}
