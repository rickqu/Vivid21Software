package LEDDriver

type LEDDriver interface {
	Initialise()
	Render()
	Stop()
}
