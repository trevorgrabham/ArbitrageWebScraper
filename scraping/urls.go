package scraping

import (
	"examples/webscraper/util"
)

const DEBUG = true

var Urls = map[string]string{
	"PlayNow":    "https://playnow.com/sports/sports/competition/7097/mixed-martial-arts/all-mma/ufc/matches",
	"DraftKings": "https://sportsbook.draftkings.com/leagues/mma/ufc",
	"MGM":        "https://sports.on.betmgm.ca/en/sports/mma-45/betting/usa-9",
	"SportsInteraction": "https://sports.sportsinteraction.com/en-ca/sports/mma-45/betting/usa-9",
	"TonyBet":	"https://tonybet.com/ca/prematch/mma",
}

var Funcs []func(*util.ThreadSafeFights, *util.ThreadSafeFighters, *util.ThreadSafeOpponents) = []func(*util.ThreadSafeFights, *util.ThreadSafeFighters, *util.ThreadSafeOpponents){
	ScrapePlayNow,
	ScrapeDraftKings,
	ScrapeMGM,
	ScrapeTonyBet,
	ScrapeSportsInteraction,
}
