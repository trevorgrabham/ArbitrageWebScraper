package scraping

import (
	"context"
	"examples/webscraper/util"
	"fmt"
	"strconv"
	"time"

	"github.com/chromedp/cdproto/network"
	cdp "github.com/chromedp/chromedp"
)

func ScrapeBet365(fights *util.ThreadSafeFights, fighters *util.ThreadSafeFighters, opponents *util.ThreadSafeOpponents, workChan chan bool) {
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
		cdp.Navigate(Urls["Bet365"]),
		cdp.WaitReady(`div[class$="_TeamName "]`),
		cdp.Sleep(3*time.Second),
		cdp.Evaluate(`document.querySelectorAll('div[class$="_TeamName "]').length`, &numFighters),
		cdp.ActionFunc(func(ctx context.Context) error {
			for i := 0;i < numFighters/2; i++ {
				var nameA, nameB, oddsStringA, oddsStringB string 
				err := cdp.Evaluate(fmt.Sprintf(`document.querySelectorAll('div[class$="_TeamName "]')[%d].textContent`, 2*i), &nameA).Do(ctx)
				if err != nil { return err }
				err = cdp.Evaluate(fmt.Sprintf(`document.querySelectorAll('div[class$="_TeamName "]')[%d].textContent`, (2*i)+1), &nameB).Do(ctx)
				if err != nil { return err }
				err = cdp.Evaluate(fmt.Sprintf(`document.querySelectorAll('span[class$="_Odds"]')[%d].textContent`, i), &oddsStringA).Do(ctx)
				if err != nil { return err }
				err = cdp.Evaluate(fmt.Sprintf(`document.querySelectorAll('span[class$="_Odds"]')[%d].textContent`, i+(numFighters/2)), &oddsStringB).Do(ctx)
				if err != nil { return err }
				oddsA, err := strconv.ParseFloat(oddsStringA, 64)
				if err != nil { return err }
				oddsB, err := strconv.ParseFloat(oddsStringB, 64)
				if err != nil { return err }
				fighterA := &util.Fighter{
					Name: util.NewName(nameA),
					Sites: []util.SiteData{{Site: "Bet365", Odds: oddsA}},
					BestSite: util.SiteData{Site: "Bet365", Odds: oddsA}}
				fighterB := &util.Fighter{
					Name: util.NewName(nameB),
					Sites: []util.SiteData{{Site: "Bet365", Odds: oddsB}},
					BestSite: util.SiteData{Site: "Bet365", Odds: oddsB}}
				fighters.AddFighters(fighterA, fighterB)
				fights.AddFight(&util.Fight{FighterA: fighterA, FighterB: fighterB})
				opponents.AddPairing(fighterA.Name, fighterB.Name)
			}
			return nil
		})}
	
	err := cdp.Run(ctx, tasks)
	if ctx.Err() == context.DeadlineExceeded {
		cancel() 
		ScrapeBet365(fights, fighters, opponents, workChan)
		return
	}
	if err != nil { panic(err) }

	cancel()
	workChan <- true
}