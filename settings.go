package main

import (
	"encoding/json"
	"io/ioutil"
	"time"
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
