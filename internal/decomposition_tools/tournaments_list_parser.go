package decomposition_tools

import (
	"github.com/mnlght/liquiwrap/content"
	"golang.org/x/net/html"
	"io"
	"strings"
)

func CollectTournamentMetaList(r io.Reader, game string) []content.Tournament {
	tokenizer := html.NewTokenizer(r)
	var err error
	var remo bool
	var el content.Tournament
	var ts []content.Tournament
	var headerLinkSeek bool
	var inLinkNow bool
	var dateSeek bool
	var prizeSeek bool
	var winnerSeek bool
	var gridI int

	for err != io.EOF {
		tokenizer.Next()
		token := tokenizer.Token()

		if token.Type == 0 {
			err = tokenizer.Err()
		}

		if token.Type == html.TextToken {
			if headerLinkSeek {
				if inLinkNow {
					el.Name = token.Data
					inLinkNow = false
				}

			}
			if winnerSeek {
				if token.Data != "" && len(token.Data) > 2 {
					el.Winner = token.Data
					winnerSeek = false
				}
			}
			if dateSeek {
				m := GetMetaTournamentDate(token.Data)
				if m != nil {
					el.DateStart = m.DateStart
					el.DateEnd = m.DateEnd
					el.Year = m.Year
				}
				dateSeek = false
			}

			if prizeSeek {
				if token.Data != "" {
					el.Prize = token.Data
					prizeSeek = false
				}

			}
		}

		if token.Type == html.StartTagToken {
			if headerLinkSeek {
				if token.Data == "a" {
					href := GetElAttribute("href", token.Attr)
					if href != "" {
						hs := strings.Split(href, game)
						if len(hs) > 1 {
							el.PageForParse = strings.TrimPrefix(hs[1], "/")
							inLinkNow = true
						}
					}

				}
			}
			if token.Data == "div" {
				if GetElAttribute("class", token.Attr) == "gridCell Tournament Header" {
					headerLinkSeek = true
					dateSeek = false
					prizeSeek = false
					winnerSeek = false
				}

				if GetElAttribute("class", token.Attr) == "gridCell EventDetails Date Header" {
					dateSeek = true
					headerLinkSeek = false
					prizeSeek = false
					winnerSeek = false
				}

				if MatchElClassByRegExp("gridCell EventDetails Prize Header", token.Attr) {
					prizeSeek = true
					dateSeek = false
					headerLinkSeek = false
					winnerSeek = false
				}

				if MatchElClassByRegExp("gridCell Placement FirstPlace", token.Attr) {
					winnerSeek = true
					prizeSeek = false
					dateSeek = false
					headerLinkSeek = false
				}

				if MatchElClassByRegExp("gridRow", token.Attr) {
					gridI++
					if remo == true {
						ts = append(ts, el)
						dateSeek = false
						headerLinkSeek = false
						prizeSeek = false
					}

					el = content.Tournament{}
					el.Game = game

					if remo == false {
						remo = true
					}
				}
			}
		}
	}

	if gridI != 0 {
		ts = append(ts, el)
	}

	return ts
}
