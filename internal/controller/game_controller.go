package controller

type Controller struct {
	ScreenWidth  int
	ScreenHeight int
	AspectRatio  float32
	IsReady      bool
}

func NewGameController(width int, height int) *Controller {
	return &Controller{
		ScreenWidth:  width,
		ScreenHeight: height,
		AspectRatio:  float32(width) / float32(height),
		IsReady:      false,
	}
}
