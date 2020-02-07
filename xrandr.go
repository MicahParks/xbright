package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os/exec"
	"regexp"
	"strconv"
)

var noDisplays = errors.New("no connected displays found in xrandr --query")
var notDisplay = errors.New("given string is not a connected display")

type xrandr struct {
	displays map[string]float64
}

func main() {
	x := xrandr{}
	if err := x.new(); err != nil {
		panic(err)
	}
	for _, display := range x.displays {
		println(display)
	}
}

func (x *xrandr) getDisplays() error {
	cmd := exec.Command("xrandr", "--verbose", "--query")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err = cmd.Start(); err != nil {
		return err
	}
	re := regexp.MustCompile(`.+\b(connected)\b.*\n`)
	out, err := ioutil.ReadAll(stdout)
	if err != nil {
		return err
	}
	if err = stdout.Close(); err != nil {
		return err
	}
	bDisplays := re.FindAll(out, -1)
	if bDisplays == nil {
		return noDisplays
	}
	displays := make([]string, 0)
	for _, bDisplay := range bDisplays {
		displays = append(displays, string(bytes.Split(bDisplay, []byte(" "))[0]))
	}
	re = regexp.MustCompile(`Brightness.*`)
	bBrightnesses := re.FindAll(out, -1)
	if bBrightnesses == nil {
		// TODO Error
	}
	brightnesses := make([]float64, 0)
	for _, bBrightness := range bBrightnesses {
		bFloat := bytes.Split(bBrightness, []byte(" "))[1]
		float, err := strconv.ParseFloat(string(bFloat), 64)
		if err != nil {
			return err
		}
		brightnesses = append(brightnesses, float)
	}
	if len(brightnesses) != len(displays) {
		// TODO Error
	}
	for i, display := range displays {
		x.displays[display] = brightnesses[i]
	}
	return nil
}

func (x *xrandr) new() error {
	if _, err := exec.LookPath("xrandr"); err != nil {
		return err
	}
	err := x.getDisplays()
	if err != nil {
		return err
	}
	return nil
}

func (x *xrandr) setBrightness(set string, val float64) error {
	good := false
	for _, display := range x.displays {
		if display == set {
			good = true
			break
		}
	}
	if !good {
		return notDisplay
	}

	return nil
}
