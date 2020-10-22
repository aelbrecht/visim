package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"image/color"
	"math"
	"strings"
	"time"
	"visim.muon.one/internal/fonts"
	"visim.muon.one/internal/inputs"
	"visim.muon.one/internal/plots"
	"visim.muon.one/internal/stocks"
	"visim.muon.one/internal/view"
)

func makeDayBuffer(data *stocks.MarketDay, screen *view.Screen) *DayBuffer {

	min, max := data.GetRange()
	minMax := max - min

	rsiImage, err := ebiten.NewImage(stocks.MinutesInDay, 100, ebiten.FilterDefault)
	handleFatal(err)
	plotImage, err := ebiten.NewImage(stocks.MinutesInDay*screen.Camera.GridSize, screen.Plot.H, ebiten.FilterDefault)
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

// draw primary plot and its axis
func plotPrimary(b *DayBuffer, o *inputs.Options, s *view.Screen, day int, plot *ebiten.Image, data *stocks.MarketDay) {

	c := s.Camera
	min, _ := data.GetRange()
	bottomDelta := (min - c.Bottom) * c.ScaleY
	gs := float64(c.GridSize)

	// draw axis
	b.Plot.Clear()
	plots.Axis(b.Plot, s)

	if o.ShowBollinger {
		op := ebiten.DrawImageOptions{}
		op.GeoM.Scale(gs, c.ScaleY/100)
		op.GeoM.Translate(0, bottomDelta)
		b.Plot.DrawImage(b.Bollinger, &op)
	}

	// draw candles
	if o.ShowQuotes {
		op := ebiten.DrawImageOptions{}
		op.GeoM.Scale(gs/5, c.ScaleY/100)
		op.GeoM.Translate(0, bottomDelta)
		b.Plot.DrawImage(b.Candles, &op)
	}

	if o.ShowSupportResistance {
		op := ebiten.DrawImageOptions{}
		op.GeoM.Scale(gs, c.ScaleY/100)
		op.GeoM.Translate(0, bottomDelta)
		b.Plot.DrawImage(b.SR, &op)
	}

	// draw plot
	op := ebiten.DrawImageOptions{}
	op.GeoM.Scale(1, -1)
	op.GeoM.Translate(0, float64(s.Program.H))
	op.GeoM.Translate(0, -300+float64(MenuHeight))
	op.GeoM.Scale(c.ScaleXF/gs, 1)
	op.GeoM.Translate(float64(stocks.MinutesInDay*day)*c.ScaleXF, 0)
	op.GeoM.Translate(-float64(c.X)*c.ScaleXF, 0)
	plot.DrawImage(b.Plot, &op)
}

// draw rsi bars
func plotRSI(b *DayBuffer, o *inputs.Options, s *view.Screen, day int, plot *ebiten.Image) {
	if !o.ShowRSI {
		return
	}
	c := s.Camera
	op := ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(c.GridSize), -1)
	op.GeoM.Translate(0, float64(s.Program.H))
	op.GeoM.Translate(0, -300+float64(MenuHeight)+102)
	op.GeoM.Scale(c.ScaleXF/float64(c.GridSize), 1)
	op.GeoM.Translate(float64(stocks.MinutesInDay*day)*c.ScaleXF, 0)
	op.GeoM.Translate(-float64(c.X)*c.ScaleXF, 0)
	plot.DrawImage(b.RSI, &op)
}

func plotDay(g *Game, day int) {

	data := g.Model.GetQuoteDay(day)
	if data == nil {
		return
	}

	// generate buffer in memory
	if g.Buffers.Day[day] == nil {
		g.Buffers.Day[day] = makeDayBuffer(data, g.Screen)
	}
	b := g.Buffers.Day[day]

	// redraw textures if needed
	if b.Update {
		data := g.Model.GetQuoteDay(day)
		plots.Candles(data, b.Candles)
		plots.RSI(RSIRange, data, b.RSI)
		plots.Bollinger(BollingerRange, data, b.Bollinger)
		plots.SR(SRRange, data, b.SR)
		b.Update = false
	}

	// draw plot and subplots
	plotPrimary(b, &g.Options, g.Screen, day, g.Buffers.Plot, data)
	plotRSI(b, &g.Options, g.Screen, day, g.Buffers.Plot)

	// draw plot dividers
	t := float64(MenuHeight)
	drawHorizontalLine(t, g)
	drawHorizontalLine(t+float64(g.Screen.Plot.H)+100, g)
	drawHorizontalLine(t+float64(g.Screen.Plot.H)-2, g)
	drawHorizontalLine(float64(g.Screen.Program.H)-24-2, g)
}

// draw date label for horizontal axis
func drawHorizontalLabels(s *view.Screen, m *stocks.Model, plot *ebiten.Image) {
	op := ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(s.Program.W), 16)
	op.GeoM.Translate(0, float64(s.Program.H)-16)
	plot.DrawImage(timelinePixel, &op)

	a, b := s.VisibleDays()
	for i := a; i <= b; i++ {

		pos := i * stocks.MinutesInDay
		dx := pos - s.Camera.X
		x := float64(dx) * s.Camera.ScaleXF

		op := ebiten.DrawImageOptions{}
		op.GeoM.Scale(s.Camera.ScaleXF, 16)
		op.GeoM.Translate(x, float64(s.Program.H)-16)
		plot.DrawImage(grayPixel, &op)

		xm := float64(dx+stocks.MinutesInDay) * s.Camera.ScaleXF
		if xm < 80 {
			x = xm - 80 + 5
		} else {
			x = math.Max(x+s.Camera.ScaleXF*2+5, 5)
		}

		q := m.GetQuote(pos)
		if q == nil {
			continue
		}

		date := time.Unix(q.Time, 0).In(time.FixedZone("GMT", 0))
		stringDate := strings.Split(date.Format(time.RFC3339), "T")
		text.Draw(
			plot,
			stringDate[0],
			fonts.FaceNormal,
			int(x),
			s.Program.H-5,
			color.RGBA{255, 255, 255, 255},
		)
	}

}

func drawMenu(s *view.Screen, plot *ebiten.Image) {
	op := ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(s.Program.W), float64(MenuHeight))
	plot.DrawImage(menuPixel, &op)

	for _, button := range menuButtons {
		button.Draw(s, plot)
	}
}

// draw price labels for vertical axis
func drawVerticalLabels(s *view.Screen, plot *ebiten.Image) {
	c := s.Camera
	ly := math.Floor(c.Bottom)
	for ly < c.Top {
		y := int((ly - c.Bottom) * c.ScaleY)
		y = s.Plot.H - y + MenuHeight
		if y > 600+MenuHeight {
			ly += 1
			continue
		}
		text.Draw(
			plot,
			fmt.Sprintf("%d", int(ly)),
			fonts.FaceHuge,
			10,
			y-10,
			color.RGBA{104, 109, 224, 150},
		)
		ly += 1
	}
}

func drawHorizontalLine(y float64, g *Game) {
	w := float64(g.Screen.Program.W)
	op := ebiten.DrawImageOptions{}
	op.GeoM.Scale(w, 2)
	op.GeoM.Translate(0, y)
	g.Buffers.Plot.DrawImage(borderPixel, &op)
}

func drawCursors(g *Game, screen *ebiten.Image) {

	// draw bot position
	op := ebiten.DrawImageOptions{}
	op.GeoM.Scale(g.Screen.Camera.ScaleXF, float64(g.Screen.Program.H))
	op.GeoM.Translate(float64(g.Model.Bot.Position-g.Screen.Camera.X)*g.Screen.Camera.ScaleXF, 0)
	screen.DrawImage(botCursorPixel, &op)

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
