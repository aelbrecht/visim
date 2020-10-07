package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"image"
	"log"
	"visim.muon.one/internal/inputs"
	"visim.muon.one/internal/plots"
	"visim.muon.one/internal/stocks"
	"visim.muon.one/internal/view"

	"github.com/hajimehoshi/ebiten"
)

type Game struct {
	History []stocks.Quote
	Screen  *view.Screen
	Plot    *image.RGBA
}

func (g *Game) Update(screen *ebiten.Image) error {
	inputs.HandleCamera(g.Screen)

	for i := 0; i < len(g.Plot.Pix)/4; i++ {
		g.Plot.Pix[i*4] = 19
		g.Plot.Pix[i*4+1] = 15
		g.Plot.Pix[i*4+2] = 64
		g.Plot.Pix[i*4+3] = 255
	}

	g.Plot.SubImage(image.Rectangle{Max: image.Point{g.Screen.Window.W, g.Screen.Window.H}})

	plots.PlotCandles(g.History, g.Plot, g.Screen)

	screen.ReplacePixels(g.Plot.Pix)

	debug := fmt.Sprintf("%d,%d\n%d", g.Screen.Camera.X, g.Screen.Camera.Y, int(ebiten.CurrentFPS()))
	ebitenutil.DebugPrint(screen, debug)
	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.Screen.Window.W, g.Screen.Window.H
}

func handleFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	data := stocks.GetDataCSV("./data/msft.csv")

	w, h := 1280, 800

	game := Game{
		History: data,
		Screen: &view.Screen{
			Camera: &view.Camera{0, 0, 0},
			Window: view.Window{w, h},
		},
		Plot: image.NewRGBA(image.Rectangle{
			Max: image.Point{w, h},
		}),
	}

	ebiten.SetWindowSize(game.Screen.Window.W, game.Screen.Window.H)
	ebiten.SetWindowTitle("Muon Trade Sim")
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
