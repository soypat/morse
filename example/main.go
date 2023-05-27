//go:build tinygo

package main

import (
	"machine"
	"time"

	"github.com/soypat/morse"
)

// This is the hardware abstraction layer (HAL) for the telegraph's Pin.
func telegraphPinHAL(b bool) {
	machine.LED.Set(b)
}

func main() {
	machine.LED.Configure(machine.PinConfig{Mode: machine.PinOutput})
	telegraph := morse.NewTelegraph(100*time.Millisecond, telegraphPinHAL)
	for {
		println(`sending "HELLO WORLD"`)
		telegraph.Send("HELLO WORLD")
		time.Sleep(5 * time.Second)
	}
}
