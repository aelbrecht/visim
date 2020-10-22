package main

import (
	"github.com/hajimehoshi/ebiten"
	"image/color"
)

var pixelBotStart *ebiten.Image
var pixelBotEnd *ebiten.Image
var pixelBotPosition *ebiten.Image

func init() {
	pixelBotStart, _ = ebiten.NewImage(1, 1, ebiten.FilterDefault)
	pixelBotStart.Fill(color.RGBA{0, 255, 0, 100})
	pixelBotEnd, _ = ebiten.NewImage(1, 1, ebiten.FilterDefault)
	pixelBotEnd.Fill(color.RGBA{255, 0, 0, 100})
	pixelBotPosition, _ = ebiten.NewImage(1, 1, ebiten.FilterDefault)
	pixelBotPosition.Fill(color.RGBA{255, 255, 0, 100})
}

func drawCursors(g *Game, screen *ebiten.Image) {

	// draw bot position
	op := ebiten.DrawImageOptions{}
	op.GeoM.Scale(g.Screen.Camera.ScaleXF, float64(g.Screen.Program.H))
	op.GeoM.Translate(float64(g.Model.Bot.Start-g.Screen.Camera.X)*g.Screen.Camera.ScaleXF, 0)
	screen.DrawImage(pixelBotStart, &op)

	// draw bot position
	op = ebiten.DrawImageOptions{}
	op.GeoM.Scale(g.Screen.Camera.ScaleXF, float64(g.Screen.Program.H))
	op.GeoM.Translate(float64(g.Model.Bot.End-g.Screen.Camera.X)*g.Screen.Camera.ScaleXF, 0)
	screen.DrawImage(pixelBotEnd, &op)

	// draw bot position
	op = ebiten.DrawImageOptions{}
	op.GeoM.Scale(g.Screen.Camera.ScaleXF, float64(g.Screen.Program.H))
	op.GeoM.Translate(float64(g.Model.Bot.Position-g.Screen.Camera.X)*g.Screen.Camera.ScaleXF, 0)
	screen.DrawImage(pixelBotPosition, &op)

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
