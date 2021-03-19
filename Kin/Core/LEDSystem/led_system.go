package LEDSystem

type LEDSystem interface {
	Start(chan LightCommand)
	Stop()
}
