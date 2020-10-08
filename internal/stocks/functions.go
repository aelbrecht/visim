package stocks

func (m *Model) GetQuote(i int) *Quote {
	if i < 0 || i >= len(m.Quotes) {
		return nil
	}
	return &m.Quotes[i]
}

func (m *Model) GetQuoteRange(a int, b int) []Quote {
	if a < 0 || b >= len(m.Quotes) {
		return nil
	}
	return m.Quotes[a:b]
}
