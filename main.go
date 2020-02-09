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
	"fyne.io/fyne/driver/desktop"
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

func (s *slideC) onChanged(val float64) {
	if val != s.prev {
		s.prev = val
		s.l.Println(val)
		val = val / 100
		if val > 1 || val <= 0 {
			return
		}
		s.x.setBrightness(s.name, val)
		s.percent.SetText(fmt.Sprintf("%.f", val*100) + "%")
		s.percent.Refresh()
	}
}

func main() {
	x := &xrandr{}
	if err := x.new(time.Millisecond * 5); err != nil {
		panic(err)
	}
	l := log.New(&bytes.Buffer{}, "", log.LUTC)
	l.SetOutput(os.Stderr)
	u, err := user.Current()
	if err != nil {
		l.Fatalln(err.Error() + "\ncouldn't get current user")
	}
	defaultPath := u.HomeDir + "/.xbright.json"
	a := app.New()
	icon, err := fyne.LoadResourceFromPath("pics/icon.png")
	if err != nil {
		l.Fatalln("Couldn't load icon.")
	}
	a.SetIcon(icon)
	w := a.NewWindow("xBright")
	sliders, sCs := makeSliders(l, x)
	s := settings{Path: defaultPath}
	s.presets(defaultPath, l, sCs, x)
	save := false
	settings := s.settingsTab(&save, sCs, x)
	slTab := widget.NewTabItem("sliders", sliders)
	tabs := widget.NewTabContainer(slTab, settings)
	shortcut(tabs, w)
	w.Resize(fyne.NewSize(400, 1))
	w.SetContent(tabs)
	w.ShowAndRun()
	close(x.death)
}

func shortcut(tabs *widget.TabContainer, w fyne.Window) {
	ctrlTab := desktop.CustomShortcut{KeyName: fyne.KeyTab, Modifier: desktop.ControlModifier}
	w.Canvas().AddShortcut(&ctrlTab, func(shortcut fyne.Shortcut) {
		switch currentTab := tabs.CurrentTabIndex(); currentTab {
		case 0:
			tabs.SelectTabIndex(1)
		case 1:
			tabs.SelectTabIndex(0)
		}
	})
}
