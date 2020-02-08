package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

type slideC struct {
	l       *log.Logger
	name    string
	prev    float64
	percent *widget.Label
	x       *xrandr
}

func makeSliders(l *log.Logger, x *xrandr) *fyne.Container {
	con := fyne.NewContainerWithLayout(layout.NewGridLayout(3))
	for k, v := range x.displays {
		percent := widget.NewLabelWithStyle(fmt.Sprintf("%.f", v*100)+"%", fyne.TextAlignCenter, fyne.TextStyle{Monospace: true})
		sC := &slideC{
			l:       l,
			name:    k,
			prev:    v * 100,
			percent: percent,
			x:       x,
		}
		sW := widget.NewSlider(0, 100)
		sW.Step = 1
		sW.Value = v * 100
		sW.OnChanged = sC.onChanged
		lW := widget.NewLabelWithStyle(k, fyne.TextAlignCenter, fyne.TextStyle{Monospace: true})
		con.AddObject(lW)
		con.AddObject(sW)
		con.AddObject(percent)
	}
	return con
}

func main() {
	x := xrandr{}
	if err := x.new(time.Millisecond * 5); err != nil {
		panic(err)
	}
	l := log.New(&bytes.Buffer{}, "", log.LUTC)
	l.SetOutput(os.Stderr)
	a := app.New()
	w := a.NewWindow("brightness")
	w.SetContent(makeSliders(l, &x))
	w.Resize(fyne.NewSize(400, 1))
	w.ShowAndRun()
	close(x.death)
}

func (s *slideC) onChanged(val float64) {
	if val != s.prev {
		s.prev = val
		s.l.Println(val)
		val = val / 100
		s.x.setBrightness(s.name, val)
		s.percent.SetText(fmt.Sprintf("%.f", val*100) + "%")
		s.percent.Refresh()
	}
}
