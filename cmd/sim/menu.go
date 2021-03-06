package main

import (
	"visim.muon.one/internal/layout"
	"visim.muon.one/internal/stocks"
	"visim.muon.one/internal/view"
)

func makeMenuButtons(m *stocks.Model, s *view.Screen) []*layout.Button {
	var menuButtons = make([]*layout.Button, 0)
	addButton := func(b *layout.Button) int {
		menuButtons = append(menuButtons, b)
		return b.GetShape().Max.X
	}

	var x int
	x = addButton(layout.NewButton(8, 8, "Run", func() {
		m.Bot.Fast = false
		m.Bot.Running = true
	}))
	x = addButton(layout.NewButton(x+8, 8, "Fast Forward", func() {
		m.Bot.Fast = true
		m.Bot.Running = true
	}))
	x = addButton(layout.NewButton(x+8, 8, "Pause", func() {
		m.Bot.Running = false
	}))
	x = addButton(layout.NewButton(x+8, 8, "Reset", func() {
		m.Bot.Running = false
		m.Bot.Position = m.Bot.Start
		m.Bot.OrderLock.Lock()
		m.Bot.Orders = make(map[int]*stocks.Order)
		m.Bot.OrderLock.Unlock()
	}))
	x = addButton(layout.NewButton(x+8, 8, "Set Start", func() {
		m.Bot.Running = false
		m.Bot.Start = m.Bot.Cursor
	}))
	x = addButton(layout.NewButton(x+8, 8, "Set End", func() {
		m.Bot.Running = false
		m.Bot.End = m.Bot.Cursor
	}))
	if m.Bot.Follow {
		x = addButton(layout.NewButton(x+8, 8, "Unfollow", func() {
			m.Bot.Follow = !m.Bot.Follow
		}))
	} else {
		x = addButton(layout.NewButton(x+8, 8, "Follow", func() {
			m.Bot.Follow = !m.Bot.Follow
		}))
	}
	x = addButton(layout.NewButton(x+8, 8, "Go to start", func() {
		s.Camera.XF = 0
	}))
	x = addButton(layout.NewButton(x+8, 8, "Go to end", func() {
		s.Camera.XF = float64(len(m.Data)*stocks.MinutesInDay - stocks.MinutesInDay/4)
	}))

	return menuButtons
}
