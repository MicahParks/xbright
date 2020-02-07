package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"regexp"
	"strconv"
)

var noBright = errors.New("xrandr isn't reporting connected diplay's brightness")
var noDisplays = errors.New("no connected displays found in xrandr --query")
var notDisplay = errors.New("given string is not a connected display")

type xrandr struct {
	displays map[string]float64
}

func brights(expected int, out []byte) ([]float64, error) {
	re := regexp.MustCompile(`Brightness.*`)
	bBrightnesses := re.FindAll(out, -1)
	if bBrightnesses == nil {
		return nil, noBright
	}
	brightnesses := make([]float64, 0)
	for _, bBrightness := range bBrightnesses {
		bFloat := bytes.Split(bBrightness, []byte(" "))[1]
		float, err := strconv.ParseFloat(string(bFloat), 64)
		if err != nil {
			return nil, err
		}
		brightnesses = append(brightnesses, float)
	}
	if len(brightnesses) != expected {
		return nil, noBright
	}
	return brightnesses, nil
}

func buildMap() (map[string]float64, error) {
	out, err := query()
	if err != nil {
		return nil, err
	}
	d, err := displays(out)
	brightnesses, err := brights(len(d), out)
	if err != nil {
		return nil, err
	}
	m := make(map[string]float64)
	for i, display := range d {
		m[display] = brightnesses[i]
	}
	return m, nil
}

func displays(out []byte) ([]string, error) {
	re := regexp.MustCompile(`.+\b(connected)\b.*\n`)
	bDisplays := re.FindAll(out, -1)
	if bDisplays == nil {
		return nil, noDisplays
	}
	d := make([]string, 0)
	for _, bDisplay := range bDisplays {
		d = append(d, string(bytes.Split(bDisplay, []byte(" "))[0]))
	}
	return d, nil
}

func main() {
	x := xrandr{}
	if err := x.new(); err != nil {
		panic(err)
	}
	for k, v := range x.displays {
		println(k, v)
	}
	if err := x.setBrightness("DVI-I-1", 1); err != nil {
		panic(err)
	}
}

func query() ([]byte, error) {
	cmd := exec.Command("xrandr", "--verbose", "--query")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err = cmd.Start(); err != nil {
		return nil, err
	}
	out, err := ioutil.ReadAll(stdout)
	if err != nil {
		return nil, err
	}
	if err = stdout.Close(); err != nil {
		return nil, err
	}
	return out, nil
}

func xrandrBright(display string, bright float64) error {
	b := fmt.Sprintf("%.2f", bright)
	println(b)
	cmd := exec.Command("xrandr", "--output", display, "--brightness", b)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func (x *xrandr) loadDisplays() error {
	d, err := buildMap()
	if err != nil {
		return err
	}
	x.displays = d
	return nil
}

func (x *xrandr) new() error {
	if _, err := exec.LookPath("xrandr"); err != nil {
		return err
	}
	err := x.loadDisplays()
	if err != nil {
		return err
	}
	return nil
}

func (x *xrandr) refresh(m map[string]float64) error {
	var err error
	for k, v := range m {
		if _, ok := x.displays[k]; !ok || x.displays[k] != v {
			x.displays[k] = v
			err = xrandrBright(k, v)
			if err != nil {
				return err
			}
			continue
		}
	}
	for k := range x.displays {
		if _, ok := m[k]; !ok {
			delete(x.displays, k)
		}
	}
	return nil
}

func (x *xrandr) setBrightness(set string, val float64) error {
	m, err := buildMap()
	if err != nil {
		return err
	}
	if _, ok := m[set]; !ok {
		return notDisplay
	}
	m[set] = val
	return x.refresh(m)
}
