package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/user"
	"sort"
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
	slider  *widget.Slider
	x       *xrandr
}

func loadPreset(m *map[string]float64, sCs []*slideC, x *xrandr) error {
	if err := x.refresh(*m); err != nil {
		return err
	} else {
		for _, sC := range sCs {
			if (*m)[sC.name] != sC.prev/100 {
				newVal := (*m)[sC.name] * 100
				sC.onChanged(newVal)
				sC.slider.Value = newVal
			}
		}
	}
	return nil
}

func makeSliders(l *log.Logger, x *xrandr) (*fyne.Container, []*slideC) {
	sCs := make([]*slideC, 0)
	con := fyne.NewContainerWithLayout(layout.NewGridLayout(3))
	str := make([]string, 0)
	for k := range x.displays {
		str = append(str, k)
	}
	sort.Strings(str)
	for _, k := range str {
		v := x.displays[k]
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
		l.Fatalln(err.Error() + "\ncouldn't get current user")
	}
	defaultPath := u.HomeDir + "/.bright.json"
	a := app.New()
	w := a.NewWindow("xBright")
	sliders, sCs := makeSliders(l, &x)
	s := settings{Path: defaultPath}
	if err := s.fromJson(); err != nil || s.DefaultPreset == nil || s.Refresh == 0 || s.Path == "" {
		s.Path = defaultPath
		s.DefaultPreset = x.displays
		s.Preset2 = make(map[string]float64)
		s.Preset3 = make(map[string]float64)
		s.Refresh = time.Millisecond * 5
		if err := s.toJson(); err != nil {
			l.Fatalln(err.Error() + fmt.Sprintf("\ncouldn't make settings file at %s", s.Path))
		}
	}
	if err := loadPreset(&s.DefaultPreset, sCs, &x); err != nil {
		// Ignore error.
	}
	save := false
	radioSwitch := func(s string) {
		switch s {
		case "load":
			save = false
		case "save":
			save = true
		}
	}
	radios := widget.NewRadio([]string{"load", "save"}, radioSwitch)
	buttons := widget.NewVBox(
		widget.NewButton("Default", func() {
			if save {
				if err := savePreset(&s.DefaultPreset, &s, &x); err != nil {
					log.Fatalln(err)
				}
			} else {
				if err := loadPreset(&s.DefaultPreset, sCs, &x); err != nil {
					log.Fatalln(err)
				}
			}
		}),
		widget.NewButton("preset 2", func() {
			if save {
				if err := savePreset(&s.Preset2, &s, &x); err != nil {
					log.Fatalln(err)
				}
			} else {
				if err := loadPreset(&s.Preset2, sCs, &x); err != nil {
					log.Fatalln(err)
				}
			}
		}),
		widget.NewButton("preset 3", func() {
			if save {
				if err := savePreset(&s.Preset3, &s, &x); err != nil {
					log.Fatalln(err)
				}
			} else {
				if err := loadPreset(&s.Preset3, sCs, &x); err != nil {
					log.Fatalln(err)
				}
			}
		}),
	)
	slTab := widget.NewTabItem("sliders", sliders)
	stTab := fyne.NewContainerWithLayout(layout.NewGridLayout(2), radios, buttons)
	settings := widget.NewTabItem("settings", stTab)
	tabs := widget.NewTabContainer(slTab, settings)
	w.Resize(fyne.NewSize(400, 1))
	w.SetContent(tabs)
	w.ShowAndRun()
	close(x.death)
}

func savePreset(m *map[string]float64, s *settings, x *xrandr) error {
	*m = make(map[string]float64)
	for k, v := range x.displays {
		(*m)[k] = v
	}
	if err := s.toJson(); err != nil {
		return err
	}
	return nil
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
