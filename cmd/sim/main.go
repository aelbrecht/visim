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

type DayBuffer struct {
	Update  bool
	RSI     *ebiten.Image
	Candles *ebiten.Image
	Plot    *ebiten.Image
}

type Buffers struct {
	Draw    *ebiten.Image
	Plot    *ebiten.Image
	Day     map[int]*DayBuffer
	Cursor  *ebiten.Image
	Tooltip *ebiten.Image
}

func makeDayBuffer(data *stocks.MarketDay, screen *view.Screen) *DayBuffer {

	min, max := data.GetRange()
	minMax := max - min

	rsiImage, err := ebiten.NewImage(view.MinutesInDay, 100, ebiten.FilterDefault)
	handleFatal(err)
	plotImage, err := ebiten.NewImage(screen.Window.W, screen.Window.H, ebiten.FilterDefault)
	handleFatal(err)
	candlesImage, err := ebiten.NewImage(view.MinutesInDay*3, int(minMax*100), ebiten.FilterDefault)
	handleFatal(err)

	return &DayBuffer{
		Update:  true,
		RSI:     rsiImage,
		Candles: candlesImage,
		Plot:    plotImage,
	}
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

func (g *Game) PlotDay(day int, screen *ebiten.Image) {

	if g.Buffers.Day[day] == nil {
		g.Buffers.Day[day] = makeDayBuffer(g.Model.GetQuoteDay(day), g.Screen)
	}

	b := g.Buffers.Day[day]
	data := g.Model.GetQuoteDay(day)
	cam := g.Screen.Camera

	if b.Update {

		data := g.Model.GetQuoteDay(day)

		fmt.Printf("plot day %d rendered\n", day)

		plots.Candles(data, b.Candles)

		plots.RSI(14, data, b.RSI)

		b.Update = false
	}

	// only update plot if moved, reduces cpu usage
	if g.Screen.HasMoved || g.ForceRender {

		// draw axis
		b.Plot.Clear()
		plots.Axis(g.Model.GetQuoteDay(day), b.Plot, g.Screen)

		// draw candles
		min, _ := data.GetRange()
		bottomDelta := (min - cam.Bottom) * cam.ScaleY
		op := ebiten.DrawImageOptions{}
		op.GeoM.Scale(1, cam.ScaleY/100)
		op.GeoM.Translate(0, bottomDelta)
		b.Plot.DrawImage(b.Candles, &op)

		// draw rsi bars
		op = ebiten.DrawImageOptions{}
		op.GeoM.Scale(3, 1)
		b.Plot.DrawImage(b.RSI, &op)

		/*

			if g.Options.ShowBollinger {
				plots.Bollinger(27, g.Model.GetQuoteDay(day), g.Plot, g.Screen)
				plotToBuffer(g)
			}

			if g.Options.ShowQuotes {
				plots.Candles(g.Model.GetQuoteDay(day), g.Plot, g.Screen)
				plotToBuffer(g)
			}

			if g.Options.ShowSupportResistance {
				plots.Resistance(5, g.Model.GetQuoteDay(day), g.Plot, g.Screen)
				plotToBuffer(g)
			}*/
	}

	// draw plot
	op := ebiten.DrawImageOptions{}
	op.GeoM.Scale(1, -1)
	op.GeoM.Translate(0, float64(g.Screen.Window.H))
	op.GeoM.Translate(float64(view.MinutesInDay*day*cam.ScaleX), 0)
	op.GeoM.Translate(-float64(cam.X*cam.ScaleX), 0)
	g.Buffers.Plot.DrawImage(b.Plot, &op)
}

func (g *Game) Update(screen *ebiten.Image) error {

	inputs.HandleCamera(g.Screen)
	g.ForceRender = inputs.HandlePlot(&g.Options) || g.ForceRender

	g.Screen.AutoYAxis(g.Model)

	// clear existing buffers
	g.Buffers.Tooltip.Clear()

	screen.Fill(color.RGBA{R: 19, G: 15, B: 64, A: 255})
	g.Buffers.Plot.Clear()
	g.PlotDay(0, screen)
	g.PlotDay(1, screen)
	screen.DrawImage(g.Buffers.Plot, nil)
	g.ForceRender = false

	debug := fmt.Sprintf("%d", int(ebiten.CurrentFPS()))

	quoteIndex := g.Screen.Camera.X + g.Screen.Cursor.X/int(g.Screen.Camera.ScaleX)
	if quoteIndex > 0 && quoteIndex < len(g.Model.Data[0].Quotes) {
		plots.TooltipCandle(quoteIndex, g.Model.Data[0].Quotes, g.Buffers.Tooltip, g.Screen)
		if quoteIndex > 20 {
			plots.TooltipRSI(quoteIndex, 20, g.Model.Data[0].Quotes, g.Buffers.Tooltip, g.Screen)
		}
	}

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
	op := ebiten.DrawImageOptions{}
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
	// data := stocks.GetData("AAPL","2020-10-01","2020-10-03")

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
			Data: data,
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
			Day:     make(map[int]*DayBuffer),
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
