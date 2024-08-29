package scraping

import (
	"context"
	"examples/webscraper/util"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	proto "github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/runtime"
	cdp "github.com/chromedp/chromedp"
)

func ScrapeDraftKings(fights *util.ThreadSafeFights, fighters *util.ThreadSafeFighters, wg *sync.WaitGroup) {
	ctx, cancel := cdp.NewExecAllocator(
		context.Background(),
		cdp.ExecPath(`C:\Program Files (x86)\Google\Chrome\Application\chrome.exe`),
	)
	defer cancel()
	ctx, cancel = cdp.NewContext(ctx)
	defer cancel()
	ctx, cancel = context.WithTimeout(ctx, time.Second * 20)
	defer cancel()

	tasks := cdp.Tasks{
		cdp.Navigate(Urls["DraftKings"]),
		cdp.WaitReady(`a.event-cell-link`),
	}
	var fighterNames []string
	var fighterOdds []float64

	// Get the names for the fighters
	// Fights are fighters[i] vs fighters[i+1]
	getName := func(ctx context.Context, id runtime.ExecutionContextID, nodes ...*proto.Node) error {
		for _, node := range nodes {
			fighterNames = append(fighterNames, strings.TrimSpace(node.Children[0].NodeValue))
		}
		return nil
	}
	tasks = append(tasks, cdp.QueryAfter(`div.event-cell__name-text`, getName, cdp.ByQueryAll))

	// Get the odds
	// The odds map so that fighter[i] has odds[i]
	getOdds := func(ctx context.Context, id runtime.ExecutionContextID, nodes ...*proto.Node) error {
		for _, node := range nodes {
			var americanOddsValue int
			var err error
			americanOdds := node.Children[0].NodeValue
			switch r, size := utf8.DecodeRune([]byte(americanOdds)); r {
			case '+':
				americanOddsValue, err = strconv.Atoi(americanOdds[size:])
				if err != nil { return err }
				fighterOdds = append(fighterOdds, float64(americanOddsValue)/100.0+1.0)
			default:
				americanOddsValue, err = strconv.Atoi(americanOdds[size:])
				if err != nil { return err }
				fighterOdds = append(fighterOdds, 100.0/float64(americanOddsValue)+1.0)
			}
		}
		return nil
	}
	tasks = append(tasks, cdp.QueryAfter(`span.sportsbook-odds`, getOdds, cdp.ByQueryAll))

	// Run the tasks
	err := cdp.Run(ctx, tasks...)
	if ctx.Err() == context.DeadlineExceeded {
		ScrapeDraftKings(fights, fighters, wg)
		return
	}
	if err != nil {
		log.Fatal(err)
	}

	// Format the data
	for i := 0; i < len(fighterNames)-1; i += 2 {
		// Populate the Fighters
		site := util.SiteData{Site: "DraftKings", Odds: fighterOdds[i]}
		fighterA := &util.Fighter{
			Name:     util.NewName(fighterNames[i]),
			Sites:    []util.SiteData{site},
			BestSite: site,
		}
		site.Odds = fighterOdds[i+1]
		fighterB := &util.Fighter{
			Name:     util.NewName(fighterNames[i+1]),
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
			site.Odds = fighterOdds[i]
			fighters.AddOdds(fighterA.Name, site)
		}

		// Populate the Fight
		fights.AddFight(&util.Fight{FighterA: fighterA, FighterB: fighterB})
	}
	cancel()
	wg.Done()
}
