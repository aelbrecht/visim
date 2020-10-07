package view

type Camera struct {
	X, Y  int
	Scale float64
}

type Window struct {
	W, H int
}

type Screen struct {
	Camera *Camera
	Window Window
}
