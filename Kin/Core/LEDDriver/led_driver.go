package LEDDriver

type LEDDriver interface {
	Initialise()
	Send()
	Stop()
}
