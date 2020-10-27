package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"visim.muon.one/internal/stocks"
)

func exitStopLosses(m *stocks.Model) {

}

func exitTakeProfits(m *stocks.Model) {

}

func exitShortPositions(m *stocks.Model) {

}

func exitLongPositions(m *stocks.Model) {

}

func RunBot(m *stocks.Model) {

	for true {
		if m.Bot.Running {

			startPos := m.Bot.Position

			if m.Bot.Position >= m.Bot.End {
				m.Bot.Running = false
				continue
			}

			if !m.Bot.Fast {
				time.Sleep(time.Second)
			}

			exitStopLosses(m)
			exitTakeProfits(m)

			// prepare data payload
			dayHistory := m.GetQuoteDay(stocks.GetDay(m.Bot.Position)).Quotes[:stocks.GetMinute(m.Bot.Position)]
			data := fmt.Sprintf("%d\n", len(dayHistory))
			for _, quote := range dayHistory {
				data += fmt.Sprintf("%d,%f,%f,%f,%f,%d\n", quote.Time, quote.Open, quote.High, quote.Low, quote.Close, quote.Volume)
			}

			// run bot
			cmd := exec.Command(os.Getenv("PYTHON_PATH"), os.Getenv("BOT_PATH"))
			cmd.Stdin = strings.NewReader(data)
			out, err := cmd.Output()
			if err != nil {
				m.Bot.Message = err.Error()
				continue
			}

			// parse bot request
			outLines := strings.Split(string(out), "\n")
			result := ""
			for _, line := range outLines {
				if line == "" {
					continue
				} else if line[0] == '#' {
					fmt.Println(line)
					continue
				}
				result = line
			}

			// split tuple
			orderRaw := strings.Split(result, ",")
			if len(orderRaw) != 5 {
				m.Bot.Message = "invalid tuple size"
				fmt.Println("invalid bot reply")
				continue
			}

			// parse price
			orderPrice, err := strconv.ParseFloat(orderRaw[1], 64)
			if err != nil {
				m.Bot.Message = "invalid buy limit"
				continue
			}

			// parse amount
			orderAmount, err := strconv.Atoi(orderRaw[2])
			if err != nil {
				m.Bot.Message = "invalid order size"
				continue
			}

			// parse order kind
			kind := orderRaw[0]
			long := false
			short := false
			exit := false
			if kind == "hold" {
				m.Bot.Position += 1
				continue
			} else if kind == "long" {
				long = true
			} else if kind == "short" {
				short = true
			} else if kind == "exit_long" {
				exitLongPositions(m)
				exit = true
				long = true
			} else if kind == "exit_short" {
				exitShortPositions(m)
				exit = true
				short = true
			}

			// parse buy limit
			takeProfitMargin, err := strconv.ParseFloat(orderRaw[3], 64)
			if err != nil {
				m.Bot.Message = "invalid buy limit"
				continue
			}

			// parse sell limit
			stopLossMargin, err := strconv.ParseFloat(orderRaw[4], 64)
			if err != nil {
				m.Bot.Message = "invalid sell limit"
				continue
			}

			q := m.GetQuote(m.Bot.Position)

			takeProfit := 0.0
			if takeProfitMargin != 0 {
				takeProfit = (1 + takeProfitMargin) * orderPrice
			}

			stopLoss := 0.0
			if stopLossMargin != 0 {
				stopLoss = (1 - stopLossMargin) * orderPrice
			}

			order := stocks.Order{
				TakeProfit: takeProfit,
				StopLoss:   stopLoss,
				Long:       long,
				Short:      short,
				Exit:       exit,
				Amount:     orderAmount,
				Quote:      m.GetQuote(m.Bot.Position),
			}

			date := time.Unix(q.Time, 0).In(time.FixedZone("GMT", 0))
			m.Bot.Message = fmt.Sprintf("%s: %s %d %f\n", date.Format(time.RFC3339), kind, orderAmount, orderPrice)

			if startPos < m.Bot.Position {
				continue
			}

			m.Bot.OrderLock.Lock()
			m.Bot.Orders[m.Bot.Position] = &order
			m.Bot.OrderLock.Unlock()
			m.Bot.Position += 1
		} else {
			time.Sleep(time.Second)
		}
	}
}

func plotTrades(g *Game, s *ebiten.Image) {

	left := g.Screen.Camera.X
	right := left + int(float64(g.Screen.Program.W)/g.Screen.Camera.ScaleXF)

	// plot long/short enter
	for i := left; i < right; i++ {
		g.Model.Bot.OrderLock.Lock()
		o := g.Model.Bot.Orders[i]
		g.Model.Bot.OrderLock.Unlock()
		if o == nil || o.Exit {
			continue
		}

		op := ebiten.DrawImageOptions{}
		if o.Long {
			op.GeoM.Scale(1, -float64(o.Amount/2))
		} else if o.Short {
			op.GeoM.Scale(1, float64(o.Amount/2))
			op.GeoM.Translate(0, 1)
		}
		op.GeoM.Translate(float64(i-left), float64(g.Screen.Plot.H)+40+100+2+50+2)
		op.GeoM.Scale(g.Screen.Camera.ScaleXF, 1)
		if o.Long {
			s.DrawImage(pixelEnter, &op)
		} else if o.Short {
			s.DrawImage(pixelExit, &op)
		}
	}

	// plot long/short exits
	for i := left; i < right; i++ {
		g.Model.Bot.OrderLock.Lock()
		o := g.Model.Bot.Orders[i]
		g.Model.Bot.OrderLock.Unlock()
		if o == nil || !o.Exit {
			continue
		}

		op := ebiten.DrawImageOptions{}
		op.GeoM.Scale(1, 37)
		op.GeoM.Translate(float64(i-left), 40+2+float64(g.Screen.Plot.H)+2+100+2+100)
		op.GeoM.Scale(g.Screen.Camera.ScaleXF, 1)
		if o.Long {
			s.DrawImage(pixelExit, &op)
		} else if o.Short {
			s.DrawImage(pixelEnter, &op)
		} else {
			s.DrawImage(pixelHold, &op)
		}
	}
}
