package decomposition_tools

import (
	"fmt"
	"github.com/mnlght/liquiwrap/content"
	"golang.org/x/net/html"
)

func TableCrosstableMatchesPickOut(mccontent MatchTable) ([]content.CompletedMatch, error) {
	var matches []content.CompletedMatch
	var line []string
	for i := 0; i < len(mccontent.TokenContent); i++ {
		if mccontent.TokenContent[i].Type == html.StartTagToken {
			if mccontent.TokenContent[i].Data == "span" {
				f := GetElAttribute("data-highlightingclass", mccontent.TokenContent[i].Attr)
				if f != "" {
					p := true
					for _, r := range line {
						if f == r {
							p = false
						}
					}

					if p {
						line = append(line, f)
					}
				}

			}
		}
	}

	var seekInTrMode bool
	var seekInTHMode bool
	excMap := make(map[string]bool)
	var teamRight string
	var teamIndex int
	for i := 0; i < len(mccontent.TokenContent); i++ {
		if mccontent.TokenContent[i].Type == html.StartTagToken {
			if seekInTHMode {
				if mccontent.TokenContent[i].Data == "span" {
					a := GetElAttribute("data-highlightingclass", mccontent.TokenContent[i].Attr)
					if a != "" {
						teamRight = a

						seekInTHMode = false
					}
				}
			}

			if seekInTrMode {
				if mccontent.TokenContent[i].Data == "th" {
					seekInTHMode = true
				}
			}

			if mccontent.TokenContent[i].Data == "tr" && GetElClass(mccontent.TokenContent[i].Attr) == "crosstable-tr" {
				seekInTrMode = true
				teamIndex = 0
			}

			if mccontent.TokenContent[i].Data == "td" {
				//seekInTdMode = true
				if GetElClass(mccontent.TokenContent[i].Attr) == "crosstable-bgc-cross" {
					teamIndex++
					continue
				}

				if _, ok := excMap[fmt.Sprintf("%s-%s", teamRight, line[teamIndex])]; !ok {
					matches = append(matches, content.CompletedMatch{
						TeamLeft:  teamRight,
						TeamRight: line[teamIndex],
						Score:     BatchScore(GetElClass(mccontent.TokenContent[i].Attr)),
					})

					excMap[fmt.Sprintf("%s-%s", line[teamIndex], teamRight)] = true
				}
				teamIndex++
			}
		}
	}

	return matches, nil
}
