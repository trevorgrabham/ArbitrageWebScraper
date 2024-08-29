package util

import (
	"fmt"
	"strings"
	"sync"
)

type Name struct {
	FirstName 		string
	LastName 			string
}

func (n Name) SameAs(other Name) (formatted Name) {
	switch {
	// Just a single name
	case n.FirstName == "":
		if n.LastName == other.LastName {
			return other
		}
	// Just a single name
	case other.FirstName == "":
		if n.LastName == other.LastName {
			return n
		}
	// Just a first initial
	case len(n.FirstName) < 3:
		if n.LastName == other.LastName && n.FirstName[0] == other.FirstName[0] {
			return other
		}
	// Just a first initial
	case len(other.FirstName) < 3:
		if n.LastName == other.LastName && n.FirstName[0] == other.FirstName[0] {
			return n
		}
	// Just a last initial
	case len(n.LastName) < 3:
		if n.FirstName == other.FirstName && n.LastName[0] == other.LastName[0] {
			return other
		}
	// Just a last initial
	case len(other.LastName) < 3:
		if n.FirstName == other.FirstName && n.LastName[0] == other.LastName[0] {
			return n
		}
	}
	return Name{}
}

func NewName(name string) Name {
	mySplitName := strings.Split(name, " ")
	var myName Name
	if len(mySplitName) > 1 {
		myName.FirstName = mySplitName[0]
		myName.LastName = strings.Join(mySplitName[1:], " ")
	} else {
		myName.LastName = name
	}
	return myName
}

func (n Name) String() string {
	if n.FirstName == "" {return n.LastName}
	return fmt.Sprintf("%s %s", n.FirstName, n.LastName)
}

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

func (f *ThreadSafeFighters) exists(name, opponent Name) (exists, matchedOpponent bool, fighter *Fighter) {
	if DEBUG {
		fmt.Printf("Searching if %s or %s exist\n\n", name, opponent)
	}
	var defaultName = Name{}
	f.Lock.Lock()
	defer f.Lock.Unlock()
	for _, fighter := range f.Fighters {
		if n := fighter.SameAs(name); n != defaultName {
			if DEBUG {
				fmt.Printf("Matched %s to %s\n", name, fighter.Name)
			}
			if n == fighter.Name {
				if DEBUG {
					fmt.Println("Keeping original name")
				}
				return true, false, fighter
			}
			if DEBUG {
				fmt.Printf("Changing name to %s\n", name)
			}
			newFighter := &Fighter{
				Name: 			n,
				Sites: 			fighter.Sites,
				BestSite: 	fighter.BestSite,
			}
			delete(f.Fighters, fighter.Name)
			f.Fighters[n] = newFighter
			return true, false, fighter
		} 
		if n := fighter.SameAs(opponent); n != defaultName {
			if DEBUG {
				fmt.Printf("Matched %s to %s\n", opponent, fighter.Name)
			}
			if n == fighter.Name {
				if DEBUG {
					fmt.Println("Keeping original name")
				}
				return true, false, fighter
			}
			if DEBUG {
				fmt.Printf("Changing name to %s\n", opponent)
			}
			newFighter := &Fighter{
				Name: 			n,
				Sites: 			fighter.Sites,
				BestSite: 	fighter.BestSite,
			}
			delete(f.Fighters, fighter.Name)
			f.Fighters[n] = newFighter
			return true, false, fighter
		}
	}
	if DEBUG {
		fmt.Printf("No match found for %s or %s\n", name, opponent)
	}
	return false, false, nil
}

func (f *ThreadSafeFighters) AddFighters(fighter, opponent *Fighter) (alreadyExisted, matchedOpponent bool, fighterRef *Fighter) {
	if exists, isOp, fighterRef := f.exists(fighter.Name, opponent.Name); exists {
		return exists, isOp, fighterRef
	}
	f.Lock.Lock()
	f.Fighters[fighter.Name] = fighter
	f.Fighters[opponent.Name] = opponent
	f.Lock.Unlock()
	return false, false, nil
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