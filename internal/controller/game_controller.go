package controller

type GameController struct {
	ScreenWidth  int
	ScreenHeight int
	AspectRatio  float32
	IsReady      bool
}

func NewGameController(width int, height int) *GameController {
	return &GameController{
		ScreenWidth:  width,
		ScreenHeight: height,
		AspectRatio:  float32(width) / float32(height),
		IsReady:      false,
	}
}
