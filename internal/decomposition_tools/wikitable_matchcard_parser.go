package decomposition_tools

import (
	"encoding/json"
	"fmt"
	"github.com/mnlght/liquiwrap"
	"golang.org/x/net/html"
)

func WikitableMatchcardMatchesPickOut(mccontent MatchTable) ([]liquiwrap.CompletedMatch, error) {
	var matches []liquiwrap.CompletedMatch
	res := make(map[string]string)
	var seekKey string
	var seekTRMode bool
	var seekTDMode bool
	var trCounter int

	fmt.Println("mc type found")
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

func MarshalWikitableMapToStruct(matches []liquiwrap.CompletedMatch, info map[string]string) ([]liquiwrap.CompletedMatch, error) {
	var match liquiwrap.CompletedMatch
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
