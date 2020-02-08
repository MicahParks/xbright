package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
)

type slideC struct {
	l    *log.Logger
	prev float64
	x    *xrandr
}

func makeSliders(l *log.Logger, x *xrandr) []*widget.Box {
	boxes := make([]*widget.Box, 0)
	for k, v := range x.displays {
		sC := &slideC{
			l:    l,
			prev: v * 100,
			x:    x,
		}
		sW := widget.NewSlider(0, 100)
		sW.Step = 1
		sW.Value = v * 100
		sW.OnChanged = sC.onChanged
		lW := widget.NewLabelWithStyle(k, fyne.TextAlignCenter, fyne.TextStyle{Monospace: true})
		percent := widget.NewLabelWithStyle(fmt.Sprintf("%.f", v*100)+"%", fyne.TextAlignCenter, fyne.TextStyle{Monospace: true})
		box := widget.NewHBox(lW, sW, percent)
		boxes = append(boxes, box)
	}
	return boxes
}

func main() {
	x := xrandr{}
	if err := x.new(time.Second / 100); err != nil {
		panic(err)
	}
	l := log.New(&bytes.Buffer{}, "", log.LUTC)
	l.SetOutput(os.Stderr)
	boxes := makeSliders(l, &x)
	box := widget.NewVBox()
	a := app.New()
	w := a.NewWindow("fyne brightness controller")
	w.SetContent(box)
	w.Resize(fyne.NewSize(500, 200))
	for _, b := range boxes {
		box.Append(b)
	}
	w.ShowAndRun()
	close(x.death)
}

func (s *slideC) onChanged(val float64) {
	if val != s.prev {
		s.prev = val
		s.l.Println(val)
		val = val / 100
		s.x.setBrightness("DP-1", val)
	}
}
