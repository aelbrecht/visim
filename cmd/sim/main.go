package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/text"
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
	Model   *stocks.Model
	Screen  *view.Screen
	Buffers Buffers
	Options inputs.Options
}

type DayBuffer struct {
	Update    bool
	RSI       *ebiten.Image
	SR        *ebiten.Image
	Bollinger *ebiten.Image
	Candles   *ebiten.Image
	Plot      *ebiten.Image
}

type Buffers struct {
	Draw    *ebiten.Image
	Menu    *ebiten.Image
	Plot    *ebiten.Image
	Day     map[int]*DayBuffer
	Cursor  *ebiten.Image
	Tooltip *ebiten.Image
}

var (
	RSIRange        = 14
	BollingerRange  = 27
	SRRange         = 20
	MenuHeight      = 40
	ColorBlack      = color.RGBA{R: 20, G: 20, B: 20, A: 255}
	ColorAxis       = color.RGBA{R: 48, G: 51, B: 107, A: 255}
	ColorBackground = color.RGBA{R: 19, G: 15, B: 64, A: 255}
)

var borderPixel *ebiten.Image
var timelinePixel *ebiten.Image
var backgroundPixel *ebiten.Image

func init() {
	borderPixel, _ = ebiten.NewImage(1, 1, ebiten.FilterDefault)
	borderPixel.Fill(ColorAxis)
	timelinePixel, _ = ebiten.NewImage(1, 1, ebiten.FilterDefault)
	timelinePixel.Fill(ColorBlack)
	backgroundPixel, _ = ebiten.NewImage(1, 1, ebiten.FilterDefault)
	backgroundPixel.Fill(ColorBackground)
}

func makeDayBuffer(data *stocks.MarketDay, screen *view.Screen) *DayBuffer {

	min, max := data.GetRange()
	minMax := max - min

	rsiImage, err := ebiten.NewImage(stocks.MinutesInDay, 100, ebiten.FilterDefault)
	handleFatal(err)
	plotImage, err := ebiten.NewImage(stocks.MinutesInDay*screen.Camera.GridSize, screen.Program.H, ebiten.FilterDefault)
	handleFatal(err)
	candlesImage, err := ebiten.NewImage(stocks.MinutesInDay*screen.Camera.GridSize, int(minMax*100), ebiten.FilterDefault)
	handleFatal(err)
	bollingerImage, err := ebiten.NewImage(stocks.MinutesInDay, int(minMax*100), ebiten.FilterDefault)
	handleFatal(err)
	srImage, err := ebiten.NewImage(stocks.MinutesInDay, int(minMax*100), ebiten.FilterDefault)
	handleFatal(err)

	return &DayBuffer{
		Update:    true,
		RSI:       rsiImage,
		Candles:   candlesImage,
		Plot:      plotImage,
		Bollinger: bollingerImage,
		SR:        srImage,
	}
}

func (g *Game) PlotDay(day int) {

	data := g.Model.GetQuoteDay(day)
	if data == nil {
		return
	}

	cam := g.Screen.Camera

	if g.Buffers.Day[day] == nil {
		g.Buffers.Day[day] = makeDayBuffer(data, g.Screen)
	}
	b := g.Buffers.Day[day]

	// update textures if needed
	if b.Update {
		data := g.Model.GetQuoteDay(day)
		plots.Candles(data, b.Candles)
		plots.RSI(RSIRange, data, b.RSI)
		plots.Bollinger(BollingerRange, data, b.Bollinger)
		plots.SR(SRRange, data, b.SR)
		b.Update = false
	}

	// draw textures to buffer in on axis

	min, _ := data.GetRange()
	bottomDelta := (min - cam.Bottom) * cam.ScaleY
	gs := float64(cam.GridSize)

	// draw axis
	b.Plot.Clear()
	plots.Axis(b.Plot, g.Screen)

	plotMargin := 300 - float64(MenuHeight)

	if g.Options.ShowBollinger {
		op := ebiten.DrawImageOptions{}
		op.GeoM.Scale(gs, cam.ScaleY/100)
		op.GeoM.Translate(0, bottomDelta)
		op.GeoM.Translate(0, plotMargin)
		b.Plot.DrawImage(b.Bollinger, &op)
	}

	// draw candles
	if g.Options.ShowQuotes {
		op := ebiten.DrawImageOptions{}
		op.GeoM.Scale(gs/5, cam.ScaleY/100)
		op.GeoM.Translate(0, bottomDelta)
		op.GeoM.Translate(0, plotMargin)
		b.Plot.DrawImage(b.Candles, &op)
	}

	if g.Options.ShowSupportResistance {
		op := ebiten.DrawImageOptions{}
		op.GeoM.Scale(gs, cam.ScaleY/100)
		op.GeoM.Translate(0, bottomDelta)
		b.Plot.DrawImage(b.SR, &op)
	}

	op := ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(stocks.MinutesInDay)*gs, 236)
	op.GeoM.Translate(0, 24)
	b.Plot.DrawImage(backgroundPixel, &op)

	// draw rsi bars
	if g.Options.ShowRSI {
		op := ebiten.DrawImageOptions{}
		op.GeoM.Scale(gs, 1)
		b.Plot.DrawImage(b.RSI, &op)
	}

	// draw plot
	op = ebiten.DrawImageOptions{}
	op.GeoM.Scale(1, -1)
	op.GeoM.Translate(0, float64(g.Screen.Program.H))
	op.GeoM.Scale(cam.ScaleXF/gs, 1)
	op.GeoM.Translate(float64(stocks.MinutesInDay*day)*cam.ScaleXF, 0)
	op.GeoM.Translate(-float64(cam.X)*cam.ScaleXF, 0)
	g.Buffers.Plot.DrawImage(b.Plot, &op)

	// draw borders
	t := float64(MenuHeight)

	drawHorizontalLine(t, g)
	drawHorizontalLine(t+float64(g.Screen.Plot.H), g)
	drawHorizontalLine(float64(g.Screen.Program.H)-24-2, g)
}

func drawHorizontalLine(y float64, g *Game) {
	w := float64(g.Screen.Program.W)
	op := ebiten.DrawImageOptions{}
	op.GeoM.Scale(w, 2)
	op.GeoM.Translate(0, y)
	g.Buffers.Plot.DrawImage(borderPixel, &op)
}

func (g *Game) Update(screen *ebiten.Image) error {

	inputs.HandleCamera(g.Screen)
	inputs.HandlePlot(&g.Options)
	inputs.HandleBot(g.Model, g.Screen)

	g.Screen.AutoYAxis(g.Model)

	// clear existing buffers
	g.Buffers.Tooltip.Clear()

	screen.Fill(ColorBackground)
	g.Buffers.Plot.Clear()

	v0, v1 := g.Screen.VisibleDays()
	for i := v0; i <= v1; i++ {
		g.PlotDay(i)
	}

	screen.DrawImage(g.Buffers.Plot, nil)

	quoteIndex := g.Screen.Camera.X + int(float64(g.Screen.Cursor.X)/g.Screen.Camera.ScaleXF)
	plots.TooltipCandle(quoteIndex, g.Model, g.Buffers.Tooltip, g.Screen)
	plots.TooltipRSI(quoteIndex, RSIRange, g.Model, g.Buffers.Tooltip, g.Screen)

	// draw text for plot
	ly := math.Floor(g.Screen.Camera.Bottom)
	for ly < g.Screen.Camera.Top {
		y := -300 + int((ly - g.Screen.Camera.Bottom +1) * g.Screen.Camera.ScaleY)
		text.Draw(
			screen,
			fmt.Sprintf("%d", int(ly)),
			fonts.FaceHuge,
			10,
			g.Screen.Plot.H-y-10,
			color.RGBA{104, 109, 224, 150},
		)
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
	op.GeoM.Scale(g.Screen.Camera.ScaleXF, 1)
	op.GeoM.Translate(float64(g.Model.Bot.Cursor-g.Screen.Camera.X)*g.Screen.Camera.ScaleXF, 0)
	g.Buffers.Cursor.Fill(color.RGBA{246, 229, 141, 50})
	screen.DrawImage(g.Buffers.Cursor, &op)

	// draw bot position
	op = ebiten.DrawImageOptions{}
	op.GeoM.Scale(g.Screen.Camera.ScaleXF, 1)
	op.GeoM.Translate(float64(g.Model.Bot.Position-g.Screen.Camera.X)*g.Screen.Camera.ScaleXF, 0)
	g.Buffers.Cursor.Fill(color.RGBA{0, 168, 255, 100})
	screen.DrawImage(g.Buffers.Cursor, &op)

	op = ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(g.Screen.Program.W), 24)
	op.GeoM.Translate(0, float64(g.Screen.Program.H)-24)
	screen.DrawImage(timelinePixel, &op)

	g.Buffers.Menu.Fill(ColorBlack)
	screen.DrawImage(g.Buffers.Menu, nil)

	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.Screen.Program.W, g.Screen.Program.H
}

func handleFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	data := stocks.GetDataCSV("./data/msft.csv")
	// data := stocks.GetData("AAPL","2020-10-01","2020-10-03")

	programWindow := view.Window{1280, 900}
	plotWindow := view.Window{W: 1280, H: 600}

	bufferDraw, _ := ebiten.NewImage(programWindow.W, programWindow.H, ebiten.FilterDefault)
	bufferPlot, _ := ebiten.NewImage(programWindow.W, programWindow.H, ebiten.FilterDefault)
	tooltipPlot, _ := ebiten.NewImage(programWindow.W, programWindow.H, ebiten.FilterDefault)
	bufferCursor, _ := ebiten.NewImage(1, programWindow.H, ebiten.FilterDefault)
	bufferMenu, _ := ebiten.NewImage(programWindow.W, MenuHeight, ebiten.FilterDefault)

	game := Game{
		Model: &stocks.Model{
			Data: data,
			Bot: stocks.Bot{
				Cursor: 0,
			},
		},
		Screen: &view.Screen{
			Camera:  &view.Camera{ScaleX: 5, ScaleXF: 5, GridSize: 5, Y: 200},
			Plot:    plotWindow,
			Program: programWindow,
		},
		Buffers: Buffers{
			Menu:    bufferMenu,
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
	}

	ebiten.SetWindowSize(game.Screen.Program.W, game.Screen.Program.H)
	ebiten.SetWindowTitle("Muon Market View")
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
