package scraping

import (
	"examples/webscraper/util"
	"sync"
)

const DEBUG = true

var Urls = map[string]string{
	"PlayNow":    "https://playnow.com/sports/sports/competition/7097/mixed-martial-arts/all-mma/ufc/matches",
	"DraftKings": "https://sportsbook.draftkings.com/leagues/mma/ufc",
	"MGM":        "https://sports.on.betmgm.ca/en/sports/mma-45/betting/usa-9",
}

var Funcs []func(*util.ThreadSafeFights, *util.ThreadSafeFighters, *sync.WaitGroup) = []func(*util.ThreadSafeFights, *util.ThreadSafeFighters, *sync.WaitGroup){
	ScrapePlayNow,
	ScrapeDraftKings,
	ScrapeMGM,
}
