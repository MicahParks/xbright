package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

type settings struct {
	Path          string
	DefaultPreset map[string]float64
	Preset2       map[string]float64
	Preset3       map[string]float64
	Refresh       time.Duration
}

func (s *settings) fromJson() error {
	b, err := ioutil.ReadFile(s.Path)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(b, s); err != nil {
		return err
	}
	return nil
}

func (s *settings) presets(defaultPath string, l *log.Logger, sCs []*slideC, x *xrandr) {
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
	if err := loadPreset(&s.DefaultPreset, sCs, x); err != nil {
		// Ignore error.
	}
}

func (s *settings) settingsTab(save *bool, sCs []*slideC, x *xrandr) *widget.TabItem {
	buttons := widget.NewVBox(
		widget.NewButton("default", func() {
			if *save {
				if err := savePreset(&s.DefaultPreset, s, x); err != nil {
					log.Fatalln(err)
				}
			} else {
				if err := loadPreset(&s.DefaultPreset, sCs, x); err != nil {
					log.Fatalln(err)
				}
			}
		}),
		widget.NewButton("preset 2", func() {
			if *save {
				if err := savePreset(&s.Preset2, s, x); err != nil {
					log.Fatalln(err)
				}
			} else {
				if err := loadPreset(&s.Preset2, sCs, x); err != nil {
					log.Fatalln(err)
				}
			}
		}),
		widget.NewButton("preset 3", func() {
			if *save {
				if err := savePreset(&s.Preset3, s, x); err != nil {
					log.Fatalln(err)
				}
			} else {
				if err := loadPreset(&s.Preset3, sCs, x); err != nil {
					log.Fatalln(err)
				}
			}
		}),
	)
	radios := widget.NewRadio([]string{"load", "save"}, func(s string) {
		switch s {
		case "load":
			*save = false
		case "save":
			*save = true
		}
	})
	radios.SetSelected("load")
	stTab := fyne.NewContainerWithLayout(layout.NewGridLayout(2), radios, buttons)
	return widget.NewTabItem("settings", stTab)
}

func (s *settings) toJson() error {
	b, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(s.Path, b, 0600); err != nil {
		return err
	}
	return nil
}

func loadPreset(m *map[string]float64, sCs []*slideC, x *xrandr) error {
	if err := x.refresh(*m); err != nil {
		return err
	}
	for _, sC := range sCs {
		if (*m)[sC.name] != sC.prev/100 {
			newVal := (*m)[sC.name] * 100
			sC.onChanged(newVal)
			sC.slider.Value = newVal
		}
	}
	return nil
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
