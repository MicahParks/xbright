package main

import (
	"bytes"
	"log"
	"os"

	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
)

type sChange struct {
	l *log.Logger
}

func main() {
	l := log.New(&bytes.Buffer{}, "", log.LUTC)
	l.SetOutput(os.Stderr)
	sC := &sChange{
		l: l,
	}
	slider := widget.NewSlider(0, 100)
	slider.Step = 1
	slider.Value = 100
	slider.OnChanged = sC.sliderChange
	app := app.New()
	w := app.NewWindow("Hello")
	w.SetContent(widget.NewVBox(
		slider,
	))
	w.ShowAndRun()
}

func (s *sChange) sliderChange(val float64) {
	s.l.Println(val)
}
