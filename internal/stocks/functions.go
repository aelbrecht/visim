package stocks

func (m *MarketDay) GetQuote(i int) *Quote {
	if i < 0 || i >= len(m.Quotes) {
		return nil
	}
	return &m.Quotes[i]
}

func (m *MarketDay) GetQuotesInRange(a int, b int) []Quote {
	if a < 0 || b >= len(m.Quotes) {
		return nil
	}
	return m.Quotes[a:b]
}

func (m *Model) GetQuoteDay(i int) *MarketDay {
	return m.Data[i]
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
