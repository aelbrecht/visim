package main

import (
	"github.com/hajimehoshi/ebiten"
	"image/color"
)

var pixelEnter *ebiten.Image
var pixelExit *ebiten.Image
var pixelHold *ebiten.Image

func init() {
	pixelEnter, _ = ebiten.NewImage(1, 1, ebiten.FilterDefault)
	pixelEnter.Fill(color.RGBA{0, 255, 0, 100})
	pixelExit, _ = ebiten.NewImage(1, 1, ebiten.FilterDefault)
	pixelExit.Fill(color.RGBA{255, 0, 0, 100})
	pixelHold, _ = ebiten.NewImage(1, 1, ebiten.FilterDefault)
	pixelHold.Fill(color.RGBA{255, 255, 0, 100})
}

func drawCursors(g *Game, screen *ebiten.Image) {

	// draw bot position
	op := ebiten.DrawImageOptions{}
	op.GeoM.Scale(g.Screen.Camera.ScaleXF, float64(g.Screen.Program.H))
	op.GeoM.Translate(float64(g.Model.Bot.Start-g.Screen.Camera.X)*g.Screen.Camera.ScaleXF, 0)
	screen.DrawImage(pixelEnter, &op)

	// draw bot position
	op = ebiten.DrawImageOptions{}
	op.GeoM.Scale(g.Screen.Camera.ScaleXF, float64(g.Screen.Program.H))
	op.GeoM.Translate(float64(g.Model.Bot.End-g.Screen.Camera.X)*g.Screen.Camera.ScaleXF, 0)
	screen.DrawImage(pixelExit, &op)

	// draw bot position
	op = ebiten.DrawImageOptions{}
	op.GeoM.Scale(g.Screen.Camera.ScaleXF, float64(g.Screen.Program.H))
	op.GeoM.Translate(float64(g.Model.Bot.Position-g.Screen.Camera.X)*g.Screen.Camera.ScaleXF, 0)
	screen.DrawImage(pixelHold, &op)

	// draw horizontal cursor
	op = ebiten.DrawImageOptions{}
	op.GeoM.Scale(1, float64(g.Screen.Program.H))
	op.GeoM.Translate(float64(g.Screen.Cursor.X), 0)
	screen.DrawImage(cursorPixel, &op)

	// draw vertical cursor
	op = ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(g.Screen.Program.W), 1)
	op.GeoM.Translate(0, float64(g.Screen.Cursor.Y))
	screen.DrawImage(cursorPixel, &op)

	// draw selection
	op = ebiten.DrawImageOptions{}
	op.GeoM.Scale(g.Screen.Camera.ScaleXF, float64(g.Screen.Program.H))
	op.GeoM.Translate(float64(g.Model.Bot.Cursor-g.Screen.Camera.X)*g.Screen.Camera.ScaleXF, 0)
	screen.DrawImage(selectionPixel, &op)

}
