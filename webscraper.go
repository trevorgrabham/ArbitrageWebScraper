package main

import (
	"fmt"
	"sync"

	"examples/webscraper/scraping"
	"examples/webscraper/util"
)

const DEBUG = true

var fights *util.ThreadSafeFights = &util.ThreadSafeFights{}
var fighters *util.ThreadSafeFighters = &util.ThreadSafeFighters{}

func init() {
	fights.Fights = make(map[util.Name]*util.Fight)
	fighters.Fighters = make(map[util.Name]*util.Fighter)
}

func main() {
	var wg sync.WaitGroup
	for _, f := range scraping.Funcs {
		wg.Add(1)
		go f(fights, fighters, &wg)
	}
	wg.Wait()

	fmt.Println(fights)
	fmt.Println(fighters)
	fmt.Println(util.FindArbitrageOpportunities(fights, fighters, 100))
}

// Figure out why it stopped working after adding the bet.go file? 
// Figure out why the bets are not being calculated propery. Think it has soemthing to do with not being able to get the underlying data properly in the first place