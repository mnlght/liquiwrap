package liquiwrap

import "time"

type OngoingMatch struct {
	ID             string
	TeamLeft       string
	TeamRight      string
	Score          string
	Format         string
	TournamentName string
	TournamentLink string
	Archive        bool
	MapNumber      int
	DateStart      time.Time
	Game           string
}
