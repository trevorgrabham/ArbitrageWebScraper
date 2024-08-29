package util

import "fmt"

type Bet struct {
	Fighters 	[2]string
	Sites    	[2]string
	Amounts  	[2]float64
	Profit 		float64
}

func (b Bet) String() string {
	return fmt.Sprintf("$%.2f PROFIT!\t$%.2f on %s @ %s, $%.2f on %s @ %s\n", b.Profit, b.Amounts[0], b.Fighters[0], b.Sites[0], b.Amounts[1], b.Fighters[1], b.Sites[1])
}

func FindArbitrageOpportunities(fights *ThreadSafeFights, fighters *ThreadSafeFighters, maxBet float64) []Bet {
	res := make([]Bet, 0)
	for _, fight := range fights.Fights {
		impliedOdds := [...]float64{1.0/fight.FighterA.BestSite.Odds, 1.0/fight.FighterB.BestSite.Odds}
		bookOdds := impliedOdds[0] + impliedOdds[1]
		if DEBUG {
			fmt.Printf("For %s book odds are %.2f\n", fight.String(), bookOdds)
		}
		if bookOdds < 1.0 {
			bet := Bet{
				Fighters: [...]string{fight.FighterA.Name.String(), fight.FighterB.Name.String()}, 
				Sites: [...]string{fight.FighterA.BestSite.Site, fight.FighterB.BestSite.Site}, 
				Amounts: [...]float64{(maxBet*impliedOdds[0])/bookOdds, (maxBet*impliedOdds[1])/bookOdds},
				Profit: fight.FighterA.BestSite.Odds * (maxBet*impliedOdds[0])/bookOdds - maxBet,
			}
			res = append(res, bet)
		}
	}
	return res
}