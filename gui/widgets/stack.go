package widgets

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Stack struct {
	widget.BaseWidget
	container *fyne.Container
}

type StackRenderer struct {
	stack     *Stack
	rectangle *canvas.Rectangle
}

func NewStack(container *fyne.Container) *Stack {
	stack := &Stack{
		container: container,
	}
	stack.ExtendBaseWidget(stack)
	return stack
}

func (s *Stack) CreateRenderer() fyne.WidgetRenderer {
	rec := canvas.NewRectangle(color.NRGBA{R: 0x34, G: 0x2F, B: 0x42, A: 0xFF})
	rec.Resize(s.container.MinSize().AddWidthHeight(20, 20))
	rec.CornerRadius = 2 * theme.InputRadiusSize()

	return &StackRenderer{
		rectangle: rec,
		stack:     s,
	}
}
func (r *StackRenderer) Refresh() {
	r.rectangle.Refresh()
	r.stack.container.Refresh()
}

func (r *StackRenderer) Layout(s fyne.Size) {
	r.stack.container.Move(fyne.NewPos(20, 20))
	r.stack.container.Resize(s.SubtractWidthHeight(40, 40))
	h := r.stack.container.MinSize().Height

	r.rectangle.Resize(fyne.NewSize(s.Width-20, h+20))
	r.rectangle.Move(fyne.NewPos(10, 10))
	r.Refresh()
}

func (r *StackRenderer) MinSize() fyne.Size {
	return r.stack.container.MinSize().AddWidthHeight(40, 40)
}

func (r *StackRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.rectangle, r.stack.container}
}

func (r *StackRenderer) Destroy() {}
