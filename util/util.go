package util

import "fmt"

const DEBUG = false

type SiteData struct {
	Site string
	Odds float64
}

func (s SiteData) String() string {
	return fmt.Sprintf("%f @ %s", s.Odds, s.Site)
}
