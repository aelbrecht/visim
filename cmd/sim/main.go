package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/text"
	"image"
	"image/color"
	"log"
	"math"
	"visim.muon.one/internal/fonts"
	"visim.muon.one/internal/inputs"
	"visim.muon.one/internal/plots"
	"visim.muon.one/internal/stocks"
	"visim.muon.one/internal/view"

	"github.com/hajimehoshi/ebiten"
)

type Game struct {
	Model       *stocks.Model
	Screen      *view.Screen
	Plot        *image.RGBA
	Buffers     Buffers
	ForceRender bool
	Options     inputs.Options
}

type Buffers struct {
	Draw    *ebiten.Image
	Plot    *ebiten.Image
	Cursor  *ebiten.Image
	Tooltip *ebiten.Image
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
	g.ForceRender = inputs.HandlePlot(&g.Options) || g.ForceRender

	g.Screen.AutoYAxis(g.Model.Quotes)

	// clear existing buffers
	g.Buffers.Tooltip.Clear()

	// only update plot if moved, reduces cpu usage
	if g.Screen.HasMoved || g.ForceRender {
		g.ForceRender = false

		g.Buffers.Plot.Fill(color.RGBA{19, 15, 64, 255})

		plots.Axis(g.Model, g.Plot, g.Screen)

		if g.Options.ShowRSI {
			plots.RSI(14, g.Model.Quotes, g.Plot, g.Screen)
			plotToBuffer(g)
		}

		if g.Options.ShowBollinger {
			plots.Bollinger(27, g.Model.Quotes, g.Plot, g.Screen)
			plotToBuffer(g)
		}

		if g.Options.ShowQuotes {
			plots.Candles(g.Model.Quotes, g.Plot, g.Screen)
			plotToBuffer(g)
		}

		if g.Options.ShowSupportResistance {
			plots.Resistance(5, g.Model, g.Plot, g.Screen)
			plotToBuffer(g)
		}
	}

	debug := fmt.Sprintf("%d", int(ebiten.CurrentFPS()))

	quoteIndex := g.Screen.Camera.X + g.Screen.Cursor.X/int(g.Screen.Camera.ScaleX)
	if quoteIndex > 0 && quoteIndex < len(g.Model.Quotes) {
		plots.TooltipCandle(quoteIndex, g.Model.Quotes, g.Buffers.Tooltip, g.Screen)
		if quoteIndex > 20 {
			plots.TooltipRSI(quoteIndex, 20, g.Model.Quotes, g.Buffers.Tooltip, g.Screen)
		}
	}

	// draw plot
	op := ebiten.DrawImageOptions{}
	op.GeoM.Scale(1, -1)
	op.GeoM.Translate(0, float64(g.Screen.Window.H))
	screen.DrawImage(g.Buffers.Plot, &op)

	// draw text for plot
	ly := math.Floor(g.Screen.Camera.Bottom)
	for ly < g.Screen.Camera.Top {
		y := int((ly - g.Screen.Camera.Bottom) * g.Screen.Camera.ScaleY)
		text.Draw(screen, fmt.Sprintf("%d", int(ly)), fonts.FaceHuge, 10, g.Screen.Window.H-y-10, color.RGBA{104, 109, 224, 150})
		ly += 1
	}

	// draw tooltip buffer
	if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		screen.DrawImage(g.Buffers.Tooltip, nil)
	}

	// draw cursor buffer
	op = ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(g.Screen.Cursor.X), 0)
	g.Buffers.Cursor.Fill(color.RGBA{104, 109, 224, 150})
	screen.DrawImage(g.Buffers.Cursor, &op)

	// draw bot cursor
	op = ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64((g.Model.Bot.Cursor-g.Screen.Camera.X)*3)+1, 0)
	g.Buffers.Cursor.Fill(color.RGBA{246, 229, 141, 150})
	screen.DrawImage(g.Buffers.Cursor, &op)

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
	bufferDraw, err := ebiten.NewImage(w, h, ebiten.FilterDefault)
	handleFatal(err)
	bufferPlot, err := ebiten.NewImage(w, h, ebiten.FilterDefault)
	handleFatal(err)
	tooltipPlot, err := ebiten.NewImage(w, h, ebiten.FilterDefault)
	handleFatal(err)
	bufferCursor, err := ebiten.NewImage(1, h, ebiten.FilterDefault)
	handleFatal(err)

	game := Game{
		Model: &stocks.Model{
			Quotes: data,
			Bot: stocks.Bot{
				Cursor: 0,
			},
		},
		Screen: &view.Screen{
			Camera: &view.Camera{ScaleX: 3},
			Window: view.Window{w, h},
		},
		Plot: image.NewRGBA(image.Rectangle{
			Max: image.Point{w, h},
		}),
		Buffers: Buffers{
			Draw:    bufferDraw,
			Plot:    bufferPlot,
			Tooltip: tooltipPlot,
			Cursor:  bufferCursor,
		},
		Options: inputs.Options{
			ShowBollinger: true,
			ShowRSI:       true,
			ShowQuotes:    true,
		},
		ForceRender: true,
	}

	ebiten.SetWindowSize(game.Screen.Window.W, game.Screen.Window.H)
	ebiten.SetWindowTitle("Muon Trade Sim")
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
