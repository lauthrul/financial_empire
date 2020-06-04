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
	Value int      `json:"value"`
	Heros []string `json:"heros"`
}

type Hero struct {
	Color     int        `json:"color"`
	Use       int        `json:"use"`
	Additions []Addition `json:"additions"`
}

type Config struct {
	Place int             `json:"place"`
	Heros map[string]Hero `json:"heros"`
}

func LoadConfig(file string, config *Config) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, config)
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

func usage(cmd string) {

}

func main() {

	usage := func() {
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

	var config Config
	err = LoadConfig("技能加成.json", &config)
	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Println(config)

	var names []string
	for name, _ := range config.Heros {
		names = append(names, name)
	}
	combinations := util.Combination(names, config.Place, threads)
	// fmt.Println(combinations)

	var result Addition
	for _, group := range combinations {
		nameMap := map[string]bool{}
		for _, name := range group {
			nameMap[name] = true
		}

		value := 0
		for name, _ := range nameMap {
			hero := config.Heros[name]
			for _, addition := range hero.Additions {
				match := true
				for _, h := range addition.Heros {
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
			result.Heros = group
		}
	}

	fmt.Println(result)

	var add, remove []string
	for _, name := range result.Heros {
		if hero, ok := config.Heros[name]; ok && hero.Use == 0 {
			add = append(add, name)
		}
	}
	for name, hero := range config.Heros {
		if hero.Use == 0 {
			continue
		}
		match := false
		for _, v := range result.Heros {
			if v == name {
				match = true
				break
			}
		}
		if !match {
			remove = append(remove, name)
		}
	}
	fmt.Printf("suggestion:\n add: %v\n remove: %v\n", add, remove)
}
