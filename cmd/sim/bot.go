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

			dayHistory := m.GetQuoteDay(stocks.GetDay(m.Bot.Position)).Quotes[:stocks.GetMinute(m.Bot.Position)]

			data := fmt.Sprintf("%d\n", len(dayHistory))
			for _, quote := range dayHistory {
				data += fmt.Sprintf("%d,%f,%f,%f,%f,%d\n", quote.Time, quote.Open, quote.High, quote.Low, quote.Close, quote.Volume)
			}

			cmd := exec.Command(os.Getenv("PYTHON_PATH"), os.Getenv("BOT_PATH"))
			cmd.Stdin = strings.NewReader(data)

			out, err := cmd.Output()
			if err != nil {
				fmt.Println(err)
				continue
			}

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

			orderRaw := strings.Split(result, ",")
			if len(orderRaw) != 2 {
				fmt.Println("invalid bot reply")
				continue
			}

			orderAmount, err := strconv.Atoi(orderRaw[1])
			if err != nil {
				fmt.Println("invalid order size")
				continue
			}

			if orderRaw[0] != "sell" && orderRaw[0] != "buy" && orderRaw[0] != "hold" {
				fmt.Println("invalid order type")
				continue

			}

			if orderRaw[0] == "hold" {
				m.Bot.Position += 1
				continue
			}

			order := stocks.Order{
				Buy:    orderRaw[0] == "buy",
				Amount: orderAmount,
				Quote:  m.GetQuote(m.Bot.Position),
			}

			date := time.Unix(m.GetQuote(m.Bot.Position).Time, 0).In(time.FixedZone("GMT", 0))
			fmt.Printf("%s order at %s for %d\n", orderRaw[0], date.Format(time.RFC3339), orderAmount)

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

	for i := left; i < right; i++ {

		g.Model.Bot.OrderLock.Lock()
		o := g.Model.Bot.Orders[i]
		g.Model.Bot.OrderLock.Unlock()
		if o == nil {
			continue
		}

		op := ebiten.DrawImageOptions{}
		if o.Buy {
			op.GeoM.Scale(1, -float64(o.Amount/2))
		} else {
			op.GeoM.Scale(1, float64(o.Amount/2))
		}
		op.GeoM.Translate(float64(i-left), float64(g.Screen.Program.H-100))
		op.GeoM.Scale(g.Screen.Camera.ScaleXF, 1)
		if o.Buy {
			s.DrawImage(pixelBotStart, &op)
		} else {
			s.DrawImage(pixelBotEnd, &op)
		}

	}

}
