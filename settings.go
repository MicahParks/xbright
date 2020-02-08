package main

import (
	"encoding/json"
	"io/ioutil"
	"time"
)

//func main() {
//	t := time.Second
//	a := settings{
//		Path:    settingsPath,
//		Preset1: make(map[string]float64),
//		Preset2: make(map[string]float64),
//		Preset3: make(map[string]float64),
//		Refresh: &t,
//	}
//	a.Preset1["dsaf}"] = 3.4
//	a.Preset2["dfalkjsdlfj"] = 93.3
//	a.Preset3["dfkajs"] = 3.1
//	if err := a.toJson(); err != nil {
//		panic(err)
//	}
//}

type settings struct {
	Path    string
	Preset1 map[string]float64
	Preset2 map[string]float64
	Preset3 map[string]float64
	Refresh time.Duration
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
