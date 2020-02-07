package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os/exec"
	"regexp"
)

var noDisplays = errors.New("no connected displays found in xrandr --query")

type xrandr struct {
	displays []string
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
	cmd := exec.Command("xrandr", "--query")
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
	x.displays = make([]string, 0)
	for _, bDisplay := range bDisplays {
		x.displays = append(x.displays, string(bytes.Split(bDisplay, []byte(" "))[0]))
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
