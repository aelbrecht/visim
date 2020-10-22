package layout

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"image"
	"image/color"
	"visim.muon.one/internal/fonts"
	"visim.muon.one/internal/view"
)

var buttonPixel *ebiten.Image
var buttonHoverPixel *ebiten.Image

func init() {
	buttonPixel, _ = ebiten.NewImage(1, 1, ebiten.FilterDefault)
	buttonPixel.Fill(color.RGBA{50, 50, 50, 255})
	buttonHoverPixel, _ = ebiten.NewImage(1, 1, ebiten.FilterDefault)
	buttonHoverPixel.Fill(color.RGBA{100, 100, 100, 255})
}

type Button struct {
	shape   image.Rectangle
	text    string
	hover   bool
	onClick func()
}

func NewButton(x int, y int, text string, onClick func()) *Button {
	return &Button{
		shape: image.Rectangle{
			Min: image.Point{X: x, Y: y},
			Max: image.Point{X: x + len(text)*6 + 16, Y: y + 8 + 8*2},
		},
		text:    text,
		hover:   false,
		onClick: onClick,
	}
}

func (b *Button) GetShape() image.Rectangle {
	return b.shape
}

func (b *Button) MouseMove(x int, y int) {
	hover := x > b.shape.Min.X && x < b.shape.Max.X && y > b.shape.Min.Y && y < b.shape.Max.Y
	b.hover = hover
}

func (b *Button) MouseClick() {
	if b.hover {
		b.onClick()
	}
}

func (b *Button) Draw(s *view.Screen, dst *ebiten.Image) {

	op := ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(b.shape.Size().X), float64(b.shape.Size().Y))
	op.GeoM.Translate(float64(b.shape.Min.X), float64(b.shape.Min.Y))
	if b.hover {
		dst.DrawImage(buttonHoverPixel, &op)
	} else {
		dst.DrawImage(buttonPixel, &op)
	}

	text.Draw(
		dst,
		b.text,
		fonts.FaceNormal,
		b.shape.Min.X+8,
		b.shape.Min.Y+16,
		color.RGBA{255, 255, 255, 255},
	)

}
