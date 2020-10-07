package indicators

import "visim.muon.one/internal/stocks"

func RelativeStrengthIndex(quotes []stocks.Quote) float64 {

	n := float64(len(quotes))
	up := 0.0
	down := 0.0

	b := -1.0
	for i := range quotes {
		a := quotes[i].Close
		if b == -1.0 {
			b = a
			continue
		}
		delta := a - b
		if delta > 0 {
			up += delta
		} else {
			down += -delta
		}
		b = a
	}

	up /= n
	down /= n

	rs := up / down
	rsi := rs/(rs+1)

	return rsi
}
