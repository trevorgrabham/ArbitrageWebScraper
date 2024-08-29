package scraping

import (
	"context"
	"examples/webscraper/util"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	proto "github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/runtime"
	cdp "github.com/chromedp/chromedp"
)

func ScrapePlayNow(fights *util.ThreadSafeFights, fighters *util.ThreadSafeFighters, wg *sync.WaitGroup) {
	ctx, cancel := cdp.NewExecAllocator(
		context.Background(),
		cdp.ExecPath(`C:\Program Files (x86)\Google\Chrome\Application\chrome.exe`),
	)
	defer cancel()
	ctx, cancel = cdp.NewContext(ctx)
	defer cancel()
	ctx, cancel = context.WithTimeout(ctx, time.Second*20)
	defer cancel()

	tasks := cdp.Tasks{}
	var nodes []*proto.Node
	var fighterNames []string
	var fighterOdds []float64

	// Check how many "Show More Events" buttons there are and click them all
	cdp.Run(ctx,
		cdp.Navigate(Urls["PlayNow"]),
		cdp.WaitReady(`span[class^="outcomeOddsCommon-"]`, cdp.NodeReady),
		cdp.Nodes(`button[data-testid="load-more"]`, &nodes, cdp.ByQueryAll))
	for range nodes {
		tasks = append(tasks, cdp.Click(`button[data-testid="load-more"]`, cdp.ByQuery))
	}
	tasks = append(tasks, cdp.Sleep(5*time.Second))

	// Get the fighter names for all of the fights
	getFight := func(ctx context.Context, id runtime.ExecutionContextID, nodes ...*proto.Node) error {
		for _, node := range nodes {
			fighterNames = append(fighterNames, strings.TrimSpace(node.Children[0].NodeValue))
		}
		return nil
	}
	tasks = append(tasks,
		cdp.QueryAfter(`div[data-testid="event-card-team-name-a"]`, getFight, cdp.ByQueryAll),
		cdp.QueryAfter(`div[data-testid="event-card-team-name-b"]`, getFight, cdp.ByQueryAll))

	// Get the name and odds for each fighter
	getOdds := func(ctx context.Context, id runtime.ExecutionContextID, nodes ...*proto.Node) error {
		for _, node := range nodes {
			odds, err := strconv.ParseFloat(node.Children[0].NodeValue, 64)
			if err != nil { return err }
			fighterOdds = append(fighterOdds, odds)
		}
		return nil
	}
	tasks = append(tasks, cdp.QueryAfter(`span[class^="outcomePriceCommon-"]`, getOdds, cdp.ByQueryAll))

	// Run the tasks
	err := cdp.Run(ctx, tasks...)
	if ctx.Err() == context.DeadlineExceeded {
		ScrapePlayNow(fights, fighters, wg)
		return
	}
	if err != nil {
		log.Fatal(err)
	}

	// Format the data
	midpoint := len(fighterNames) / 2
	for i := 0; i < midpoint; i++ {
		// Populate FighterA and FighterB
		var fighterA, fighterB *util.Fighter
		site := util.SiteData{Site: "PlayNow", Odds: fighterOdds[2*i]}
		fighterA = &util.Fighter{
			Name:     util.NewName(fighterNames[i]),
			Sites:    []util.SiteData{site},
			BestSite: site,
		}
		site.Odds = fighterOdds[(2*i)+1]
		fighterB = &util.Fighter{
			Name:     util.NewName(fighterNames[i+midpoint]),
			Sites:    []util.SiteData{site},
			BestSite: site,
		}
		if exitsts, isOp, fighter := fighters.AddFighters(fighterA, fighterB); exitsts {
			if isOp {
				fighterA = fights.Opponent(fighter)
				fighterB = fighter
			} else {
				fighterA = fighter
				fighterB = fights.Opponent(fighterA)
			}
			fighters.AddOdds(fighterB.Name, site)
			site.Odds = fighterOdds[2*i]
			fighters.AddOdds(fighterA.Name, site)
		}

		// Populate fights
		f := &util.Fight{FighterA: fighterA, FighterB: fighterB}
		fights.AddFight(f)
	}
	cancel()
	wg.Done()
}
