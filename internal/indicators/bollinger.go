package indicators

import (
	"math"
	"visim.muon.one/internal/stocks"
)

func SimpleMeanAverage(quotes []stocks.Quote) float64 {
	total := 0.0
	for i := range quotes {
		total += quotes[i].Close
	}
	return total / float64(len(quotes))
}

func StandardDeviation(quotes []stocks.Quote) float64 {
	mean := SimpleMeanAverage(quotes)
	numerator := 0.0
	for i := range quotes {
		numerator += math.Pow(quotes[i].Close-mean, 2)
	}
	return math.Sqrt(numerator / (float64(len(quotes)) - 1))
}
