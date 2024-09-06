package scraping

import (
	"context"
	"examples/webscraper/util"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/cdproto/network"
	cdp "github.com/chromedp/chromedp"
)

func ScrapeTonyBet(fights *util.ThreadSafeFights, fighters *util.ThreadSafeFighters, opponents *util.ThreadSafeOpponents, workChan chan bool) {
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
		cdp.Navigate(Urls["TonyBet"]),
		cdp.WaitReady(`div[data-test="teamName"] > div > span`),
		cdp.Sleep(3*time.Second),
		cdp.Evaluate(`document.querySelectorAll('div[data-test="teamName"] > div > span').length`, &numFighters),
		cdp.ActionFunc(func(ctx context.Context) error {
			for i := 0;i < numFighters-1; i += 2 {
				var nameA, nameB, oddsStringA, oddsStringB string 
				err := cdp.Evaluate(fmt.Sprintf(`document.querySelectorAll('div[data-test="teamName"] > div > span')[%d].textContent`, i), &nameA).Do(ctx)
				if err != nil { return err }
				err = cdp.Evaluate(fmt.Sprintf(`document.querySelectorAll('div[data-test="teamName"] > div > span')[%d].textContent`, i+1), &nameB).Do(ctx)
				if err != nil { return err }
				err = cdp.Evaluate(fmt.Sprintf(`document.querySelectorAll('div[data-test="marketItem"]:nth-of-type(1) span')[%d].textContent`, i), &oddsStringA).Do(ctx)
				if err != nil { return err }
				err = cdp.Evaluate(fmt.Sprintf(`document.querySelectorAll('div[data-test="marketItem"]:nth-of-type(1) span')[%d].textContent`, i+1), &oddsStringB).Do(ctx)
				if err != nil { return err }
				oddsA, err := strconv.ParseFloat(oddsStringA, 64)
				if err != nil { return err }
				oddsB, err := strconv.ParseFloat(oddsStringB, 64)
				if err != nil { return err }
				var splitNameA, splitNameB = strings.Split(nameA, ", "), strings.Split(nameB, ", ")
				if len(splitNameA) < 2 {
					nameA = splitNameA[0]
				} else {
					nameA = strings.Join([]string{splitNameA[1], splitNameA[0]}, " ")
				}
				if len(splitNameB) < 2 {
					nameB = splitNameB[0]
				} else {
					nameB = strings.Join([]string{splitNameB[1], splitNameB[0]}, " ")
				}
				fighterA := &util.Fighter{
					Name: util.NewName(nameA),
					Sites: []util.SiteData{{Site: "TonyBet", Odds: oddsA}},
					BestSite: util.SiteData{Site: "TonyBet", Odds: oddsA}}
				fighterB := &util.Fighter{
					Name: util.NewName(nameB),
					Sites: []util.SiteData{{Site: "TonyBet", Odds: oddsB}},
					BestSite: util.SiteData{Site: "TonyBet", Odds: oddsB}}
				fighters.AddFighters(fighterA, fighterB)
				fights.AddFight(&util.Fight{FighterA: fighterA, FighterB: fighterB})
				opponents.AddPairing(fighterA.Name, fighterB.Name)
			}
			return nil
		})}

	err := cdp.Run(ctx, tasks)
	if ctx.Err() == context.DeadlineExceeded {
		cancel()
		ScrapeTonyBet(fights, fighters, opponents, workChan)
		return
	}
	if err != nil { panic(err) }

	workChan <- true
	cancel()
}