package inputs

import "github.com/hajimehoshi/ebiten"

var (
	Pressed1 = false
	Pressed2 = false
	Pressed3 = false
)

func HandlePlot(options *Options) bool {

	update := false

	if ebiten.IsKeyPressed(ebiten.KeyF1) {
		Pressed1 = true
	} else {
		if Pressed1 {
			Pressed1 = false
			options.ShowQuotes = !options.ShowQuotes
			update = true
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyF2) {
		Pressed2 = true
	} else {
		if Pressed2 {
			Pressed2 = false
			options.ShowBollinger = !options.ShowBollinger
			update = true
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyF3) {
		Pressed3 = true
	} else {
		if Pressed3 {
			Pressed3 = false
			options.ShowRSI = !options.ShowRSI
			update = true
		}
	}

	return update
}

type Options struct {
	ShowBollinger bool
	ShowRSI       bool
	ShowQuotes    bool
}
