package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/text"
	"image/color"
	"log"
	"sync"
	"visim.muon.one/internal/fonts"
	"visim.muon.one/internal/inputs"
	"visim.muon.one/internal/layout"
	"visim.muon.one/internal/plots"
	"visim.muon.one/internal/stocks"
	"visim.muon.one/internal/view"

	"github.com/hajimehoshi/ebiten"
)

type Game struct {
	Buttons []*layout.Button
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
	Plot    *ebiten.Image
	Day     map[int]*DayBuffer
	Tooltip *ebiten.Image
}

var (
	RSIRange        = 14
	BollingerRange  = 27
	SRRange         = 20
	MenuHeight      = 40
	ColorBlack      = color.RGBA{R: 20, G: 20, B: 20, A: 255}
	ColorGray       = color.RGBA{R: 50, G: 50, B: 50, A: 255}
	ColorAxis       = color.RGBA{R: 48, G: 51, B: 107, A: 255}
	ColorBackground = color.RGBA{R: 19, G: 15, B: 64, A: 255}
	ColorCursor     = color.RGBA{R: 104, G: 109, B: 224, A: 150}
)

var borderPixel *ebiten.Image
var timelinePixel *ebiten.Image
var backgroundPixel *ebiten.Image
var cursorPixel *ebiten.Image
var botCursorPixel *ebiten.Image
var selectionPixel *ebiten.Image
var menuPixel *ebiten.Image
var grayPixel *ebiten.Image

func init() {
	borderPixel, _ = ebiten.NewImage(1, 1, ebiten.FilterDefault)
	borderPixel.Fill(ColorAxis)
	timelinePixel, _ = ebiten.NewImage(1, 1, ebiten.FilterDefault)
	timelinePixel.Fill(ColorBlack)
	backgroundPixel, _ = ebiten.NewImage(1, 1, ebiten.FilterDefault)
	backgroundPixel.Fill(ColorBackground)
	cursorPixel, _ = ebiten.NewImage(1, 1, ebiten.FilterDefault)
	cursorPixel.Fill(ColorCursor)
	botCursorPixel, _ = ebiten.NewImage(1, 1, ebiten.FilterDefault)
	botCursorPixel.Fill(color.RGBA{249, 202, 36, 100})
	selectionPixel, _ = ebiten.NewImage(1, 1, ebiten.FilterDefault)
	selectionPixel.Fill(ColorCursor)
	menuPixel, _ = ebiten.NewImage(1, 1, ebiten.FilterDefault)
	menuPixel.Fill(ColorBlack)
	grayPixel, _ = ebiten.NewImage(1, 1, ebiten.FilterDefault)
	grayPixel.Fill(ColorGray)
}

func (g *Game) Update(screen *ebiten.Image) error {

	// handle inputs
	inputs.HandleMouseLeft(g.Screen, g.Buttons)
	inputs.HandlePlot(&g.Options)
	inputs.HandleBot(g.Model, g.Screen)

	// scale axis to always fit visible data
	g.Screen.AutoYAxis(g.Model)

	// follow bot
	if g.Model.Bot.Follow && g.Model.Bot.Running {
		g.Screen.Camera.XF = float64(g.Model.Bot.Position) - float64(g.Screen.Program.W/2)/g.Screen.Camera.ScaleXF
	}

	// clear screen and buffers
	screen.Fill(ColorBackground)
	g.Buffers.Tooltip.Clear()
	g.Buffers.Plot.Clear()

	// plot stock
	v0, v1 := g.Screen.VisibleDays()
	for i := v0; i <= v1; i++ {
		plotDay(g, i)
	}
	screen.DrawImage(g.Buffers.Plot, nil)

	// plot trade indicators
	plotTrades(g, screen)

	// plot tooltips
	quoteIndex := g.Screen.Camera.X + int(float64(g.Screen.Cursor.X)/g.Screen.Camera.ScaleXF)
	plots.TooltipCandle(quoteIndex, g.Model, g.Buffers.Tooltip, g.Screen)
	plots.TooltipRSI(quoteIndex, RSIRange, g.Model, g.Buffers.Tooltip, g.Screen)
	if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		screen.DrawImage(g.Buffers.Tooltip, nil)
	}

	// plot interface
	drawCursors(g, screen)
	drawVerticalLabels(g.Screen, screen)
	drawHorizontalLabels(g.Screen, g.Model, screen)
	drawMenu(g.Buttons, g.Screen, screen)

	debugMessage := g.Model.Bot.Message
	portfolio := BotPortfolio(g.Model, g.Model.Bot.Cursor)
	msgPortfolio := fmt.Sprintf(
		"stocks: %0.2f | settled: %0.2f | investment: %0.2f | long: %d | short: %d | P/L: %0.2f",
		portfolio.Stocks, portfolio.Settled, portfolio.Invested, portfolio.Long, portfolio.Short, portfolio.Profit)
	x := g.Screen.Program.W - len(msgPortfolio)*6 - 20
	text.Draw(screen, msgPortfolio, fonts.FaceNormal, x, 26, color.White)
	x = x - len(debugMessage)*6 - 20
	text.Draw(screen, debugMessage, fonts.FaceNormal, x, 26, ColorGray)

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
	//data := stocks.GetData("AAPL", "2020-02-01", "2020-05-01")

	programWindow := view.Window{W: 1400, H: 840}
	plotWindow := view.Window{W: 1400, H: 500}

	bufferDraw, _ := ebiten.NewImage(programWindow.W, programWindow.H, ebiten.FilterDefault)
	bufferPlot, _ := ebiten.NewImage(programWindow.W, programWindow.H, ebiten.FilterDefault)
	tooltipPlot, _ := ebiten.NewImage(programWindow.W, programWindow.H, ebiten.FilterDefault)

	model := &stocks.Model{
		Data: data,
		Bot: stocks.Bot{
			Cursor:    0,
			Position:  0,
			Start:     0,
			End:       stocks.MinutesInDay,
			Running:   false,
			Orders:    make(map[int]*stocks.Order),
			OrderLock: sync.Mutex{},
		},
	}
	screen := &view.Screen{
		Camera:  &view.Camera{ScaleX: 5, ScaleXF: 5, GridSize: 5, Y: 200},
		Plot:    plotWindow,
		Program: programWindow,
	}
	game := Game{
		Buttons: makeMenuButtons(model, screen),
		Model:   model,
		Screen:  screen,
		Buffers: Buffers{
			Draw:    bufferDraw,
			Plot:    bufferPlot,
			Tooltip: tooltipPlot,
			Day:     make(map[int]*DayBuffer),
		},
		Options: inputs.Options{
			ShowBollinger: true,
			ShowRSI:       true,
			ShowQuotes:    true,
		},
	}

	go RunBot(model)

	ebiten.SetRunnableOnUnfocused(true)
	ebiten.SetWindowSize(game.Screen.Program.W, game.Screen.Program.H)
	ebiten.SetWindowTitle("Muon Market View")
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
