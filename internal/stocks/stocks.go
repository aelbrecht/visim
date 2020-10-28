package stocks

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Model struct {
	Data []*MarketDay
	Bot  Bot
}

type Order struct {
	Finished   bool
	Long       bool
	Short      bool
	Exit       bool
	StopLoss   float64
	TakeProfit float64
	Amount     int
	Leverage   int
	EnterQuote *Quote
	ExitQuote  *Quote
}

type DailyHistory = []*MarketDay

type MarketDay struct {
	modified bool
	Quotes   []Quote
	min      float64
	max      float64
}

type Bot struct {
	Message   string
	Cursor    int
	Position  int
	Start     int
	End       int
	Running   bool
	Fast      bool
	Orders    map[int]*Order
	OrderLock sync.Mutex
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

func GetDataCSV(file string) DailyHistory {
	body, err := ioutil.ReadFile(file)
	handleFatal(err)
	return formatData(string(body))
}

func formatData(rawText string) DailyHistory {

	days := make([]*MarketDay, 0)

	rawRows := strings.Split(rawText, "\n")
	fmt.Printf("parsing dataset containing %d data points\n", len(rawRows))

	benchStart := time.Now()

	rows := make([]Quote, 0)
	points := 0

	// set start of data
	startTimestamp, err := strconv.ParseInt(strings.Split(rawRows[0], ",")[0], 10, 64)
	handleFatal(err)
	startTime := time.Unix(startTimestamp, 0).In(time.FixedZone("GMT", 0))

	startMinute, _ := strconv.Atoi(startTime.Format("04"))
	for startMinute > 30 || startMinute < 30 {
		startTime = startTime.Add(time.Minute)
		startMinute, _ = strconv.Atoi(startTime.Format("04"))
	}
	startHour, _ := strconv.Atoi(startTime.Format("15"))
	for startHour > 16 || startHour < 9 {
		startTime = startTime.Add(time.Hour)
		startHour, _ = strconv.Atoi(startTime.Format("15"))
	}

	for _, row := range rawRows {
		if startTime.Format("15:04") == "16:01" {
			startTime = startTime.Add(time.Hour*17 + time.Minute*30)
			if points > 0 {
				days = append(days, &MarketDay{
					Quotes:   rows,
					modified: true,
				})
			}
			rows = make([]Quote, 0)
			points = 0
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

		startTime = startTime.Add(time.Minute)

		rows = append(rows, Quote{
			Time:   timestamp,
			Open:   priceOpen,
			High:   priceHigh,
			Low:    priceLow,
			Close:  priceClose,
			Volume: volume,
		})
		points++
	}

	elapsed := time.Now().Sub(benchStart).Milliseconds()
	fmt.Printf("final set is %d day(s), cleanup took %dms\n", len(days), elapsed)

	return days
}

func GetData(symbol string, from string, to string) DailyHistory {

	host := os.Getenv("API_URL")

	resp, err := http.Get(fmt.Sprintf("%s/history?symbol=%s&from=%s&to=%s", host, symbol, from, to))
	handleFatal(err)

	body, err := ioutil.ReadAll(resp.Body)
	handleFatal(err)

	rawText := string(body)
	return formatData(rawText)
}
