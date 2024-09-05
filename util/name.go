package util

import (
	"fmt"
	"strings"
)

type Name struct {
	FirstName string
	LastName  string
}

func (n Name) SameAs(other Name) (formatted Name) {
	if DEBUG {
		if n.FirstName != "" {
			fmt.Printf("Name %v. Has first initial?\t%v\n", n, n.FirstName[1] == '.')
		}
	}
	switch {
	// Just a single name
	case n.FirstName == "":
		otherCombined := other.FirstName + other.LastName
		if n.LastName == otherCombined { return other }

		otherReverseCombined := other.LastName + other.FirstName
		if n.LastName == otherReverseCombined { return other }

		if other.FirstName != "" { return Name{} }

		if n.LastName == other.LastName { return n }
	case other.FirstName == "":
		nCombined := n.FirstName + n.LastName
		if other.LastName == nCombined { return n }

		nReverseCombined := n.LastName + n.FirstName
		if other.LastName == nReverseCombined { return n }

		if n.FirstName != "" { return Name{} }

		if n.LastName == other.LastName { return n }
	// Just a first initial
	case n.FirstName[1] == '.':
		if n.FirstName[0] == other.FirstName[0] && distanceBetweenNames(n.LastName, other.LastName) < 2 { return other }
	case other.FirstName[1] == '.':
		if n.FirstName[0] == other.FirstName[0] && distanceBetweenNames(n.LastName, other.LastName) < 2 { return n }
	// Just a last initial
	case n.LastName[1] == '.':
		if n.FirstName == other.FirstName && distanceBetweenNames(n.LastName, other.LastName) < 2 { return other }
	// Just a last initial
	case other.LastName[1] == '.':
		if n.FirstName == other.FirstName && distanceBetweenNames(n.LastName, other.LastName) < 2 { return n }
	// First and last names
	default:
		if n.FirstName == other.FirstName && distanceBetweenNames(n.LastName, other.LastName) < 3 { return n }

		if n.LastName == other.LastName && distanceBetweenNames(n.FirstName, other.FirstName) < 3 { return n }
	}
	return Name{}
}

func NewName(name string) Name {
	name = strings.ReplaceAll(name, "\u00A0", " ")
	if DEBUG {
		fmt.Printf("Making name for %s\n", name)
	}
	name = strings.ToLower(strings.TrimSpace(name))
	mySplitName := strings.Split(name, " ")
	var myName Name
	if len(mySplitName) > 1 {
		myName.FirstName = mySplitName[0]
		myName.LastName = strings.Join(mySplitName[1:], " ")
	} else {
		myName.LastName = name
	}
	if DEBUG {
		fmt.Printf("Created name %v for %s\n", myName, name)
	}
	return myName
}

// func (n Name) String() string {
// 	if n.FirstName == "" {
// 		return n.LastName
// 	}
// 	return fmt.Sprintf("%s %s", n.FirstName, n.LastName)
// }

func distanceBetweenNames(first, second string) (numDiffs int) {
	if first == second { return 0 }
	dist := make([][]int, len(first)+1)
	for i := range dist {
		dist[i] = make([]int, len(second)+1)
		dist[i][0] = i
	} 
	for i := range dist[0] {
		dist[0][i] = i
	}
	for row := 1;row < len(dist); row++ {
		for col := 1;col < len(dist[row]); col++ {
			if first[row-1] == second[col-1] {
				dist[row][col] = dist[row-1][col-1]
				continue
			}
			dist[row][col] = 1 + min(dist[row-1][col-1], dist[row-1][col], dist[row][col-1])
		}
	}
	return dist[len(dist)-1][len(dist[0])-1]
}