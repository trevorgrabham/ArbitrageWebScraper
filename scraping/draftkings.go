package scraping

import (
	"context"
	"examples/webscraper/util"
	"fmt"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/chromedp/cdproto/network"
	cdp "github.com/chromedp/chromedp"
)

func ScrapeDraftKings(fights *util.ThreadSafeFights, fighters *util.ThreadSafeFighters, opponents *util.ThreadSafeOpponents, workChan chan bool) {
	// Setup the driver
	ctx, cancel := cdp.NewExecAllocator(
		context.Background(),
		cdp.ExecPath(`C:\Program Files (x86)\Google\Chrome\Application\chrome.exe`),
	)
	defer cancel()
	ctx, cancel = cdp.NewContext(ctx)
	defer cancel()
	ctx, cancel = context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	var numFighters int 
	tasks := cdp.Tasks{
		network.Enable(),
		network.SetExtraHTTPHeaders(network.Headers(util.Headers)),
		cdp.Navigate(Urls["DraftKings"]),
		cdp.WaitReady(`//tr/td[3]/div/div/div/div/div[2]/span`),
		cdp.Sleep(3*time.Second),
		cdp.Evaluate(`document.querySelectorAll('span.sportsbook-odds.no-margin').length`, &numFighters),
		cdp.ActionFunc(func(ctx context.Context) error {
			for i := 0; i < numFighters-1; i += 2 {
				var nameA, nameB, oddsStringA, oddsStringB string
				err := cdp.Evaluate(fmt.Sprintf(`document.querySelectorAll('div.event-cell__name-text')[%d].textContent`, i), &nameA).Do(ctx)
				if err != nil { return err }
				err = cdp.Evaluate(fmt.Sprintf(`document.querySelectorAll('span.sportsbook-odds.no-margin')[%d].textContent`, i), &oddsStringA).Do(ctx)
				if err != nil { return err }
				err = cdp.Evaluate(fmt.Sprintf(`document.querySelectorAll('div.event-cell__name-text')[%d].textContent`, i+1), &nameB).Do(ctx)
				if err != nil { return err }
				err = cdp.Evaluate(fmt.Sprintf(`document.querySelectorAll('span.sportsbook-odds.no-margin')[%d].textContent`, i+1), &oddsStringB).Do(ctx)
				if err != nil { return err }
				var oddsA, oddsB float64
				switch sign, size := utf8.DecodeRuneInString(oddsStringA); {
				case sign == '+':
					oddsA, err = strconv.ParseFloat(oddsStringA[size:], 64)
					oddsA = americanToDecimalOdds(true, oddsA)
					if err != nil { return err }
				default:
					oddsA, err = strconv.ParseFloat(oddsStringA[size:], 64)
					oddsA = americanToDecimalOdds(false, oddsA)
					if err != nil { return err }
				}
				switch sign, size := utf8.DecodeRuneInString(oddsStringB); {
				case sign == '+':
					oddsB, err = strconv.ParseFloat(oddsStringB[size:], 64)
					oddsB = americanToDecimalOdds(true, oddsB)
					if err != nil { return err }
				default:
					oddsB, err = strconv.ParseFloat(oddsStringB[size:], 64)
					oddsB = americanToDecimalOdds(false, oddsB)
					if err != nil { return err }
				}
				fighterA := &util.Fighter{
					Name: util.NewName(nameA), 
					Sites: []util.SiteData{{Site: "DraftKings", Odds: oddsA }},
					BestSite: util.SiteData{Site: "DraftKings", Odds: oddsA }}
				fighterB := &util.Fighter{
					Name: util.NewName(nameB), 
					Sites: []util.SiteData{{Site: "DraftKings", Odds: oddsB }},
					BestSite: util.SiteData{Site: "DraftKings", Odds: oddsB }}
				fighters.AddFighters(fighterA, fighterB)
				fights.AddFight(&util.Fight{FighterA: fighterA, FighterB: fighterB})
				opponents.AddPairing(fighterA.Name, fighterB.Name)
			}
			return nil
		})}
	
	err := cdp.Run(ctx, tasks)
	if ctx.Err() == context.DeadlineExceeded {
		cancel()
		ScrapeDraftKings(fights, fighters, opponents, workChan)
		return
	}
	if err != nil { panic(err) }
	
	workChan <- true
	cancel()
}

func americanToDecimalOdds(isPositive bool, american float64) (decimal float64) {
	switch isPositive {
	case true:
		return 1.0 + (american/100.0)
	case false:
		return 1.0 + (100.0/american)
	}
	return -1
}