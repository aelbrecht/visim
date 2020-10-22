package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"io/ioutil"
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

			if m.Bot.Position >= m.Bot.End {
				m.Bot.Running = false
				continue
			}

			if !m.Bot.Fast {
				time.Sleep(time.Second)
			}

			dayHistory := m.GetQuoteDay(stocks.GetDay(m.Bot.Position)).Quotes[:stocks.GetMinute(m.Bot.Position)]

			data := ""
			for _, quote := range dayHistory {
				data += fmt.Sprintf("%f,%f,%f,%f,%d\n", quote.Open, quote.High, quote.Low, quote.Close, quote.Volume)
			}

			tmpDir := os.Getenv("TMP_DIR")
			ioutil.WriteFile(tmpDir+"/quotes.txt", []byte(data), 0644)

			out, err := exec.Command(os.Getenv("PYTHON_PATH"), os.Getenv("BOT_PATH"), tmpDir+"/quotes.txt").Output()
			if err != nil {
				fmt.Println(err)
				continue
			}

			orderRaw := strings.Split(strings.ReplaceAll(string(out), "\n", ""), ",")
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

			fmt.Printf("%s order at %d for %d\n", orderRaw[0], m.Bot.Position, orderAmount)

			m.Bot.Orders[m.Bot.Position] = &order

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

		o := g.Model.Bot.Orders[i]
		if o == nil {
			continue
		}

		op := ebiten.DrawImageOptions{}
		op.GeoM.Scale(1, 5)
		op.GeoM.Translate(float64(i-left), float64(g.Screen.Program.H-100))
		op.GeoM.Scale(g.Screen.Camera.ScaleXF, 1)
		if o.Buy {
			s.DrawImage(pixelBotStart, &op)
		} else {
			s.DrawImage(pixelBotEnd, &op)
		}

	}

}
