package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

const (
	Red    = 5
	Orange = 4
	Purple = 3
	Blue   = 2
	Green  = 1
)

type Addition struct {
	Value  int      `json:"value"`
	Heroes []string `json:"heroes"`
}

type Hero struct {
	Name      string     `json:"name"`
	Color     int        `json:"color"`
	Additions []Addition `json:"additions"`
}

type Config struct {
	Seats  int             `json:"seats"`
	InUse  []string        `json:"in_use"`
	Heroes map[string]Hero `json:"heroes"`
}

func (this *Config) IsInUse(hero string) bool {
	for _, name := range this.InUse {
		if name == hero {
			return true
		}
	}
	return false
}

func LoadConfig(file string, config *Config) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, config)
	if err != nil {
		return err
	}
	if config.Seats != len(config.InUse) {
		return fmt.Errorf("heroes in use not match seats")
	}
	return nil
}
