package oscillator

import (
	"math"
	"math/rand"

	"github.com/ebitengine/oto/v3"
)

type Oscillator struct {
	amplitude float64
	duty      float64
	frequency float64
	phase     float64
	pan       float64
	skew      float64
	shape     uint8

	position   float64
	sampleRate float64
	player     *oto.Player
	context    *oto.Context
}

func NewOscillator(context *oto.Context, sampleRate int) *Oscillator {
	o := &Oscillator{
		frequency:  800,
		sampleRate: float64(sampleRate),
	}

	player := context.NewPlayer(o)
	o.player = player
	o.player.Play()

	o.context = context

	return o
}

func (O *Oscillator) Shape() string {
	switch O.shape {
	case 0:
		return "SINE"
	case 1:
		return "TRIANGLE"
	case 2:
		return "RECTANGLE"
	case 255:
		return "NOISE"
	default:
		return "SINE"
	}
}

func (O *Oscillator) SetShape(shape string) {
	switch shape {
	case "SINE":
		O.shape = 0
	case "TRIANGLE":
		O.shape = 1
	case "RECTANGLE":
		O.shape = 2
	case "NOISE":
		O.shape = 255
	}
}

func (O *Oscillator) Amplitude() float64 {
	return O.amplitude
}

// Sets the amplitude of the signal.
// Accepts a real number from [0, 1].
func (O *Oscillator) SetAmplitude(amplitude float64) {
	if amplitude >= 0 && amplitude <= 1 {
		O.amplitude = amplitude
	}
}

func (O *Oscillator) Duty() float64 {
	return O.duty
}

// Sets the duty cycle of rectangular wave.
// Accepts a real number from [0, 1].
func (O *Oscillator) SetDuty(duty float64) {
	if duty >= 0 && duty <= 1 {
		O.duty = duty
	}
}

func (O *Oscillator) Frequency() float64 {
	return O.frequency
}

// Sets the duty frequency of periodic signals.
// Accepts a real number from [10, 20000].
func (O *Oscillator) SetFrequency(frequency float64) {
	if frequency >= 10 && frequency <= 20000 {
		O.frequency = frequency
	}
}

func (O *Oscillator) Pan() float64 {
	return O.pan
}

// Sets the pan of the signal.
// Accepts a real number from [-1, 1],
// -1 meaning total left, and 1 total right pan.
func (O *Oscillator) SetPan(pan float64) {
	if pan >= -1 && pan <= 1 {
		O.pan = pan
	}
}

func (O *Oscillator) Phase() float64 {
	return math.Mod(O.phase, 360) / 360
}

func (O *Oscillator) SetPhase(phase float64) {
	O.phase = phase * 360
}

func (O *Oscillator) Skew() float64 {
	return O.skew
}

// Sets the skewness of the triangular wave.
// Accepts a real number from [0, 1].
func (O *Oscillator) SetSkew(skew float64) {
	if skew >= 0 && skew <= 1 {
		O.skew = skew
	}
}

func (O *Oscillator) Read(buf []byte) (int, error) {

	var signal float64
	var sample int16

	leftAmplitude := O.amplitude * (-0.5) * (O.pan - 1) * 32767
	rightAmplitude := O.amplitude * (0.5) * (O.pan + 1) * 32767

	bufLen := len(buf) / (2 * 2)

	switch O.shape {
	case 0: // Sine
		for i := 0; i < bufLen; i++ {
			signal = math.Sin(O.frequency * (O.position + float64(i) + 2*math.Pi*O.phase) / O.sampleRate * 2 * math.Pi)

			sample = int16(signal * leftAmplitude)
			buf[4*i] = byte(sample)
			buf[4*i+1] = byte(sample >> 8)
			sample = int16(signal * rightAmplitude)
			buf[4*i+2] = byte(sample)
			buf[4*i+3] = byte(sample >> 8)
		}
	case 1: // Triangle
		for i := 0; i < bufLen; i++ {
			signal = math.Mod(O.frequency*(O.position+float64(i)+O.phase)/O.sampleRate, 1.0)
			if signal <= O.skew {
				signal = 2 / O.skew * (signal - O.skew/2)
			} else {
				signal = 2 / (1 - O.skew) * (signal - (O.skew+1)/2)
			}

			sample = int16(signal * leftAmplitude)
			buf[4*i] = byte(sample)
			buf[4*i+1] = byte(sample >> 8)
			sample = int16(signal * rightAmplitude)
			buf[4*i+2] = byte(sample)
			buf[4*i+3] = byte(sample >> 8)
		}
	case 2: // Rectangle
		for i := 0; i < bufLen; i++ {
			signal = math.Mod(O.frequency*(O.position+float64(i)+O.phase)/O.sampleRate, 1.0)
			if signal < O.duty {
				signal = 1
			} else {
				signal = -1
			}

			sample = int16(signal * leftAmplitude)
			buf[4*i] = byte(sample)
			buf[4*i+1] = byte(sample >> 8)
			sample = int16(signal * rightAmplitude)
			buf[4*i+2] = byte(sample)
			buf[4*i+3] = byte(sample >> 8)
		}
	default:
		for i := 0; i < bufLen; i++ {
			signal = rand.Float64()*2 - 1

			sample = int16(signal * leftAmplitude)
			buf[4*i] = byte(sample)
			buf[4*i+1] = byte(sample >> 8)
			sample = int16(signal * rightAmplitude)
			buf[4*i+2] = byte(sample)
			buf[4*i+3] = byte(sample >> 8)
		}
	}

	O.position += float64((bufLen))
	return len(buf), nil
}
