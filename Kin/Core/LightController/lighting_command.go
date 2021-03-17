package LightController

type LightCommand struct {
	command string
	data    []float64
}

func NewLightCommand(command string, data []float64) *LightCommand {
	return &LightCommand{command, data}
}
