package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Chooser struct {
	widget.BaseWidget
	options  []string
	chosen   int
	onChange func(string)
	width    float32
	hover    bool
}

type ChooserRenderer struct {
	chooser   *Chooser
	rectangle *canvas.Rectangle
	iconLeft  *canvas.Image
	iconRight *canvas.Image
	label     *widget.RichText
}

func NewChooser(options []string, onChange func(string)) *Chooser {
	ws := &Chooser{
		chosen:   0,
		onChange: onChange,
		options:  options,
	}
	ws.ExtendBaseWidget(ws)
	return ws
}

func (c *Chooser) padding() float32 {
	return 20
}

func (c *Chooser) buttonDimension() float32 {
	return 1.5 * theme.TextSize()
}

func (c *Chooser) CreateRenderer() fyne.WidgetRenderer {
	rec := canvas.NewRectangle(theme.InputBackgroundColor())
	rec.StrokeWidth = 0
	rec.StrokeColor = theme.DisabledButtonColor()
	rec.Move(fyne.NewPos(0, 10))
	rec.Resize(fyne.NewSize(200, theme.TextSize()))
	rec.CornerRadius = theme.InputRadiusSize()

	iconLeft := canvas.NewImageFromResource(theme.NavigateBackIcon())
	iconLeft.Resize(fyne.NewSize(20, 20))
	iconLeft.Move(fyne.NewPos(5, theme.TextSize()))

	iconRight := canvas.NewImageFromResource(theme.NavigateNextIcon())
	iconRight.Resize(fyne.NewSize(20, 20))
	iconRight.Move(fyne.NewPos(195, theme.TextSize()))

	seg := &widget.TextSegment{Text: c.options[0], Style: widget.RichTextStyle{Alignment: fyne.TextAlignCenter}}
	seg.Style.Alignment = fyne.TextAlignCenter
	label := widget.NewRichText(seg)
	label.Move(fyne.NewPos(100, theme.TextSize()))

	return &ChooserRenderer{
		chooser:   c,
		rectangle: rec,

		iconLeft:  iconLeft,
		iconRight: iconRight,
		label:     label,
	}
}

func (c *Chooser) MouseUp(m *desktop.MouseEvent) {
}

func (c *Chooser) MouseDown(m *desktop.MouseEvent) {
	if m.Position.X <= 20 {
		c.chosen = (c.chosen + len(c.options) - 1) % len(c.options)
		c.onChange(c.options[c.chosen])
		c.Refresh()
	} else if m.Position.X >= c.width-20 {
		c.chosen = (c.chosen + 1) % len(c.options)
		c.onChange(c.options[c.chosen])
		c.Refresh()
	}
}

func (c *Chooser) MouseIn(m *desktop.MouseEvent) {
	c.hover = false
	if m.Position.X <= 20 || m.Position.X >= c.width-20 {
		c.hover = true
	}
}

func (c *Chooser) MouseMoved(m *desktop.MouseEvent) {
	c.hover = false
	if m.Position.X <= 20 || m.Position.X >= c.width-20 {
		c.hover = true
	}
}

func (c *Chooser) MouseOut() {
	c.hover = false
}

func (c *Chooser) Cursor() desktop.Cursor {
	if c.hover {
		return desktop.PointerCursor
	}
	return desktop.DefaultCursor
}

func (r *ChooserRenderer) Refresh() {
	r.rectangle.Refresh()
	r.iconRight.Refresh()
	r.label.Segments[0].(*widget.TextSegment).Text = r.chooser.options[r.chooser.chosen]
	r.label.Refresh()
}

func (r *ChooserRenderer) Layout(s fyne.Size) {
	r.chooser.width = s.Width
	r.rectangle.Resize(fyne.NewSize(s.Width, 2*theme.TextSize()))
	r.iconRight.Move(fyne.NewPos(s.Width-25, theme.TextSize()))
	r.label.Move(fyne.NewPos(s.Width/2-r.label.Size().Width/2, theme.TextSize()/2))
	r.Refresh()
}

func (r *ChooserRenderer) MinSize() fyne.Size {
	dim := r.chooser.buttonDimension()
	padding := r.chooser.padding()
	return fyne.NewSize(200, padding+dim)
}

func (r *ChooserRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.rectangle, r.iconLeft, r.iconRight, r.label}
}

func (r *ChooserRenderer) Destroy() {}
