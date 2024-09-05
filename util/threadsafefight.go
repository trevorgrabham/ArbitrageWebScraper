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
	Fights []*Fight
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

func (f *ThreadSafeFights) AddFight(fight *Fight) (err error) {
	if fight == nil { return fmt.Errorf("addfight(): got nil value") }
	f.Lock.Lock()
	f.Fights = append(f.Fights, fight)
	f.Lock.Unlock()
	return nil
}