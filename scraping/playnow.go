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

func ScrapePlayNow(fights *util.ThreadSafeFights, fighters *util.ThreadSafeFighters, opponents *util.ThreadSafeOpponents, wg *sync.WaitGroup) {
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

	var numLoadButtons, numFighters int

	clickButtons := cdp.Tasks{
		cdp.Navigate(Urls["PlayNow"]),
		cdp.WaitReady(`button[data-testid="load-more"]`),
		cdp.Evaluate(`document.querySelectorAll('button[data-testid="load-more"]').length`, &numLoadButtons),
		cdp.ActionFunc(func(ctx context.Context) error {
			for range numLoadButtons {
				err := cdp.Click(`button[data-testid="load-more"]`).Do(ctx)
				if err != nil { return err }
			}
			return nil
		}),
		cdp.Sleep(5*time.Second)}
	
	getFighters := cdp.Tasks{
		cdp.Evaluate(`document.querySelectorAll('span[data-testid="outcome-odds"]').length`, &numFighters),
		cdp.ActionFunc(func(ctx context.Context) error {
			var fighterA, fighterB string
			for i := 0;i < numFighters-1;i += 2 {
				err := cdp.Evaluate(fmt.Sprintf(`document.querySelectorAll('span[data-testid="outcome-odds"]')[%d].textContent`, i), &fighterA).Do(ctx)
				if err != nil { return err }
				err = cdp.Evaluate(fmt.Sprintf(`document.querySelectorAll('span[data-testid="outcome-odds"]')[%d].textContent`, i+1), &fighterB).Do(ctx)
				if err != nil { return err }
				var oddsAIndex, oddsBIndex = strings.IndexFunc(fighterA, unicode.IsDigit), strings.IndexFunc(fighterB, unicode.IsDigit)
				if oddsAIndex < 0 || oddsBIndex < 0 { return fmt.Errorf("could not find odds for one of %s or %s", fighterA, fighterB) }
				oddsA, err := strconv.ParseFloat(fighterA[oddsAIndex:], 64)
				if err != nil { return err }
				oddsB, err := strconv.ParseFloat(fighterB[oddsBIndex:], 64)
				if err != nil { return err }
				var nameA, nameB = util.NewName(fighterA[:oddsAIndex]), util.NewName(fighterB[:oddsBIndex])
				a := &util.Fighter{
					Name: nameA, 
					Sites: []util.SiteData{{Site: "PlayNow", Odds: oddsA}},
					BestSite: util.SiteData{Site: "PlayNow", Odds: oddsA}}
				b := &util.Fighter{
					Name: nameB,
					Sites: []util.SiteData{{Site: "PlayNow", Odds: oddsB}},
					BestSite: util.SiteData{Site: "PlayNow", Odds: oddsB}}
				fighters.AddFighters(a, b)
				fights.AddFight(&util.Fight{FighterA: a, FighterB: b})
				opponents.AddPairing(a.Name, b.Name)
			}
			return nil
		})}

			
	err := cdp.Run(ctx, clickButtons, getFighters)
	if ctx.Err() == context.DeadlineExceeded {
		cancel()
		ScrapePlayNow(fights, fighters, opponents, wg)
		return
	}
	if err != nil { panic(err) }

	cancel()
	wg.Done()
}
