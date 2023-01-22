package main

import (
	"image/color"
	"machine"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/fogleman/ease"
	"tinygo.org/x/drivers/ws2812"
)

type Dot struct {
	Origin         uint32
	Destination    uint32
	Duration       time.Duration
	Start          time.Time
	End            time.Time
	Color          color.RGBA
	EasingFunction func(float64) float64
}

var stripLeds [96]color.RGBA
var stripLedsMin = uint32(0)
var stripLedsMax = uint32(cap(stripLeds))

func main() {
	stripPin := machine.GPIO15
	stripPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	strip := ws2812.New(stripPin)

	easingFunctions := []func(float64) float64{
		// ease.InOutQuad,
		ease.InOutCubic,
		// ease.InOutQuart,
		// ease.InOutQuint,
		// ease.InOutExpo,
		// ease.InOutCirc,
	}

	dotsMutex := sync.Mutex{}
	dots := []Dot{}
	go func() {
		for {
			dot := Dot{
				Duration: time.Duration(rand.Int31n(10000)+2000) * time.Millisecond,
				Start:    time.Now(),
				Color: color.RGBA{
					R: uint8(rand.Int31n(0x4)),
					G: uint8(rand.Int31n(0x4)),
					B: uint8(rand.Int31n(0x4)),
				},
				EasingFunction: easingFunctions[rand.Int31n(int32(len(easingFunctions)))],
			}
			dot.End = dot.Start.Add(dot.Duration)
			if rand.Int31n(2) == 1 {
				dot.Origin = stripLedsMin
				dot.Destination = stripLedsMax
			} else {
				dot.Origin = stripLedsMax
				dot.Destination = stripLedsMin
			}
			dotsMutex.Lock()
			dots = append(dots, dot)
			dotsMutex.Unlock()
			time.Sleep(time.Duration(rand.Int31n(2000)+20) * time.Millisecond)
		}
	}()

	for {
		for i := 0; i < len(stripLeds); i++ {
			stripLeds[i] = color.RGBA{R: 0, G: 0, B: 0}
		}

		dotsMutex.Lock()
		for i := 0; i < len(dots); i++ {
			dot := &dots[i]
			if time.Now().After(dot.End) {
				dots = append(dots[:i], dots[i+1:]...)
				i--
				continue
			}
			elapsed := math.Min(1, float64(time.Now().Sub(dot.Start))/float64(dot.Duration))
			elapsed = dot.EasingFunction(elapsed)
			var position uint32
			if dot.Origin < dot.Destination {
				position = uint32(float64(dot.Origin) + (float64(dot.Destination-dot.Origin) * elapsed))
			} else {
				position = uint32(float64(dot.Origin) - (float64(dot.Origin-dot.Destination) * elapsed))
			}
			if position >= stripLedsMin && position < stripLedsMax {
				stripLeds[position] = dot.Color
			}
		}
		dotsMutex.Unlock()

		strip.WriteColors(stripLeds[:])
		time.Sleep(10 * time.Millisecond)
	}
}
