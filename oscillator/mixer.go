package oscillator

import (
	"fmt"

	"github.com/ebitengine/oto/v3"
)

type Mixer struct {
	oscillators []*Oscillator
	context     *oto.Context
	sampleRate  int
}

func NewMixer(sampleRate int) (*Mixer, error) {
	op := &oto.NewContextOptions{
		SampleRate:   sampleRate,
		ChannelCount: 2,
		Format:       oto.FormatSignedInt16LE,
	}

	context, ready, err := oto.NewContext(op)
	if err != nil {
		return nil, fmt.Errorf("creating context: %w", err)
	}
	<-ready

	mixer := Mixer{
		context:    context,
		sampleRate: sampleRate,
	}

	return &mixer, nil
}

func (m *Mixer) Oscillators() []*Oscillator {
	return m.oscillators
}

func (m *Mixer) NewOscillator() *Oscillator {
	osc := NewOscillator(m.context, m.sampleRate)
	m.oscillators = append(m.oscillators, osc)
	return osc
}

func (m *Mixer) DeleteOscillator(index int) {
	if index < 0 || index >= len(m.oscillators) {
		return
	}

	m.oscillators = append(m.oscillators[:index], m.oscillators[index+1:]...)
}
