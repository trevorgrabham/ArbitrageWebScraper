package scraping

import (
	"context"
	"examples/webscraper/util"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/chromedp/cdproto/network"
	cdp "github.com/chromedp/chromedp"
)

func ScrapeMGM(fights *util.ThreadSafeFights, fighters *util.ThreadSafeFighters, opponents *util.ThreadSafeOpponents) {
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
		cdp.Navigate(Urls["MGM"]),
		cdp.WaitReady(`div.participant`),
		cdp.Sleep(3*time.Second),
		cdp.Evaluate(`document.querySelectorAll('div.participant').length`, &numFighters),
		cdp.ActionFunc(func(ctx context.Context) error {
			for i := 0; i < numFighters-1; i += 2 {
				var nameAndCountryA, nameAndCountryB, oddsStringA, oddsStringB string 
				err := cdp.Evaluate(fmt.Sprintf(`document.querySelectorAll('div.participant')[%d].textContent`, i), &nameAndCountryA).Do(ctx)
				if err != nil { return err }
				err = cdp.Evaluate(fmt.Sprintf(`document.querySelectorAll('div.participant')[%d].textContent`, i+1), &nameAndCountryB).Do(ctx)
				if err != nil { return err }
				err = cdp.Evaluate(fmt.Sprintf(`document.querySelectorAll('span.custom-odds-value-style')[%d].textContent`, i), &oddsStringA).Do(ctx)
				if err != nil { return err }
				err = cdp.Evaluate(fmt.Sprintf(`document.querySelectorAll('span.custom-odds-value-style')[%d].textContent`, i+1), &oddsStringB).Do(ctx)
				if err != nil { return err }
				nameEndIndexA := strings.LastIndexFunc(nameAndCountryA, unicode.IsLower)
				if nameEndIndexA < 0 { return fmt.Errorf("could not find where name ended in %s", nameAndCountryA) }
				nameEndIndexB := strings.LastIndexFunc(nameAndCountryB, unicode.IsLower)
				if nameEndIndexB < 0 { return fmt.Errorf("could not find where name ended in %s", nameAndCountryB) }
				oddsA, err := strconv.ParseFloat(oddsStringA, 64)
				if err != nil { return err }
				oddsB, err := strconv.ParseFloat(oddsStringB, 64)
				if err != nil { return err }
				var nameA, nameB = nameAndCountryA[:nameEndIndexA+1], nameAndCountryB[:nameEndIndexB+1]
				nameA, nameB = strings.TrimSpace(nameA), strings.TrimSpace(nameB)
				fighterA := &util.Fighter{
					Name: util.NewName(nameA),
					Sites: []util.SiteData{{Site: "MGM", Odds: oddsA}},
					BestSite: util.SiteData{Site: "MGM", Odds: oddsA}}
				fighterB := &util.Fighter{
					Name: util.NewName(nameB),
					Sites: []util.SiteData{{Site: "MGM", Odds: oddsB}},
					BestSite: util.SiteData{Site: "MGM", Odds: oddsB}}
				fighters.AddFighters(fighterA, fighterB)
				fights.AddFight(&util.Fight{FighterA: fighterA, FighterB: fighterB})
				opponents.AddPairing(fighterA.Name, fighterB.Name)
			}
			return nil
		})}
	
	err := cdp.Run(ctx, tasks)
	if ctx.Err() == context.DeadlineExceeded {
		cancel() 
		ScrapeMGM(fights, fighters, opponents)
		return
	}
	if err != nil { panic(err) }

	cancel()
}