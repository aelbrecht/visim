package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"image"
	"image/color"
	"log"
	"visim.muon.one/internal/indicators"
	"visim.muon.one/internal/inputs"
	"visim.muon.one/internal/plots"
	"visim.muon.one/internal/stocks"
	"visim.muon.one/internal/view"

	"github.com/hajimehoshi/ebiten"
)

type Game struct {
	History     []stocks.Quote
	Screen      *view.Screen
	Plot        *image.RGBA
	Buffers     Buffers
	ForceRender bool
}

type Buffers struct {
	Draw   *ebiten.Image
	Plot   *ebiten.Image
	Cursor *ebiten.Image
}

func clearPlot(plot *image.RGBA) {
	for i := 0; i < len(plot.Pix)/4; i++ {
		plot.Pix[i*4+3] = 0
	}
}

func plotToBuffer(g *Game) {
	g.Buffers.Draw.ReplacePixels(g.Plot.Pix)
	g.Buffers.Plot.DrawImage(g.Buffers.Draw, nil)
	clearPlot(g.Plot)
}

func (g *Game) Update(screen *ebiten.Image) error {

	inputs.HandleCamera(g.Screen)

	g.Screen.AutoYAxis(g.History)

	// only update plot if moved, reduces cpu usage
	if g.Screen.HasMoved || g.ForceRender {
		g.ForceRender = false

		g.Buffers.Plot.Fill(color.RGBA{19, 15, 64, 255})

		plots.RSI(20, g.History, g.Plot, g.Screen)
		plotToBuffer(g)

		plots.Bollinger(20, g.History, g.Plot, g.Screen)
		plotToBuffer(g)

		plots.Candles(g.History, g.Plot, g.Screen)
		plotToBuffer(g)
	}
	op := ebiten.DrawImageOptions{}
	op.GeoM.Scale(1, -1)
	op.GeoM.Translate(0, float64(g.Screen.Window.H))
	screen.DrawImage(g.Buffers.Plot, &op)

	debug := fmt.Sprintf("%d,%d\n%d\n", g.Screen.Camera.X, g.Screen.Camera.Y, int(ebiten.CurrentFPS()))
	debug += fmt.Sprintf("%d,%d\n", g.Screen.Cursor.X, g.Screen.Cursor.Y)

	quoteIndex := g.Screen.Camera.X + g.Screen.Cursor.X/3
	quoteDebug := "no quote"
	if quoteIndex > 0 && quoteIndex < len(g.History) {
		quote := g.History[quoteIndex]
		mean := 0.0
		if quoteIndex > 20 {
			mean = indicators.SimpleMeanAverage(g.History[quoteIndex-20 : quoteIndex])
		}
		quoteDebug = fmt.Sprintf("%d: %d\n%f %f %f %f %d\n%f\n", quoteIndex, quote.Time, quote.Open, quote.High, quote.Low, quote.Close, quote.Volume, mean)
	}
	op = ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(g.Screen.Cursor.X), 0)
	screen.DrawImage(g.Buffers.Cursor, &op)

	ebitenutil.DebugPrint(screen, debug+quoteDebug)
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
	bufferDraw, err := ebiten.NewImage(w, h, ebiten.FilterDefault)
	handleFatal(err)
	bufferPlot, err := ebiten.NewImage(w, h, ebiten.FilterDefault)
	handleFatal(err)
	bufferCursor, err := ebiten.NewImage(1, h, ebiten.FilterDefault)
	handleFatal(err)
	bufferCursor.Fill(color.RGBA{104, 109, 224, 150})

	game := Game{
		History: data,
		Screen: &view.Screen{
			Camera: &view.Camera{},
			Window: view.Window{w, h},
		},
		Plot: image.NewRGBA(image.Rectangle{
			Max: image.Point{w, h},
		}),
		Buffers: Buffers{
			Draw:   bufferDraw,
			Plot:   bufferPlot,
			Cursor: bufferCursor,
		},
		ForceRender: true,
	}

	ebiten.SetWindowSize(game.Screen.Window.W, game.Screen.Window.H)
	ebiten.SetWindowTitle("Muon Trade Sim")
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
