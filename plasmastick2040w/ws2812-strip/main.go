package main

import (
	"image/color"
	"machine"
	"math/rand"
	"time"

	"tinygo.org/x/drivers/ws2812"
)

var stripLeds [64]color.RGBA

func main() {
	stripPin := machine.GPIO15
	stripPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	strip := ws2812.New(stripPin)

	for {
		for i := 0; i < len(stripLeds); i++ {
			red := uint8(rand.Uint32() & uint32(0x16))
			green := uint8(rand.Uint32() & uint32(0x16))
			blue := uint8(rand.Uint32() & uint32(0x16))
			stripLeds[i] = color.RGBA{R: red, G: green, B: blue}
		}
		strip.WriteColors(stripLeds[:])
		time.Sleep(time.Duration(rand.Uint32()&500) * time.Millisecond)
	}
}
