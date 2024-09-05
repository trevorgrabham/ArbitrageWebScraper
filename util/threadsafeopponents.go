package util

import (
	"fmt"
	"strings"
	"sync"
)

type ThreadSafeOpponents struct {
	Lock sync.Mutex
	Opponents 	map[Name]Name
}

func (o *ThreadSafeOpponents) AddPairing(fighterA, fighterB Name) {
	o.Lock.Lock()
	defer o.Lock.Unlock()
	o.Opponents[fighterA] = fighterB
	o.Opponents[fighterB] = fighterA
}

func (o *ThreadSafeOpponents) GetOpponent(fighter Name) Name {
	o.Lock.Lock()
	defer o.Lock.Unlock()
	return o.Opponents[fighter]
}

func (o *ThreadSafeOpponents) String() string {
	var res strings.Builder
	res.WriteString("{\n")
	o.Lock.Lock()
	defer o.Lock.Unlock()
	for key, value := range o.Opponents {
		res.WriteString(fmt.Sprintf("\t%s vs %s\n", key, value))
	}
	res.WriteByte('}')
	return res.String()
}