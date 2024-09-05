package util

import (
	"fmt"
	"strings"
	"sync"
)


type Fighter struct {
	Name     
	Sites    []SiteData
	BestSite SiteData
}

func (f *Fighter) String() string {
	var formatted strings.Builder
	formatted.WriteString(fmt.Sprintf("{%s: %.2f @ %s", f.Name, f.BestSite.Odds, f.BestSite.Site))
	for _, site := range f.Sites {
		if site.Site == f.BestSite.Site {
			continue
		}
		formatted.WriteString(fmt.Sprintf(", %.2f @ %s", site.Odds, site.Site))
	}
	formatted.WriteByte('}')
	return formatted.String()
}

type ThreadSafeFighters struct {
	Fighters map[Name]*Fighter
	Lock     sync.Mutex
}

func (f *ThreadSafeFighters) String() string {
	var valueString strings.Builder
	valueString.WriteByte('{')
	valueString.WriteByte('\n')
	for _, value := range f.Fighters {
		valueString.WriteByte('\t')
		valueString.WriteString(value.String())
		valueString.WriteByte('\n')
	}
	valueString.WriteByte('}')
	return valueString.String()
}

func (f *ThreadSafeFighters) exists(fighter, opponent *Fighter) (fighterA, fighterB *Fighter) {
	if DEBUG {
		fmt.Printf("Searching if %s or %s exist\n", fighter.Name, opponent.Name)
	}
	var defaultName = Name{}
	f.Lock.Lock()
	defer f.Lock.Unlock()
	for _, otherFighter := range f.Fighters {
		if n := fighter.Name.SameAs(otherFighter.Name); n != defaultName {
			if DEBUG {
				fmt.Printf("Matched %s to %s\n", fighter.Name, otherFighter.Name)
			}
			if n == otherFighter.Name {
				if DEBUG {
					fmt.Println("Keeping original name")
				}
				fighterA = otherFighter
				fighter.Name = otherFighter.Name
				continue
			}
			if DEBUG {
				fmt.Printf("Changing name to %s\n", n)
			}
			newFighter := &Fighter{
				Name: 			n,
				Sites: 			otherFighter.Sites,
				BestSite: 	otherFighter.BestSite,
			}
			delete(f.Fighters, otherFighter.Name)
			f.Fighters[n] = newFighter
			fighterA = newFighter
			continue
		} 
		if n := opponent.Name.SameAs(otherFighter.Name); n != defaultName {
			if DEBUG {
				fmt.Printf("Matched %s to %s\n", opponent, otherFighter.Name)
			}
			if n == otherFighter.Name {
				if DEBUG {
					fmt.Println("Keeping original name")
				}
				fighterB = otherFighter
				opponent.Name = n
				continue
			}
			if DEBUG {
				fmt.Printf("Changing name to %s\n", opponent.Name)
			}
			newFighter := &Fighter{
				Name: 			n,
				Sites: 			otherFighter.Sites,
				BestSite: 	otherFighter.BestSite,
			}
			delete(f.Fighters, otherFighter.Name)
			f.Fighters[n] = newFighter
			fighterB = newFighter
		}
	}
	if DEBUG {
		if fighterA == nil {
			fmt.Printf("No match for %s\n", fighter.Name)
		}
		if fighterB == nil {
			fmt.Printf("No match for %s\n", opponent.Name)
		}
		if fighterA != nil && fighterB != nil {
			fmt.Printf("Matched (%s, %s) to (%s, %s)\n", fighter.Name, opponent.Name, fighterA.Name, fighterB.Name)
		}
	}
	return
}

func (f *ThreadSafeFighters) AddFighters(fighter, opponent *Fighter) (fighterA, fighterB *Fighter) {
	if fighterA, fighterB = f.exists(fighter, opponent); fighterA != nil && fighterB != nil {
		f.AddOdds(fighterA.Name, fighter.BestSite)
		f.AddOdds(fighterB.Name, opponent.BestSite)
		return fighterA, fighterB
	}
	f.Lock.Lock()
	f.Fighters[fighter.Name] = fighter
	f.Fighters[opponent.Name] = opponent
	f.Lock.Unlock()
	return fighter, opponent
}

func (f *ThreadSafeFighters) AddOdds(name Name, site SiteData) (existed bool) {
	f.Lock.Lock()
	fighter, ok := f.Fighters[name]
	if site.Odds > fighter.BestSite.Odds {
		fighter.BestSite = site
	}
	fighter.Sites = append(fighter.Sites, site)
	f.Lock.Unlock()
	return ok
}

func (f *ThreadSafeFighters) BestSite(name Name) SiteData {
	f.Lock.Lock()
	best := f.Fighters[name].BestSite
	f.Lock.Unlock()
	return best
}

func (f *ThreadSafeFighters) FighterHasOdds(name Name, site SiteData) bool {
	f.Lock.Lock()
	defer f.Lock.Unlock()
	for _, s := range f.Fighters[name].Sites {
		if s.Site == site.Site { return true }
	}
	return false
}

func (f *ThreadSafeFighters) GetFighter(name Name) *Fighter {
	defaultName := Name{}
	if name == defaultName { return nil }
	f.Lock.Lock()
	fighter := f.Fighters[name]
	f.Lock.Unlock()
	return fighter
}