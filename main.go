package main

import (
	"fmt"
	"math"
	"sound-synth/soundStream"
	"time"

	"github.com/gazed/vu/audio"
)

var dFrequencyOutput = 440.0

func makeNoise(dTime float64) float64 {
	return 0.5 * math.Sin(dFrequencyOutput*2*3.14159*dTime)
}

func main() {

	fmt.Println("Sound synthesizer v16.3.1")

	ss := soundStream.New(audio.New(), 44100, 1, 8, 32, 256) // sampleRate or frequency, channels, sampleBits, block count, block samples

	go ss.ProduceSound(makeNoise)

	time.Sleep(2000 * time.Millisecond)
}
