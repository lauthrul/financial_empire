package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"jinrongdiguo/util"
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
		return fmt.Errorf("invalid config: hero in use not match place")
	}
	return nil
}

func TimeCost() func() {
	start := time.Now()
	fmt.Println(start)
	return func() {
		end := time.Now()
		tc := time.Since(start)
		fmt.Printf("%v\ntime cost = %v\n", end, tc)
	}
}

func main() {

	usage := func() {
		// ./jinrongdiguo 10
		fmt.Printf("param error. \nusage: %s <thread_nums>", os.Args[0])
	}

	if len(os.Args) < 2 {
		usage()
		return
	}

	threads, err := strconv.Atoi(os.Args[1])
	if err != nil || threads < 1 {
		usage()
		return
	}

	defer TimeCost()()

	// load config
	var config Config
	err = LoadConfig("技能加成.json", &config)
	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Println(config)

	// gen all combinations within the given thread numbers
	var names []string
	for name, _ := range config.Heroes {
		names = append(names, name)
	}
	combinations := util.Combination(names, config.Seats, threads)
	// fmt.Println(combinations)

	// choose the highest score heroes from all the combinations
	var result Addition
	for _, group := range combinations {
		nameMap := map[string]bool{}
		for _, name := range group {
			nameMap[name] = true
		}

		value := 0
		for name, _ := range nameMap {
			hero := config.Heroes[name]
			for _, addition := range hero.Additions {
				match := true
				for _, h := range addition.Heroes {
					if !nameMap[h] {
						match = false
						break
					}
				}
				if match {
					value += hero.Color * 100 * (1 + addition.Value/100)
				}
			}
		}
		if value > result.Value {
			result.Value = value
			result.Heroes = group
		}
	}
	fmt.Println(result)

	// make suggestion for add new heroes and remove old heroes
	var add, remove []string
	for _, name := range result.Heroes {
		if hero, ok := config.Heroes[name]; ok && !config.IsInUse(name) {
			add = append(add, hero.Name)
		}
	}
	for name, hero := range config.Heroes {
		if !config.IsInUse(name) {
			continue
		}
		match := false
		for _, v := range result.Heroes {
			if v == name {
				match = true
				break
			}
		}
		if !match {
			remove = append(remove, hero.Name)
		}
	}
	fmt.Printf("suggestion:\n add: %v\n remove: %v\n", add, remove)
}
