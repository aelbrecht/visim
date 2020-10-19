package inputs

import (
	"github.com/hajimehoshi/ebiten"
	"visim.muon.one/internal/stocks"
	"visim.muon.one/internal/view"
)

var cursorPressed = false

func HandleBot(model *stocks.Model, screen *view.Screen) {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		cursorPressed = true
	} else if cursorPressed {
		cursorPressed = false
		model.Bot.Cursor = screen.Camera.X + int(float64(screen.Cursor.X)/screen.Camera.ScaleXF)
	}

}
