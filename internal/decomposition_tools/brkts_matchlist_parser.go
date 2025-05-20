package decomposition_tools

import (
	"fmt"
	"github.com/mnlght/liquiwrap/content"
	"golang.org/x/net/html"
	"strings"
)

func BrktsMatchlistMatchesPickOut(mccontent MatchTable, game string) ([]content.CompletedMatch, error) {
	var matches []content.CompletedMatch
	var el content.CompletedMatch
	var inMatchSeekMode bool
	var opponentIndex int
	var opponentSeek bool
	var scoreIndex int
	var scoreSeek bool
	var dateSeek bool
	var spaceSeek bool
	var spaceIndex int
	var vodsSeek bool

	var picksGameSeek bool
	var picksGameIndex int
	var picksThumbsSeek bool
	var picksThumbsIndex int
	var bansGameSeek bool
	var bansGameIndex int
	var bansThumbsSeek bool
	var bansThumbsIndex int
	var tr int

	combineCh := make(chan content.CompletedMatch)
	combinedCh := make(chan []content.CompletedMatch)

	go func() {
		for v := range combineCh {
			v.Game = game
			matches = append(matches, v)
		}
		combinedCh <- matches
		close(combinedCh)
	}()

	go func() {
		for i := 0; i < len(mccontent.TokenContent); i++ {
			if dateSeek {
				if mccontent.TokenContent[i].Type == html.TextToken {
					if mccontent.TokenContent[i].Data != "" {
						el.Date = GetMetaMatchDateWithTime(mccontent.TokenContent[i].Data).DateStart
						dateSeek = false
					}
				}
			}
			if spaceSeek {
				if mccontent.TokenContent[i].Type == html.TextToken {
					if mccontent.TokenContent[i].Data != "" {
						f := strings.Split(mccontent.TokenContent[i].Data, ":")
						if len(f) == 1 {
							el.SeriesDuration = append(el.SeriesDuration, mccontent.TokenContent[i].Data)
						}
						if len(f) == 2 {
							el.SeriesScore = append(el.SeriesScore, mccontent.TokenContent[i].Data)
						}
						spaceSeek = false
					}
				}
			}
			if scoreSeek {
				if mccontent.TokenContent[i].Type == html.TextToken {
					if mccontent.TokenContent[i].Data != "" {
						if scoreIndex == 1 {
							el.Score = mccontent.TokenContent[i].Data
						}
						if scoreIndex == 2 {
							el.Score = fmt.Sprintf("%s:%s", el.Score, mccontent.TokenContent[i].Data)
						}

						scoreSeek = false
					}
				}
			}
			if mccontent.TokenContent[i].Type == html.StartTagToken {
				if vodsSeek {
					if mccontent.TokenContent[i].Data == "a" {
						el.SeriesVods = append(el.SeriesVods, GetElAttribute("href", mccontent.TokenContent[i].Attr))
					}
				}
				if bansThumbsSeek {
					if mccontent.TokenContent[i].Data == "a" {
						if bansThumbsIndex == 1 {
							el.BansLeft[bansGameIndex-1] = append(el.BansLeft[bansGameIndex-1], GetElAttribute("title", mccontent.TokenContent[i].Attr))
						}
						if bansThumbsIndex == 2 {
							el.BansRight[bansGameIndex-1] = append(el.BansRight[bansGameIndex-1], GetElAttribute("title", mccontent.TokenContent[i].Attr))
						}
					}
				}
				if picksThumbsSeek {
					if mccontent.TokenContent[i].Data == "a" {
						if picksThumbsIndex == 1 {
							el.HeroesLeft[picksGameIndex-1] = append(el.HeroesLeft[picksGameIndex-1], GetElAttribute("title", mccontent.TokenContent[i].Attr))
						}
						if picksThumbsIndex == 2 {
							el.HeroesRight[picksGameIndex-1] = append(el.HeroesRight[picksGameIndex-1], GetElAttribute("title", mccontent.TokenContent[i].Attr))
						}
					}
				}
				if opponentSeek {
					if mccontent.TokenContent[i].Data == "span" {
						a := GetElAttribute("data-highlightingclass", mccontent.TokenContent[i].Attr)
						if a != "" {
							if opponentIndex == 1 {
								el.TeamLeft = a
							}
							if opponentIndex == 2 {
								el.TeamRight = a
							}
							opponentSeek = false
						}
					}
				}
				if bansGameIndex != 0 {
					if mccontent.TokenContent[i].Data == "td" && GetElAttribute("style", mccontent.TokenContent[i].Attr) == "float:left" {
						bansThumbsSeek = true
						bansThumbsIndex = 1
					}
					if mccontent.TokenContent[i].Data == "td" && GetElAttribute("style", mccontent.TokenContent[i].Attr) == "float:right" {
						bansThumbsSeek = true
						bansThumbsIndex = 2
					}
				}
				if bansGameSeek {
					if mccontent.TokenContent[i].Data == "tr" {
						if tr > 0 {
							el.BansLeft = append(el.BansLeft, []string{})
							el.BansRight = append(el.BansRight, []string{})
							bansGameIndex++
						}

						tr++
					}
				}

				if mccontent.TokenContent[i].Data == "span" {
					if MatchElClassByRegExp("timer-object", mccontent.TokenContent[i].Attr) {
						dateSeek = true
					}
				}
				if mccontent.TokenContent[i].Data == "div" {

					if vodsSeek {
						vodsSeek = false
					}

					if GetElClass(mccontent.TokenContent[i].Attr) == "brkts-popup-spaced vodlink" {
						vodsSeek = true
					}
					if GetElClass(mccontent.TokenContent[i].Attr) == "brkts-popup-spaced" {
						spaceSeek = true
						spaceIndex++
					}
					if MatchElClassByRegExp("brkts-popup-body-element brkts-popup-body-game", mccontent.TokenContent[i].Attr) {
						picksGameSeek = true
						picksThumbsIndex = 0
						el.HeroesLeft = append(el.HeroesLeft, []string{})
						el.HeroesRight = append(el.HeroesRight, []string{})
						picksGameIndex++
					}
					if picksGameSeek {
						if MatchElClassByRegExp("brkts-popup-body-element-thumbs", mccontent.TokenContent[i].Attr) {
							picksThumbsSeek = true
							picksThumbsIndex++
						}
					}
					if MatchElClassByRegExp("brkts-matchlist-score", mccontent.TokenContent[i].Attr) {
						scoreSeek = true
						scoreIndex++
					}
					if MatchElClassByRegExp("brkts-popup-footer", mccontent.TokenContent[i].Attr) || MatchElClassByRegExp("brkts-popup-comment", mccontent.TokenContent[i].Attr) {
						bansGameSeek = false
						opponentSeek = false
						scoreSeek = false
						picksGameSeek = false
						bansThumbsSeek = false
						picksThumbsSeek = false
						dateSeek = false
						spaceSeek = false
						bansGameIndex = 0
						spaceIndex = 0
						picksGameIndex = 0
						scoreIndex = 0
						tr = 0
						opponentIndex = 0
						bansThumbsIndex = 0
					}

					if MatchElClassByRegExp("brkts-popup-mapveto", mccontent.TokenContent[i].Attr) {
						bansGameSeek = true
						bansThumbsIndex = 0
					}
					if MatchElClassByRegExp("brkts-matchlist-opponent", mccontent.TokenContent[i].Attr) {
						opponentSeek = true
						opponentIndex++
					}

					if MatchElClassByRegExp("brkts-matchlist-match", mccontent.TokenContent[i].Attr) {
						bansGameSeek = false
						opponentSeek = false
						scoreSeek = false
						picksGameSeek = false
						bansThumbsSeek = false
						picksThumbsSeek = false
						spaceSeek = false
						bansGameIndex = 0
						spaceIndex = 0
						picksGameIndex = 0
						scoreIndex = 0
						tr = 0
						opponentIndex = 0
						bansThumbsIndex = 0
						if inMatchSeekMode == true {
							combineCh <- el
							//matches = append(matches, el)
						}

						el = content.CompletedMatch{}
						if inMatchSeekMode == false {
							inMatchSeekMode = true
						}
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
