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

func ScrapeSportsInteraction(fights *util.ThreadSafeFights, fighters *util.ThreadSafeFighters, opponents *util.ThreadSafeOpponents) {
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

	headers := util.Headers 
	headers["Cookie"] = `tdpehcd=covers_roc_nosweat; btagcd=; vn-ld-session=-1725518700; tdpeh=covers_roc_nosweat; trc.cid=b76bdbce2c9d4bea95495f36f19a752e; _rdt_uuid=1725518717825.6bb6771c-ea1d-445e-bf64-5cb4d9d33e2d; _scid=S2qfN5M3ltUTG_SSVogKYGd_W7BgUwwv; _scid_r=S2qfN5M3ltUTG_SSVogKYGd_W7BgUwwv; _fbp=fb.1.1725518717956.110973966380227184; __qca=P0-1509651890-1725518717947; _tt_enable_cookie=1; _ttp=uu-Ic5bCZCbWn8CdENCod641w4P; _uetsid=69f197c06b5211ef82546f00fc624792|1en55h6|2|fox|0|1709; _ScCbts=%5B%5D; _uetvid=69f18d506b5211efb5c70b08d091aaea|9fwk0d|1725518718526|1|1|bat.bing.com/p/insights/c/x; _gcl_au=1.1.651626615.1725518719; _ga=GA1.1.949022709.1725518719; _sctr=1%7C1725433200000; tdpeh2=covers_roc_nosweat; trackerId2=5481549; DAPROPS="bS%3A0%7CscsVersion%3Abug-DA-6298-add-option-to-encode-cs-string-SNAPSHOT%7CsdeviceAspectRatio%3A16%2F9%7CsdevicePixelRatio%3A1%7Cbhtml.video.ap4x%3A0%7Cbhtml.video.av1%3A1%7Cbjs.deviceMotion%3A1%7Csjs.webGlRenderer%3AANGLE%20(NVIDIA%2C%20NVIDIA%20GeForce%20GT%20710%20(0x0000128B)%20Direct3D11%20vs_5_0%20ps_5_0%2C%20D3D11)%7CsrendererRef%3A01633483327%7CsscreenWidthHeight%3A1920%2F1080%7CsaudioRef%3A781311942%7CbE%3A0"; lastKnownProduct=%7B%22url%22%3A%22https%3A%2F%2Fsports.sportsinteraction.com%2Fen-ca%22%2C%22name%22%3A%22sports%22%2C%22previous%22%3A%22unknown%22%2C%22platformProductId%22%3A%22SPORTSBOOK%22%7D; LPVID=A0ZWVkMzA3NGYxM2RhY2Qx; _hjSessionUser_929373=eyJpZCI6ImEzYzdkNWNlLTE1MDctNTk3ZC05NDFiLWZiNWY0NDhmMjYyYiIsImNyZWF0ZWQiOjE3MjU1MTg3MTgyMTksImV4aXN0aW5nIjp0cnVlfQ==; hq=%5B%7B%22name%22%3A%22homescreen%22%2C%22shouldShow%22%3Afalse%7D%5D; _sp_id.187c=6b9e792b-3b7c-4810-8faa-01eaf5524fb5.1725518718.1.1725518784.1725518718.31e4c9b6-2348-4a98-a2a4-1f15feeaa8cc; _ga_KFTYM8CSCC=GS1.1.1725518717.1.1.1725518784.56.0.0; _ga_SM5BJ4XV8X=GS1.1.1725518717.1.1.1725518784.58.0.0; RT="z=1&dm=sportsinteraction.com&si=c3b88629-37ba-4fd8-ab53-d58e338c479d&ss=m0oxam4b&sl=0&tt=0&bcn=%2F%2F68794912.akstat.io%2F"; OptanonConsent=isIABGlobal=false&datestamp=Thu+Sep+05+2024+00%3A36%3A35+GMT-0700+(Pacific+Daylight+Time)&version=6.14.0&hosts=&consentId=22fb97fe-5be5-407a-8b4c-12567024904f&interactionCount=1&landingPath=https%3A%2F%2Fsports.sportsinteraction.com%2Fen-ca%2Fsports%3FproductId%3DSPORTSBOOK%26rurl%3Dhttps%3A%252F%252Fsports.sportsinteraction.com%252Fen-ca%252Fsports%26wm%3D5481549%26tdpeh%3Dcovers_roc_nosweat&groups=C0001%3A1%2CC0004%3A1%2CC0002%3A1%2CC0003%3A1; lang=en-ca; trackerId=5418580; seoLandingUrl=http%3A%2F%2Fsports.sportsinteraction.com%2Fen-ca%2Fsports; vnSession=85887d46-579e-4717-95aa-fcf2c23ad324; usersettings=cid%3Den-CA%26vc%3D2%26sst%3D2024-09-06T01%3A08%3A52.9039393Z%26psst%3D2024-09-05T06%3A45%3A16.5765226Z; __cf_bm=X.l.xDCnZVKLXvvETqgqzJFiVFQeMvrmmxmlhl0ljbQ-1725584933-1.0.1.1-9BNs99ASH3aN1fL2Y6GzwbuSPEOwxaPr9pDZEeV1EEN6R7Fghnds29IAM6XuwAnwB4oti4QO9OPbupaD1sg0Qw; cf_clearance=SKCnJvxaDL1ZXa9N6zltwWFTTd4.Nx2oBHl_Vc1HegQ-1725584934-1.2.1.1-_M9trblvdsDZdtlsf9IUwfA33Gi0rCCrZWMPJ6FxaUtyxmN.wuqYm713VFomJn2GaGsDpoLK4xGERshvZUbDVIUcgUHuwL7XTtSdHqYx1fDHanWBvEyXOXqi2COEIf6rMjZW4Z4zikBpiSBkBnSXF2kLZLbwP836Ghd11R9GajHg7viYmPJgNJnjEZJ3mpzNcfJxMjgIrC3549vIM6QdbpfW.sRf2oUxnvxwY98ZvvD_pSeJIl90_nHiOERRzSYJkJvN26Sv7rlD4.Fgj8fF_Khhlie6JXaOVfzyDXTSg0FkO8vWWphHmNe2g0icBWfPmJ0rNC.sXO17aObPyGuwrAfFufMFdI7_a7xieXHFMj2E1UhG4NeWNYxZfDLfGrtZ; tq=%5B%5D`

	tasks := cdp.Tasks{
		network.Enable(),
		network.SetExtraHTTPHeaders(network.Headers(headers)),
		cdp.Navigate(Urls["SportsInteraction"]),
		cdp.WaitReady(`div.participant`),
		cdp.Sleep(3*time.Second),
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
					BestSite: util.SiteData{Site: "SportsInteraction", Odds: oddsA}}
				fighterB := &util.Fighter{
					Name: util.NewName(nameB),
					Sites: []util.SiteData{{Site: "SportsInteraction", Odds: oddsB}},
					BestSite: util.SiteData{Site: "SportsInteraction", Odds: oddsB}}
				fighters.AddFighters(fighterA, fighterB)
				fights.AddFight(&util.Fight{FighterA: fighterA, FighterB: fighterB})
				opponents.AddPairing(fighterA.Name, fighterB.Name)
			}
			return nil
		})}
	
	err := cdp.Run(ctx, tasks)
	if ctx.Err() == context.DeadlineExceeded {
		cancel()
		ScrapeSportsInteraction(fights, fighters, opponents)
		return 
	}
	if err != nil { panic(err) }

	cancel()
}