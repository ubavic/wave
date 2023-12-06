package main

import (
	"embed"

	"github.com/ubavic/wave/gui"
	"github.com/ubavic/wave/oscillator"
)

//go:embed resources/*
var resourcesFS embed.FS

func main() {
	sampleRate := int(48000)

	mixer, _ := oscillator.NewMixer(sampleRate)
	mixer.NewOscillator()
	mixer.NewOscillator()

	gui.SetResources(resourcesFS)
	gui.StartGui(mixer)
}
