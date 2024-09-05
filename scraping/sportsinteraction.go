package scraping

import (
	"context"
	"examples/webscraper/util"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"

	cdp "github.com/chromedp/chromedp"
)

func ScrapeSportsInteraction(fights *util.ThreadSafeFights, fighters *util.ThreadSafeFighters, opponents *util.ThreadSafeOpponents, wg *sync.WaitGroup) {
	// Setup driver
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
		cdp.Navigate(Urls["SportsInteraction"]),
		cdp.WaitReady(`div.participant`),
		cdp.Evaluate(`document.querySelectorAll('div.participant').length`, &numFighters),
		cdp.ActionFunc(func(ctx context.Context) error {
			for i := 0;i < numFighters-1; i += 2 {
				var nameAndCountryA, nameAndCountryB, oddsStringA, oddsStringB string 
				err := cdp.Evaluate(fmt.Sprintf(`document.querySelectorAll('div.participant')[%d].textContent`, i), &nameAndCountryA).Do(ctx)
				if err != nil { return err }
				err = cdp.Evaluate(fmt.Sprintf(`document.querySelectorAll('div.participant')[%d].textContent`, i+1), &nameAndCountryB).Do(ctx)
				if err != nil { return err }
				err = cdp.Evaluate(fmt.Sprintf(`document.querySelectorAll('div.option-value')[%d].textContent`, i), &oddsStringA).Do(ctx)
				if err != nil { return err }
				err = cdp.Evaluate(fmt.Sprintf(`document.querySelectorAll('div.option-value')[%d].textContent`, i+1), &oddsStringB).Do(ctx)
				if err != nil { return err }
				var nameEndIndexA, nameEndIndexB = strings.LastIndexFunc(nameAndCountryA, unicode.IsLower), strings.LastIndexFunc(nameAndCountryB, unicode.IsLower)
				if nameEndIndexA < 0 { return fmt.Errorf("unable to find name in %s", nameAndCountryA) }
				if nameEndIndexB < 0 { return fmt.Errorf("unable to find name in %s", nameAndCountryB) }
				var nameA, nameB = nameAndCountryA[:nameEndIndexA+1], nameAndCountryB[:nameEndIndexB+1]
				oddsA, err := strconv.ParseFloat(oddsStringA, 64)
				if err != nil { return err }
				oddsB, err := strconv.ParseFloat(oddsStringB, 64)
				if err != nil { return err }
				fighterA := &util.Fighter{
					Name: util.NewName(nameA),
					Sites: []util.SiteData{{Site: "SportsInteraction", Odds: oddsA}},
					BestSite: util.SiteData{Site: "SportsInteractin", Odds: oddsA}}
				fighterB := &util.Fighter{
					Name: util.NewName(nameB),
					Sites: []util.SiteData{{Site: "SportsInteraction", Odds: oddsB}},
					BestSite: util.SiteData{Site: "SportsInteractin", Odds: oddsB}}
				fighters.AddFighters(fighterA, fighterB)
				fights.AddFight(&util.Fight{FighterA: fighterA, FighterB: fighterB})
				opponents.AddPairing(fighterA.Name, fighterB.Name)
			}
			return nil
		})}
	
	err := cdp.Run(ctx, tasks)
	if ctx.Err() == context.DeadlineExceeded {
		cancel()
		ScrapeSportsInteraction(fights, fighters, opponents, wg)
		return 
	}
	if err != nil { panic(err) }

	cancel()
	wg.Done()
}