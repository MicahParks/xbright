package main

import (
	"fmt"
	"log"
	"sort"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

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
