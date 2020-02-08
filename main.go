package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/user"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

var haveSettings = true

type slideC struct {
	l       *log.Logger
	name    string
	prev    float64
	percent *widget.Label
	slider  *widget.Slider
	x       *xrandr
}

func makeSliders(l *log.Logger, x *xrandr) (*fyne.Container, []*slideC) {
	sCs := make([]*slideC, 0)
	con := fyne.NewContainerWithLayout(layout.NewGridLayout(3))
	for k, v := range x.displays {
		percent := widget.NewLabelWithStyle(fmt.Sprintf("%.f", v*100)+"%", fyne.TextAlignCenter, fyne.TextStyle{Monospace: true})
		sW := widget.NewSlider(0, 100)
		sC := &slideC{
			l:       l,
			name:    k,
			prev:    v * 100,
			percent: percent,
			slider:  sW,
			x:       x,
		}
		sCs = append(sCs, sC)
		sW.Step = 1
		sW.Value = v * 100
		sW.OnChanged = sC.onChanged
		lW := widget.NewLabelWithStyle(k, fyne.TextAlignCenter, fyne.TextStyle{Monospace: true})
		con.AddObject(lW)
		con.AddObject(sW)
		con.AddObject(percent)
	}
	return con, sCs
}

func main() {
	x := xrandr{}
	if err := x.new(time.Millisecond * 5); err != nil {
		panic(err)
	}
	l := log.New(&bytes.Buffer{}, "", log.LUTC)
	l.SetOutput(os.Stderr)
	u, err := user.Current()
	if err != nil {
		l.Println("couldn't get current user")
		panic(err)
	}
	defaultPath := u.HomeDir + "/.bright.json"
	a := app.New()
	w := a.NewWindow("xBright")
	sliders, sCs := makeSliders(l, &x)
	w.SetContent(sliders)
	w.Resize(fyne.NewSize(400, 1))
	s := settings{Path: defaultPath}
	if err := s.fromJson(); err != nil {
		s.Path = defaultPath
		s.Preset1 = x.displays
		s.Preset2 = make(map[string]float64)
		s.Preset3 = make(map[string]float64)
		s.Refresh = time.Millisecond * 5
		if err := s.toJson(); err != nil {
			haveSettings = false
			l.Println(err.Error() + "\n" + fmt.Sprintf("couldn't make settings file at %s", s.Path))
		}
	}
	if haveSettings {
		if err := x.refresh(s.Preset1); err != nil {
			// Monitor from settings is missing.
		} else {
			for _, sC := range sCs {
				if s.Preset1[sC.name] != sC.prev/100 {
					newVal := s.Preset1[sC.name] * 100
					sC.prev = newVal
					sC.onChanged(newVal)
					sC.slider.Value = newVal
				}
			}
		}
	}
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
