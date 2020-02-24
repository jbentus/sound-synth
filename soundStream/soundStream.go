package soundStream

import (
	"fmt"
	"math"
	"time"

	"github.com/gazed/vu/audio"
)

type soundStream struct {
	audio        audio.Audio
	sampleRate   uint32
	channels     uint16
	sampleBits   uint16
	blockCount   int
	blockSamples int
}

// New ctor
func New(a audio.Audio, sampleRate uint32, channels uint16, sampleBits uint16, blockCount int, blockSamples int) soundStream {
	ss := soundStream{}

	ss.audio = a
	ss.sampleRate = sampleRate
	ss.channels = channels
	ss.sampleBits = sampleBits
	ss.blockCount = blockCount
	ss.blockSamples = blockSamples

	ss.initAudio()

	return ss
}

func (ss *soundStream) initAudio() {
	err := ss.audio.Init()

	if err != nil {
		fmt.Println("Error initializing audio layer", err)
	}
}

type createWaveFn func(float64) float64

// ProduceSound Produce baby!
func (ss *soundStream) ProduceSound(createWave createWaveFn) {

	var sound, buf uint64 = 0, 0

	for {

		sndWave := ss.createSoundData(createWave)

		soundData := audio.Data{
			Frequency:  ss.sampleRate,
			Channels:   ss.channels,
			SampleBits: ss.sampleBits,
			AudioData:  sndWave,
			DataSize:   uint32(len(sndWave)),
		}

		if err := ss.audio.BindSound(&sound, &buf, &soundData); err != nil {
			fmt.Println("Error binding sound data to the sound card", err)
		}

		ss.audio.PlaySound(sound, 0, 0, 0)

		time.Sleep(2000 * time.Millisecond)
	}
}

func (ss *soundStream) createSoundData(createWave createWaveFn) []byte {
	blockSamples := ss.blockCount * ss.blockSamples
	globalTime := float64(0)
	dTimeStep := 1.0 / float64(ss.sampleRate)

	// Goofy hack to get maximum integer for a type at run-time
	dMaxSample := math.Pow(2, (1*8)-1) - 1

	slice := make([]byte, blockSamples)

	for i := 0; i < blockSamples; i++ {
		newSample := clip(createWave(globalTime), 1.0) * dMaxSample

		slice[i] = byte(newSample)

		globalTime += float64(dTimeStep)
	}

	return slice
}

func clip(dSample, dMax float64) float64 {
	if dSample >= 0.0 {
		return math.Min(dSample, dMax)
	}

	return math.Max(dSample, -dMax)
}
