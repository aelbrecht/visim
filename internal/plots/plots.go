package plots

import (
	"image"
	"image/color"
	"visim.muon.one/internal/view"
)

func PlotX(f func(i int), s *view.Screen) {
	for x := s.Camera.X; x < s.Camera.X+s.Window.W/s.Camera.ScaleX; x++ {
		f(x)
	}
}

func PlotY(f func(y int, v float64), s *view.Screen) {
	dy := s.Camera.Top-s.Camera.Bottom
	dy /= float64(s.Window.H)
	for y := 0; y < s.Window.H; y++ {
		v := s.Camera.Bottom + float64(y) * dy
		f(y, v)
	}
}

func SetPixel(x int, y int, c color.Color, p *image.RGBA, s *view.Screen) {
	p.Set((x-s.Camera.X)*s.Camera.ScaleX+s.Camera.ScaleX/2, y, c)
}

func SetDash(x int, y int, w int, c color.Color, p *image.RGBA, s *view.Screen) {
	for j := -w/2; j < w/2; j++ {
		p.Set((x-s.Camera.X)*s.Camera.ScaleX+s.Camera.ScaleX/2+j, y, c)
	}
}
