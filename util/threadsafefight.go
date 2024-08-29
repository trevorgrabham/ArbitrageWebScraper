package util

import (
	"fmt"
	"strings"
	"sync"
)

type Fight struct {
	FighterA *Fighter
	FighterB *Fighter
}

func (f *Fight) String() string {
	return fmt.Sprintf("%s vs %s", f.FighterA.Name, f.FighterB.Name)
}

type ThreadSafeFights struct {
	Fights map[Name]*Fight
	Lock   sync.Mutex
}

func (f *ThreadSafeFights) String() string {
	var valueString strings.Builder
	valueString.WriteString("{\n")
	for _, f := range f.Fights {
		valueString.WriteString(f.String())
		valueString.WriteByte('\n')
	}
	valueString.WriteByte('}')
	return valueString.String()
}

func (f *ThreadSafeFights) exists(name Name) bool {
	f.Lock.Lock()
	_, ok := f.Fights[name]
	f.Lock.Unlock()
	return ok
}

func (f *ThreadSafeFights) AddFight(fight *Fight) (alreadyExists bool) {
	if f.exists(fight.FighterA.Name) && f.exists(fight.FighterB.Name) {
		return true
	}
	f.Lock.Lock()
	f.Fights[fight.FighterA.Name] = fight
	f.Fights[fight.FighterB.Name] = fight
	f.Lock.Unlock()
	return false
}

func (f *ThreadSafeFights) Opponent(fighter *Fighter) *Fighter {
	f.Lock.Lock()
	defer f.Lock.Unlock()
	fight := f.Fights[fighter.Name]
	if fight.FighterA.Name == fighter.Name {
		return fight.FighterB
	} else {
		return fight.FighterA
	}
}
