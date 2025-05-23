package decomposition_tools

import (
	"encoding/json"
	"github.com/mnlght/liquiwrap/content"
	"golang.org/x/net/html"
)

func WikitableMatchcardMatchesPickOut(mccontent MatchTable) ([]content.CompletedMatch, error) {
	var matches []content.CompletedMatch
	res := make(map[string]string)
	var seekKey string
	var seekTRMode bool
	var seekTDMode bool
	var trCounter int

	for i := 0; i < len(mccontent.TokenContent); i++ {
		if i == len(mccontent.TokenContent)-1 {

			matchesWithFinalTR, err := MarshalWikitableMapToStruct(matches, res)
			if err != nil {
				return nil, err
			}
			matches = matchesWithFinalTR
		}

		if mccontent.TokenContent[i].Type == html.TextToken {
			if seekTDMode {
				if mccontent.TokenContent[i].Data != "" && mccontent.TokenContent[i].Data != " " {
					res[seekKey] = mccontent.TokenContent[i].Data
					seekTDMode = false
				}
			}
		}
		if mccontent.TokenContent[i].Type == html.StartTagToken {
			if mccontent.TokenContent[i].Data == "tr" && GetElClass(mccontent.TokenContent[i].Attr) != "HeaderRow" {
				if seekTRMode {
					matchesWithLineTR, err := MarshalWikitableMapToStruct(matches, res)
					if err != nil {
						return nil, err
					}

					matches = matchesWithLineTR
				}
				if seekTRMode == false {
					seekTRMode = true
				}

				trCounter++
			}

			if mccontent.TokenContent[i].Data == "td" {
				if seekTRMode {
					seekTDMode = true
					seekKey = GetElClass(mccontent.TokenContent[i].Attr)
				}
			}
		}
	}

	return matches, nil
}

func MarshalWikitableMapToStruct(matches []content.CompletedMatch, info map[string]string) ([]content.CompletedMatch, error) {
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
