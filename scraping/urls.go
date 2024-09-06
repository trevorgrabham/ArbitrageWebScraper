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
	"Bet365": "https://www.bet365.com/#/AC/B162/C20860037/D1/E162148/F2/",
	"Pinnacle": "https://www.pinnacle.com/en/mixed-martial-arts/ufc/matchups",
}

var Funcs []func(*util.ThreadSafeFights, *util.ThreadSafeFighters, *util.ThreadSafeOpponents, chan bool) = []func(*util.ThreadSafeFights, *util.ThreadSafeFighters, *util.ThreadSafeOpponents, chan bool){
	ScrapePlayNow,
	ScrapeDraftKings,
	ScrapeMGM,
	ScrapeTonyBet,
	ScrapeBet365,
	ScrapePinnacle,

	// Getting blocked for some reason
	// ScrapeSportsInteraction,
}
