package stocks

func (m *MarketDay) GetQuote(minute int) *Quote {
	if m == nil || minute < 0 || minute >= len(m.Quotes) {
		return nil
	}
	return &m.Quotes[minute]
}

func (m *MarketDay) GetQuotesInRange(lb int, ub int) []Quote {
	lb = GetMinute(lb)
	ub = GetMinute(ub)
	if m == nil || lb < 0 || ub >= len(m.Quotes) || ub < lb {
		return nil
	}
	return m.Quotes[lb:ub]
}

func (m *Model) GetQuoteDay(day int) *MarketDay {
	if day < 0 || day >= len(m.Data) {
		return nil
	}
	return m.Data[day]
}

func (m *Model) GetQuote(x int) *Quote {
	day, minute := GetDay(x), GetMinute(x)
	return m.GetQuoteDay(day).GetQuote(minute)
}

func (m *MarketDay) update() {
	min := 999999.9
	max := 0.0
	for _, quote := range m.Quotes {
		if quote.High > max {
			max = quote.High
		}
		if quote.Low < min {
			min = quote.Low
		}
	}
	m.min = min
	m.max = max
	m.modified = false
}

func (m *MarketDay) GetRange() (float64, float64) {
	if m.modified {
		m.update()
	}
	return m.min, m.max
}

const MinutesInDay = 6*60 + 30

func GetDay(x int) int {
	return x / (MinutesInDay)
}

func GetMinute(x int) int {
	return x % MinutesInDay
}
