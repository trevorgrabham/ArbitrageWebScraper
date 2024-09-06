package util

import "fmt"

const DEBUG = false

var Headers = map[string]interface{}{
	"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36"}

type SiteData struct {
	Site string
	Odds float64
}

func (s SiteData) String() string {
	return fmt.Sprintf("%f @ %s", s.Odds, s.Site)
}
