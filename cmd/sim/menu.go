package main

import (
	"fmt"
	"visim.muon.one/internal/layout"
)

var menuButtons []*layout.Button

func runBot() {
	fmt.Println("test")
}

func init() {
	menuButtons = make([]*layout.Button, 0)
	addButton := func(b *layout.Button) int {
		menuButtons = append(menuButtons, b)
		return b.GetShape().Max.X
	}

	var x int
	x = addButton(layout.NewButton(8, 8, "Run", func() {
		runBot()
	}))
	x = addButton(layout.NewButton(x+8, 8, "Pause", func() {
		runBot()
	}))
	x = addButton(layout.NewButton(x+8, 8, "Reset", func() {
		runBot()
	}))
	x = addButton(layout.NewButton(x+8, 8, "Set Start", func() {
		runBot()
	}))
	x = addButton(layout.NewButton(x+8, 8, "Set End", func() {
		runBot()
	}))
}
