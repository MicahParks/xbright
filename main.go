package main

import (
	"bytes"
	"log"
	"os"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
)

type sChange struct {
	l    *log.Logger
	prev float64
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
	box := widget.NewVBox(
		slider,
	)
	a := app.New()
	w := a.NewWindow("Hello")
	w.SetContent(box)
	w.Resize(fyne.NewSize(500, 200))
	w.ShowAndRun()
}

func (s *sChange) sliderChange(val float64) {
	if val != s.prev {
		s.prev = val
		s.l.Println(val)
	}
}
