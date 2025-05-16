package decomposition_tools

import (
	"encoding/json"
	"github.com/mnlght/liquiwrap/content"
	"golang.org/x/net/html"
)

func WikitableSwisstableMatchesPickOut(mccontent MatchTable) ([]content.CompletedMatch, error) {
	var matches []content.CompletedMatch
	res := make(map[string]string)
	var seekTRMode bool
	var seekTDMode bool
	var trCounter int
	var tdCounter int
	var intoTrCounter int

	for i := 0; i < len(mccontent.TokenContent); i++ {
		if mccontent.TokenContent[i].Type == html.TextToken {
			if mccontent.TokenContent[i].Data != "" && mccontent.TokenContent[i].Data != " " {
				if seekTDMode {
					res["Score"] = mccontent.TokenContent[i].Data

					matchesWithFinalTR, err := MarshalSwissMapToStruct(matches, res)
					if err != nil {
						return nil, err
					}
					matches = matchesWithFinalTR

					intoTrCounter = 0
					seekTDMode = false
				}
			}
		}

		if mccontent.TokenContent[i].Type == html.StartTagToken {
			if seekTDMode {
				if mccontent.TokenContent[i].Data == "span" {
					d := GetElAttribute("data-highlightingclass", mccontent.TokenContent[i].Attr)
					if d != "" {
						res["TeamRight"] = d
					}
					intoTrCounter++
				}
			}

			if mccontent.TokenContent[i].Data == "tr" {
				if seekTRMode == false && trCounter > 0 {
					seekTRMode = true
				}
				tdCounter = 0
				trCounter++
			}

			if mccontent.TokenContent[i].Data == "td" {
				if seekTRMode {
					if tdCounter == 0 {
						res["TeamLeft"] = GetElAttribute("data-highlightingkey", mccontent.TokenContent[i].Attr)
					}

					if tdCounter == 1 {
						tdCounter++
						continue
					}

					if tdCounter > 1 {
						if MatchElClassByRegExp("swisstable-bgc", mccontent.TokenContent[i].Attr) {
							seekTDMode = true
						}

					}
					tdCounter++
				}
			}
		}
	}
	return matches, nil
}

func MarshalSwissMapToStruct(matches []content.CompletedMatch, info map[string]string) ([]content.CompletedMatch, error) {
	var match content.CompletedMatch
	b, err := json.Marshal(info)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &match)
	if err != nil {
		return nil, err
	}

	return append(matches, match), nil
}
