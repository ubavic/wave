package gui

import (
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/ubavic/wave/gui/widgets"
	"github.com/ubavic/wave/oscillator"
)

var oscs []fyne.CanvasObject
var columns *fyne.Container

func StartGui(mixer *oscillator.Mixer) {
	app := app.New()
	win := app.NewWindow("Microwave")
	app.Settings().SetTheme(MyTheme{})

	oscs = make([]fyne.CanvasObject, 0)
	mixerOscs := mixer.Oscillators()

	for i := range mixerOscs {
		oscs = append(oscs, oscillatorUI(mixerOscs[i]))
	}

	columns = container.New(layout.NewHBoxLayout(), oscs...)

	insertNewOscillator := func() {
		osc := mixer.NewOscillator()
		oscGui := oscillatorUI(osc)
		columns.Add(oscGui)
	}

	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.ContentAddIcon(), insertNewOscillator),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.InfoIcon(), ShowAboutBox(win, "0.1.0")),
	)

	rows := container.New(layout.NewVBoxLayout(), toolbar, columns)

	win.SetContent(rows)

	win.ShowAndRun()
}

func oscillatorUI(osc *oscillator.Oscillator) *widgets.Stack {
	shapes := []string{
		"SINE",
		"TRIANGLE",
		"RECTANGLE",
		"NOISE",
	}

	freq := binding.NewFloat()
	freq.Set(osc.Frequency())
	freq.AddListener(binding.NewDataListener(func() {
		v, _ := freq.Get()
		osc.SetFrequency(v)
	}))
	freqLbl := widget.NewLabel("FREQUENCY")
	freqValue := widget.NewLabelWithData(binding.FloatToStringWithFormat(freq, "%0.0fHz"))
	freqTopRow := container.New(layout.NewHBoxLayout(), freqLbl, layout.NewSpacer(), freqValue)
	freqSlider := widget.NewSliderWithData(200, 20000, freq)
	freqGroup := container.New(layout.NewVBoxLayout(), freqTopRow, freqSlider)

	amp := binding.NewFloat()
	amp.AddListener(binding.NewDataListener(func() {
		v, _ := amp.Get()
		osc.SetAmplitude(v / 100)
	}))
	ampLbl := widget.NewLabel("AMPLITUDE")
	ampValue := widget.NewLabelWithData(binding.FloatToStringWithFormat(amp, "%0.0f%%"))
	ampTopRow := container.New(layout.NewHBoxLayout(), ampLbl, layout.NewSpacer(), ampValue)
	ampSlider := widget.NewSliderWithData(0, 100, amp)
	ampGroup := container.New(layout.NewVBoxLayout(), ampTopRow, ampSlider)

	phase := binding.NewFloat()
	phase.Set(osc.Phase())
	phase.AddListener(binding.NewDataListener(func() {
		v, _ := phase.Get()
		osc.SetPhase(v)
	}))
	phaseLbl := widget.NewLabel("PHASE")
	shortPhase := binding.FloatToStringWithFormat(phase, "%0.0f°")
	phaseValue := widget.NewLabelWithData(shortPhase)
	phaseContainer := container.New(layout.NewHBoxLayout(), phaseLbl, layout.NewSpacer(), phaseValue)
	phaseSlider := widget.NewSliderWithData(0, 360, phase)
	phaseGroup := container.New(layout.NewVBoxLayout(), phaseContainer, phaseSlider)

	duty := binding.NewFloat()
	duty.Set(osc.Duty())
	dutyLbl := widget.NewLabel("DUTY")
	dutyValue := widget.NewLabelWithData(binding.FloatToStringWithFormat(duty, "%0.0f%%"))
	dutyContainer := container.New(layout.NewHBoxLayout(), dutyLbl, layout.NewSpacer(), dutyValue)
	dutySlider := widget.NewSliderWithData(0, 100, duty)
	duty.AddListener(binding.NewDataListener(func() {
		v, _ := duty.Get()
		osc.SetDuty(v / 100)
	}))
	dutyGroup := container.New(layout.NewVBoxLayout(), dutyContainer, dutySlider)
	dutyGroup.Hide()

	skew := binding.NewFloat()
	skew.AddListener(binding.NewDataListener(func() {
		v, _ := skew.Get()
		osc.SetSkew(v / 100)
	}))
	skewLbl := widget.NewLabel("BIAS")
	skewValue := widget.NewLabelWithData(binding.FloatToStringWithFormat(skew, "%0.0f%%"))
	skewContainer := container.New(layout.NewHBoxLayout(), skewLbl, layout.NewSpacer(), skewValue)
	skewSlider := widget.NewSliderWithData(0, 100, skew)
	skewGroup := container.New(layout.NewVBoxLayout(), skewContainer, skewSlider)
	skewGroup.Hide()

	pan := binding.NewFloat()
	pan.Set(osc.Pan())
	pan.AddListener(binding.NewDataListener(func() {
		v, _ := pan.Get()
		osc.SetPan(v / 100)
	}))
	panLbl := widget.NewLabel("PAN")
	panValue := widget.NewLabelWithData(binding.FloatToStringWithFormat(pan, "%0.0f%%"))
	panContainer := container.New(layout.NewHBoxLayout(), panLbl, layout.NewSpacer(), panValue)
	panSlider := widget.NewSliderWithData(-100, 100, pan)
	panGroup := container.New(layout.NewVBoxLayout(), panContainer, panSlider)

	shapesChooser := widgets.NewChooser(shapes, func(s string) {
		osc.SetShape(s)
		switch s {
		case "SINE":
			freqGroup.Show()
			phaseGroup.Show()
			dutyGroup.Hide()
			skewGroup.Hide()
		case "TRIANGLE":
			freqGroup.Show()
			phaseGroup.Show()
			dutyGroup.Hide()
			skewGroup.Show()
		case "RECTANGLE":
			freqGroup.Show()
			phaseGroup.Show()
			dutyGroup.Show()
			skewGroup.Hide()
		case "NOISE":
			freqGroup.Hide()
			phaseGroup.Hide()
			dutyGroup.Hide()
			skewGroup.Hide()
		}
	})

	c := container.New(layout.NewVBoxLayout(), shapesChooser, freqGroup, ampGroup, phaseGroup, dutyGroup, skewGroup, panGroup)
	return widgets.NewStack(c)
}

func ShowAboutBox(win fyne.Window, version string) func() {
	verLabel := widget.NewLabelWithStyle("Version: "+version, fyne.TextAlignLeading, fyne.TextStyle{Italic: true})
	moreLabel := widget.NewLabel("More about the program:")
	url, _ := url.Parse("https://github.com/ubavic/microwave")
	linkLabel := widget.NewHyperlink("github.com/ubavic/microwave", url)
	spacer := widgets.NewSpacer()
	hBox := container.NewHBox(moreLabel, linkLabel)
	vBox := container.NewVBox(verLabel, hBox, spacer)

	return func() {
		dialog.ShowCustom(
			"MICROWAVE",
			"Close",
			vBox,
			win,
		)
	}
}
