package stocks

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type History struct {
	Quotes []Quote
	Bot    Bot
}

type Bot struct {
	Cursor int
}

type Quote struct {
	Time   int64
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume int64
}

func handleFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func GetDataCSV(file string) []Quote {
	body, err := ioutil.ReadFile(file)
	handleFatal(err)
	return formatData(string(body))
}

func formatData(rawText string) []Quote {
	rawRows := strings.Split(rawText, "\n")
	rows := make([]Quote, len(rawRows))
	for i, row := range rawRows {
		items := strings.Split(row, ",")

		time, err := strconv.ParseInt(items[0], 10, 64)
		handleFatal(err)
		priceOpen, err := strconv.ParseFloat(items[1], 64)
		handleFatal(err)
		priceHigh, err := strconv.ParseFloat(items[2], 64)
		handleFatal(err)
		priceLow, err := strconv.ParseFloat(items[3], 64)
		handleFatal(err)
		priceClose, err := strconv.ParseFloat(items[4], 64)
		handleFatal(err)
		volume, err := strconv.ParseInt(items[5], 10, 64)
		handleFatal(err)

		rows[i] = Quote{
			Time:   time,
			Open:   priceOpen,
			High:   priceHigh,
			Low:    priceLow,
			Close:  priceClose,
			Volume: volume,
		}
	}

	return rows
}

func GetData(symbol string, from string, to string) []Quote {

	host := os.Getenv("API_URL")

	resp, err := http.Get(fmt.Sprintf("%s/history?symbol=%s&from=%s&to=%s", host, symbol, from, to))
	handleFatal(err)

	body, err := ioutil.ReadAll(resp.Body)
	handleFatal(err)

	rawText := string(body)
	return formatData(rawText)
}
