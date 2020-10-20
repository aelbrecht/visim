package main

import (
	"github.com/hajimehoshi/ebiten"
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

func drawHorizontalLine(y float64, g *Game) {
	w := float64(g.Screen.Program.W)
	op := ebiten.DrawImageOptions{}
	op.GeoM.Scale(w, 2)
	op.GeoM.Translate(0, y)
	g.Buffers.Plot.DrawImage(borderPixel, &op)
}