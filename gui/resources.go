package gui

import (
	"embed"

	"fyne.io/fyne/v2"
)

var uiFont fyne.Resource

func SetResources(fs embed.FS) {
	data, _ := fs.ReadFile("resources/NovaSquare.ttf")
	uiFont = fyne.NewStaticResource("NovaSquare", data)
}
