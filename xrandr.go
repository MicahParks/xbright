package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"regexp"
	"strconv"
	"sync"
	"time"
)

var noBright = errors.New("xrandr isn't reporting connected diplay's brightness")
var noDisplays = errors.New("no connected displays found in xrandr --query")
var notDisplay = errors.New("given string is not a connected display")
var tooBright = errors.New("trying to set value to over 100% brightness")

type xrandr struct {
	death    chan struct{}
	displays map[string]float64
	muxQ     *sync.Mutex
	queued   map[string]float64
	wait     time.Duration
}

func (x *xrandr) brightLoop() {
	for {
		select {
		case <-x.death:
			return
		case <-time.After(x.wait):
			x.muxQ.Lock()
			for k, v := range x.queued {
				if v > 1 {
					panic(tooBright)
				}
				m, err := buildMap()
				if err != nil {
					panic(err)
				}
				if _, ok := m[k]; !ok {
					panic(notDisplay)
				}
				m[k] = v
				if err = x.refresh(m); err != nil {
					panic(err)
				}
			}
			x.queued = make(map[string]float64)
			x.muxQ.Unlock()
		}
	}
}

func (x *xrandr) loadDisplays() error {
	d, err := buildMap()
	if err != nil {
		return err
	}
	x.displays = d
	return nil
}

func (x *xrandr) new(wait time.Duration) error {
	x.death = make(chan struct{})
	x.muxQ = &sync.Mutex{}
	x.wait = wait
	x.queued = make(map[string]float64)
	if _, err := exec.LookPath("xrandr"); err != nil {
		return err
	}
	err := x.loadDisplays()
	if err != nil {
		return err
	}
	go x.brightLoop()
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

func (x *xrandr) setBrightness(set string, val float64) {
	x.muxQ.Lock()
	x.queued[set] = val
	x.muxQ.Unlock()
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
	cmd := exec.Command("xrandr", "--output", display, "--brightness", fmt.Sprintf("%.2f", bright))
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
