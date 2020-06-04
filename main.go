package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"jinrongdiguo/config"
	"jinrongdiguo/util"
)

func timeCost() func() {
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

	defer timeCost()()

	// load config
	var cfg config.Config
	err = config.LoadConfig("技能加成.json", &cfg)
	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Println(cfg)

	// gen all combinations within the given thread numbers
	var names []string
	for name, _ := range cfg.Heroes {
		names = append(names, name)
	}
	combinations := util.Combination(names, cfg.Seats, threads)
	// fmt.Println(combinations)

	// choose the highest score heroes from all the combinations
	var result config.Addition
	for _, group := range combinations {
		nameMap := map[string]bool{}
		for _, name := range group {
			nameMap[name] = true
		}

		value := 0
		for name, _ := range nameMap {
			hero := cfg.Heroes[name]
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
		if hero, ok := cfg.Heroes[name]; ok && !cfg.IsInUse(name) {
			add = append(add, hero.Name)
		}
	}
	for name, hero := range cfg.Heroes {
		if !cfg.IsInUse(name) {
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
