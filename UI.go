package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

func Info(info string) {
	myApp := app.New()
	myWindow := myApp.NewWindow("提示")
	myWindow.SetContent(widget.NewLabel(info))
	myWindow.Resize(fyne.NewSize(300, 200))
	myWindow.ShowAndRun()
}
