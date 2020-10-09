package stocks

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Model struct {
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
	fmt.Printf("cleaning dataset containing %d data points\n", len(rawRows))
	benchStart := time.Now()

	rows := make([]Quote, 0)

	// set start of data
	startTimestamp, err := strconv.ParseInt(strings.Split(rawRows[0], ",")[0], 10, 64)
	handleFatal(err)
	startTime := time.Unix(startTimestamp, 0)

	startMinute, _ := strconv.Atoi(startTime.Format("04"))
	for startMinute > 30 || startMinute < 30 {
		startTime = startTime.Add(time.Minute)
		startMinute, _ = strconv.Atoi(startTime.Format("04"))
	}
	startHour, _ := strconv.Atoi(startTime.Format("15"))
	for startHour > 18 || startHour < 11 {
		startTime = startTime.Add(time.Hour)
		startHour, _ = strconv.Atoi(startTime.Format("15"))
	}

	for _, row := range rawRows {

		if startTime.Format("15:04") == "18:01" {
			startTime = startTime.Add(time.Hour*17 + time.Minute*30)
		}

		if startTime.Format("Mon") == "Sat" {
			startTime = startTime.Add(time.Hour * 24 * 2)
		}

		items := strings.Split(row, ",")

		timestamp, err := strconv.ParseInt(items[0], 10, 64)
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

		quoteTime := time.Unix(timestamp, 0)
		stu := startTime.Unix()
		qtu := quoteTime.Unix()
		for stu != qtu {
			if qtu < stu {
				break
			} else {
				j := len(rows) - 1
				rows = append(rows, Quote{
					Time:   startTime.Unix(),
					Open:   rows[j].Close,
					High:   rows[j].Close,
					Low:    rows[j].Close,
					Close:  rows[j].Close,
					Volume: 0,
				})
				startTime = startTime.Add(time.Minute)
			}
			stu = startTime.Unix()
		}
		if qtu < stu {
			continue
		}
		//if startTime.Format("15:04") == ""
		startTime = startTime.Add(time.Minute)

		rows = append(rows, Quote{
			Time:   timestamp,
			Open:   priceOpen,
			High:   priceHigh,
			Low:    priceLow,
			Close:  priceClose,
			Volume: volume,
		})
	}

	elapsed := time.Now().Sub(benchStart).Milliseconds()
	fmt.Printf("final set contains %d data points, cleanup took %dms\n", len(rows), elapsed)

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
