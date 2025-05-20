package decomposition_tools

import (
	"fmt"
	"github.com/mnlght/liquiwrap/content"
	"golang.org/x/net/html"
	"io"
	"strconv"
	"strings"
)

func OngoingMatchesCSPickOut(r io.Reader) ([]content.OngoingMatch, error) {
	tokenizer := html.NewTokenizer(r)
	var err error
	var matches []content.OngoingMatch
	var el content.OngoingMatch
	var inMatchSeekMode bool
	var teamleftSeekMode bool
	var teamrightSeekMode bool
	var scoreSeekMode bool
	var formatSeekMode bool
	var timeSeekMode bool
	var tournamentInfoSeekMode bool
	var scoreStr []string
	var divScoreCount int

	combineCh := make(chan content.OngoingMatch)
	combinedCh := make(chan []content.OngoingMatch)

	go func() {
		am := make(map[string]bool)
		for v := range combineCh {
			d := fmt.Sprintf("%s-%s-%s-%s", v.Game, v.TeamLeft, v.TeamRight, v.DateStart)
			v.Game = "counterstrike"
			if _, ok := am[d]; !ok {
				matches = append(matches, v)
			}
			am[d] = true
		}
		combinedCh <- matches
		close(combinedCh)
	}()

	go func() {
		for err != io.EOF {
			tokenizer.Next()
			token := tokenizer.Token()

			if token.Type == 0 {
				err = tokenizer.Err()
			}
			if timeSeekMode {
				if token.Type == html.TextToken {
					el.DateStart = GetMetaMatchDateWithTime(token.Data).DateStart
					timeSeekMode = false
				}
			}
			if scoreSeekMode {
				if token.Data == "div" {
					if divScoreCount == 1 && scoreSeekMode {
						if len(scoreStr) == 1 {
							score := strings.Join(scoreStr, "")
							el.Score = score
							el.MapNumber = 1
							scoreDivided := strings.Split(score, ":")
							if len(scoreDivided) == 2 {
								n1, err := strconv.Atoi(scoreDivided[0])
								n2, err := strconv.Atoi(scoreDivided[1])
								if err == nil {
									el.MapNumber = n1 + n2 + 1
								}
							}
						}
						scoreSeekMode = false
						formatSeekMode = true

						divScoreCount = 0
					}
				}
				if token.Type == html.TextToken {
					if token.Data != "vs" {
						scoreStr = append(scoreStr, token.Data)
					}
					divScoreCount++
				}
			}

			if token.Type == html.StartTagToken {
				if token.Data == "a" {
					if tournamentInfoSeekMode {
						el.TournamentName = GetElAttribute("title", token.Attr)
						el.TournamentLink = GetElAttribute("href", token.Attr)

						tournamentInfoSeekMode = false
					}
				}
				if token.Data == "abbr" {
					if formatSeekMode {
						el.Format = GetElAttribute("title", token.Attr)
						formatSeekMode = false
					}
				}
				if token.Data == "span" {
					if MatchElClassByRegExp("league-icon", token.Attr) {
						tournamentInfoSeekMode = true
					}
					if MatchElClassByRegExp("timer-object", token.Attr) {
						timeSeekMode = true
					}
					if teamleftSeekMode == true {
						el.TeamLeft = GetElAttribute("data-highlightingclass", token.Attr)
						teamleftSeekMode = false
					}
					if teamrightSeekMode == true {
						el.TeamRight = GetElAttribute("data-highlightingclass", token.Attr)
						teamrightSeekMode = false
					}
				}
				if token.Data == "table" {
					if inMatchSeekMode == true {
						combineCh <- el
						//matches = append(matches, el)
					}

					el = content.OngoingMatch{
						Game: "counterstrike",
					}
					if inMatchSeekMode == false {
						inMatchSeekMode = true
					}
				}
				if token.Data == "td" {
					//if GetElAttribute("class", token.Attr) == "match" {
					//	if inMatchSeekMode == true {
					//		matches = append(matches, el)
					//	}
					//
					//	el = entity.OngoingMatch{
					//		Game: game,
					//	}
					//	if inMatchSeekMode == false {
					//		inMatchSeekMode = true
					//	}
					//}
					if MatchElClassByRegExp("team-left", token.Attr) {
						teamleftSeekMode = true
					}
					if GetElAttribute("class", token.Attr) == "versus" {
						scoreStr = []string{}
						scoreSeekMode = true
					}
					if GetElAttribute("class", token.Attr) == "versus-lowe" {

					}
					if MatchElClassByRegExp("team-right", token.Attr) {

						teamrightSeekMode = true
					}

					//if GetElAttribute("class", token.Attr) == "match-tournament" {
					//	tournamentInfoSeekMode = true
					//}
					if GetElAttribute("data-filter-expansion-template", token.Attr) == "MainPageMatches/Upcoming" {
						//upcomingSeekMode = true
						//continue
					}
					if GetElAttribute("data-filter-expansion-template", token.Attr) == "MainPageMatches/Completed" {
						//upcomingSeekMode = true
						//continue
					}

				}
			}
		}

		combineCh <- el
		close(combineCh)
		//matches = append(matches, el)
	}()

	return <-combinedCh, nil
}
