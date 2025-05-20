package content

import "time"

type CompletedMatch struct {
	MatchPage      string     `json:"matchPage"`
	Date           time.Time  `json:"date"`
	Game           string     `json:"game"`
	Round          string     `json:"round"`
	TeamLeft       string     `json:"teamLeft"`
	TeamRight      string     `json:"teamRight"`
	Score          string     `json:"score"`
	SeriesVods     []string   `json:"seriesVods"`
	SeriesScore    []string   `json:"matchScore"`
	SeriesDuration []string   `json:"matchDuration"`
	HeroesLeft     [][]string `json:"heroesLeft"`
	HeroesRight    [][]string `json:"heroesRight"`
	BansLeft       [][]string `json:"bansLeft"`
	BansRight      [][]string `json:"bansRight"`
}
