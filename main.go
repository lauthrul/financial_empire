package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"financial_empire/cache"
	"financial_empire/config"
	"financial_empire/logger"
	"financial_empire/util"
)

const USE_CACHE = false

func timeCost() func() {
	start := time.Now()
	return func() {
		tc := time.Since(start)
		logger.Logger.Printf("time cost = %v\n", tc)
	}
}

func showResult(result config.Addition, cfg config.Config) {
	show := config.Addition{
		Value: result.Value,
	}
	for _, h := range result.Heroes {
		for name, hero := range cfg.Heroes {
			if h == name {
				show.Heroes = append(show.Heroes, hero.Name)
			}
		}
	}
	logger.Logger.Println(show)
}

func main() {

	usage := func() {
		// ./financial_empire 10
		fmt.Printf("usage: %s <config_file> <thread_nums>\n", os.Args[0])
	}

	if len(os.Args) < 3 {
		usage()
		return
	}

	configFile := os.Args[1]
	threads, err := strconv.Atoi(os.Args[2])
	if err != nil || threads < 1 {
		usage()
		return
	}

	defer timeCost()()

	// load config
	logger.Logger.Printf("load config ...")
	var cfg config.Config
	err = config.LoadConfig(configFile, &cfg)
	if err != nil {
		logger.Logger.Printf("load config fail: %s\n", err.Error())
		return
	}
	// log.Logger.Println(cfg)

	// gen all combinations within the given thread numbers
	logger.Logger.Printf("calculate combinations ...")
	var names []string
	for name, _ := range cfg.Heroes {
		names = append(names, name)
	}

	var combinations [][]string
	var cch cache.Cache

	if USE_CACHE {
		err = cache.LoadCache(&cch)
		if err != nil {
			logger.Logger.Printf("load cache fail: %s\n", err.Error())
		} else {
			if cch.Exist(names) {
				logger.Logger.Printf("read cache\n")
				err = cch.Read(names, &combinations)
				if err != nil {
					logger.Logger.Printf("read cache fail: %s\n", err.Error())
				}
			}
		}
	}

	if len(combinations) <= 0 {
		util.Combination(names, uint64(cfg.Seats), uint64(threads), &combinations)
		if USE_CACHE {
			err = cch.Save(names, combinations)
			if err != nil {
				logger.Logger.Printf("save cache fail: %s\n", err.Error())
			}
		}
	}
	// log.Logger.Println(combinations)

	// choose the highest score heroes from all the combinations
	logger.Logger.Printf("choose best solution ...")
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
					value += int(math.Pow(float64(hero.Color), float64(10))) * (1 + addition.Value/100)
				}
			}
		}
		if value > result.Value {
			result.Value = value
			result.Heroes = group
		}
	}
	showResult(result, cfg)

	// make suggestion for add new heroes and remove old heroes
	logger.Logger.Printf("make suggestion ...")
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
	logger.Logger.Printf("suggestion:\n add: %v\n remove: %v\n", add, remove)
}
