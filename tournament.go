package liquiwrap

import "time"

type Tournament struct {
	Tier         string
	Name         string
	Prize        string
	Winner       string
	Year         int
	DateStart    time.Time
	DateEnd      time.Time
	Game         string
	PageForParse string
}
