package LEDDriver

type LEDDriverMatrix interface {
}

type LEDDriver interface {
	Initialise()
	Send()
	Stop()
}
