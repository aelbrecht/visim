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
	ColorCursor     = color.RGBA{R: 104, G: 109, B: 224, A: 150}
)

var borderPixel *ebiten.Image
var timelinePixel *ebiten.Image
var backgroundPixel *ebiten.Image
var cursorPixel *ebiten.Image

func init() {
	borderPixel, _ = ebiten.NewImage(1, 1, ebiten.FilterDefault)
	borderPixel.Fill(ColorAxis)
	timelinePixel, _ = ebiten.NewImage(1, 1, ebiten.FilterDefault)
	timelinePixel.Fill(ColorBlack)
	backgroundPixel, _ = ebiten.NewImage(1, 1, ebiten.FilterDefault)
	backgroundPixel.Fill(ColorBackground)
	cursorPixel, _ = ebiten.NewImage(1, 1, ebiten.FilterDefault)
	cursorPixel.Fill(ColorCursor)
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
		plotDay(g, i)
	}

	screen.DrawImage(g.Buffers.Plot, nil)

	quoteIndex := g.Screen.Camera.X + int(float64(g.Screen.Cursor.X)/g.Screen.Camera.ScaleXF)
	plots.TooltipCandle(quoteIndex, g.Model, g.Buffers.Tooltip, g.Screen)
	plots.TooltipRSI(quoteIndex, RSIRange, g.Model, g.Buffers.Tooltip, g.Screen)

	// draw text for plot
	ly := math.Floor(g.Screen.Camera.Bottom)
	for ly < g.Screen.Camera.Top {
		y := int((ly - g.Screen.Camera.Bottom) * g.Screen.Camera.ScaleY)
		y = g.Screen.Plot.H - y + MenuHeight
		if y > 600+MenuHeight {
			ly += 1
			continue
		}
		text.Draw(
			screen,
			fmt.Sprintf("%d", int(ly)),
			fonts.FaceHuge,
			10,
			y-10,
			color.RGBA{104, 109, 224, 150},
		)
		ly += 1
	}

	// draw tooltip buffer
	if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		screen.DrawImage(g.Buffers.Tooltip, nil)
	}

	// draw horizontal cursor
	op := ebiten.DrawImageOptions{}
	op.GeoM.Scale(1, float64(g.Screen.Program.H))
	op.GeoM.Translate(float64(g.Screen.Cursor.X), 0)
	screen.DrawImage(cursorPixel, &op)

	// draw vertical cursor
	op = ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(g.Screen.Program.W), 1)
	op.GeoM.Translate(0, float64(g.Screen.Cursor.Y))
	screen.DrawImage(cursorPixel, &op)

	// draw bot cursor
	/*op = ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(g.Screen.Program.W), 1)
	op.GeoM.Translate(float64(g.Model.Bot.Cursor-g.Screen.Camera.X)*g.Screen.Camera.ScaleXF, 0)
	cursorPixel.Fill(ColorCursor)
	screen.DrawImage(cursorPixel, &op)

	// draw bot position
	op = ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(g.Screen.Program.W), 1)
	op.GeoM.Translate(float64(g.Model.Bot.Position-g.Screen.Camera.X)*g.Screen.Camera.ScaleXF, 0)
	cursorPixel.Fill(color.RGBA{0, 168, 255, 100})
	screen.DrawImage(cursorPixel, &op)*/

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

	programWindow := view.Window{1400, 900}
	plotWindow := view.Window{W: 1400, H: 600}

	bufferDraw, _ := ebiten.NewImage(programWindow.W, programWindow.H, ebiten.FilterDefault)
	bufferPlot, _ := ebiten.NewImage(programWindow.W, programWindow.H, ebiten.FilterDefault)
	tooltipPlot, _ := ebiten.NewImage(programWindow.W, programWindow.H, ebiten.FilterDefault)
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
