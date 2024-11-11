package ui

import (
	"main/simulation"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
)

type AppUI struct{}

func NewAppUI() *AppUI {
	return &AppUI{}
}

func (ui *AppUI) Start() {
	application := app.New()
	application.Settings().SetTheme(theme.DarkTheme())
	window := application.NewWindow("Parking Simulator")
	window.Resize(fyne.NewSize(620, 720))
	window.CenterOnScreen()

	simView := simulation.NewSimView(window)
	simView.Render()
	go simView.Execute()
	window.ShowAndRun()
}
