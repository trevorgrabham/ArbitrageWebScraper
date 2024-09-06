package main

import (
	"fmt"

	"examples/webscraper/scraping"
	"examples/webscraper/util"
)

const DEBUG = true

var fights *util.ThreadSafeFights = &util.ThreadSafeFights{}
var fighters *util.ThreadSafeFighters = &util.ThreadSafeFighters{}
var opponents *util.ThreadSafeOpponents = &util.ThreadSafeOpponents{}

func init() {
	fights.Fights = make([]*util.Fight, 0)
	fighters.Fighters = make(map[util.Name]*util.Fighter)
	opponents.Opponents = make(map[util.Name]util.Name)
}

func main() {
	for _, f := range scraping.Funcs {
		f(fights, fighters, opponents)
	}

	fmt.Println(fights)
	fmt.Println(fighters)
	fmt.Println(util.FindArbitrageOpportunities(fights, fighters, 100))
}

// Figure out why it stopped working after adding the bet.go file? 
// Figure out why the bets are not being calculated propery. Think it has soemthing to do with not being able to get the underlying data properly in the first place