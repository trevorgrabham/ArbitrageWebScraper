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

func ScrapeMGM(fights *util.ThreadSafeFights, fighters *util.ThreadSafeFighters, wg *sync.WaitGroup) {
	ctx, cancel := cdp.NewExecAllocator(
		context.Background(), 
		cdp.ExecPath(`C:\Program Files (x86)\Google\Chrome\Application\chrome.exe`),
	)
	defer cancel()
	ctx, cancel = cdp.NewContext(ctx)
	defer cancel()
	ctx, cancel = context.WithTimeout(ctx, time.Second * 20)
	defer cancel()

	tasks := cdp.Tasks{}
	var fighterNames []string
	var fighterOdds []float64

	tasks = append(tasks, 
		cdp.Navigate(Urls["MGM"]),
		cdp.WaitReady(`div.participant`, cdp.ByQuery),
		cdp.Sleep(5*time.Second))
	
	// Get fighter names
	// Fights of form fighterNames[i] vs fighterNames[i+1]
	getName := func(ctx context.Context, id runtime.ExecutionContextID, nodes ...*proto.Node) error {
		for _, node := range nodes {
			fighterNames = append(fighterNames, strings.TrimSpace(node.Children[0].NodeValue))
		}
		return nil
	}
	tasks = append(tasks, cdp.QueryAfter(`div.participant`, getName, cdp.ByQueryAll, cdp.Populate(2, false), cdp.NodeReady))

	// Get fighter odds 
	// Fighters of form fighterName[i] has fighterOdds[i]
	getOdds := func(ctx context.Context, id runtime.ExecutionContextID, nodes ...*proto.Node) error {
		for _, node := range nodes {
			odds, err := strconv.ParseFloat(node.Children[0].NodeValue, 64)
			if err != nil { return err }
			fighterOdds = append(fighterOdds, odds)
		}
		return nil
	}
	tasks = append(tasks, cdp.QueryAfter(`span.custom-odds-value-style`, getOdds, cdp.ByQueryAll))

	// Run tasks
	err := cdp.Run(ctx, tasks...)
	if ctx.Err() == context.DeadlineExceeded {
		cancel()
		ScrapeMGM(fights, fighters, wg)
		return
	} 
	if err != nil {
		log.Fatal(err)
	}

	// Format the data 
	for i := 0; i < len(fighterNames)-1; i += 2 {
		// Populate Fighters
		var fighterA, fighterB *util.Fighter
		site := util.SiteData{Site: "MGM", Odds: fighterOdds[i]}
		fighterA = &util.Fighter{
			Name: 		util.NewName(fighterNames[i]),
			Sites: 		[]util.SiteData{site},
			BestSite: site,
		}
		site.Odds = fighterOdds[i+1]
		fighterB = &util.Fighter{
			Name: 		util.NewName(fighterNames[i+1]),
			Sites: 		[]util.SiteData{site},
			BestSite: site,
		}
		if exists, isOp, fighter := fighters.AddFighters(fighterA, fighterB); exists {
			if isOp {
				fighterA = fights.Opponent(fighter)
				fighterB = fighter
			} else {
				fighterA = fighter 
				fighterB = fights.Opponent(fighterA)
			}
			fighters.AddOdds(fighterB.Name, site)
			site.Odds = fighterOdds[i]
			fighters.AddOdds(fighterA.Name, site)
		}

		// Populate Fight
		fights.AddFight(&util.Fight{FighterA: fighterA, FighterB: fighterB})
	}
	cancel()
	wg.Done()
}