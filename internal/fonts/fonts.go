package fonts

import (
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"golang.org/x/image/font"
	"image/color"
	"io/ioutil"
	"log"
)

var FaceNormal font.Face
var FaceLarge font.Face
var FaceHuge font.Face
var bg *ebiten.Image

func init() {

	bg, _ = ebiten.NewImage(1, 1, ebiten.FilterDefault)

	monogram, err := ioutil.ReadFile("./assets/monogram.ttf")
	if err != nil {
		log.Fatal(err)
	}

	tt, err := truetype.Parse(monogram)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72

	FaceNormal = truetype.NewFace(tt, &truetype.Options{
		Size:    16,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	FaceLarge = truetype.NewFace(tt, &truetype.Options{
		Size:    32,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	FaceHuge = truetype.NewFace(tt, &truetype.Options{
		Size:    48,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
}

func Background(x int, y int, w int, h int, c color.Color, b *ebiten.Image) {
	op := ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(w), float64(h))
	op.GeoM.Scale(1, -1)
	op.GeoM.Translate(float64(x), float64(y))
	bg.Fill(c)
	b.DrawImage(bg, &op)
}
