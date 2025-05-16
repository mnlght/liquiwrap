package decomposition_tools

import (
	"golang.org/x/net/html"
	"io"
	"regexp"
	"strings"
)

type StackElement struct {
	IsParent bool
	Token    html.Token
}

type ClassifiedStage struct {
	External     bool
	ExternalLink string
	MatchTables  []MatchTable
}

type MatchTable struct {
	Type         string
	MainStage    string
	SubStage     string
	ExternalLink string
	TokenContent []html.Token
}

func ClassifyTypeOfTournamentStage(r io.Reader) (map[string]ClassifiedStage, error) {
	tokenizer := html.NewTokenizer(r)
	var err error
	var matchTables []MatchTable
	var elementsStack []StackElement
	var grabMode bool
	var currentTableIndex int
	var mainStageSeek bool
	var currentMainStage string
	var subStageSeek bool
	var currentSubStage string
	var externalLink string
	var waitText bool
	var waitHref bool

	for err != io.EOF {
		tokenizer.Next()
		token := tokenizer.Token()

		if token.Type == 0 {
			err = tokenizer.Err()
		}

		switch token.Type {
		case html.StartTagToken:
			element := StackElement{
				Token: token,
			}

			if token.Data == "i" {
				waitText = true
				continue
			}

			if waitHref {
				if token.Data == "a" {
					if href := GetElHref(token.Attr); href != "" {
						externalLink = href
						matchTables[currentTableIndex-1].ExternalLink = externalLink
						externalLink = ""
					}

					waitHref = false
				}
			}

			if token.Data == "h2" {
				mainStageSeek = true
				currentSubStage = ""

			}
			if token.Data == "h3" {
				subStageSeek = true
			}

			if grabMode == true {
				matchTables[currentTableIndex].TokenContent = append(matchTables[currentTableIndex].TokenContent, token)
			}

			if token.Data == "div" {
				if MatchElClassByRegExp("brkts-matchlist brkts-matchlist-collapsible", token.Attr) {
					matchTables = append(matchTables, MatchTable{
						Type:         "bm",
						MainStage:    currentMainStage,
						SubStage:     currentSubStage,
						TokenContent: []html.Token{},
					})
					grabMode = true
					element.IsParent = true
				}

				if GetElClass(token.Attr) == "brkts-bracket-wrapper" {
					matchTables = append(matchTables, MatchTable{
						Type:         "br",
						MainStage:    currentMainStage,
						SubStage:     currentSubStage,
						TokenContent: []html.Token{},
					})
					grabMode = true
					element.IsParent = true
				}
			}

			if token.Data == "table" {
				if MatchElClassByRegExp("wikitable wikitable-striped sortable match-card", token.Attr) {

					matchTables = append(matchTables, MatchTable{
						Type:         "mc",
						MainStage:    currentMainStage,
						SubStage:     currentSubStage,
						TokenContent: []html.Token{},
					})
					grabMode = true
					element.IsParent = true
				}

				if GetElClass(token.Attr) == "wikitable wikitable-bordered wikitable-striped swisstable" {

					matchTables = append(matchTables, MatchTable{
						Type:         "sw",
						MainStage:    currentMainStage,
						SubStage:     currentSubStage,
						TokenContent: []html.Token{},
					})
					grabMode = true
					element.IsParent = true
				}

				if GetElClass(token.Attr) == "table table-bordered table-condensed crosstable" {
					matchTables = append(matchTables, MatchTable{
						Type:         "cr",
						MainStage:    currentMainStage,
						SubStage:     currentSubStage,
						TokenContent: []html.Token{},
					})
					grabMode = true
					element.IsParent = true
				}
			}

			elementsStack = append(elementsStack, element)

		case html.EndTagToken:
			for i := len(elementsStack) - 1; i >= 0; i-- {
				// we are try to close parent element
				if strings.EqualFold(elementsStack[i].Token.Data, token.Data) {
					if elementsStack[i].IsParent == true {
						grabMode = false
						elementsStack = elementsStack[:i]
						currentTableIndex++
					} else {
						elementsStack = elementsStack[:i]
					}
					if grabMode == true {
						matchTables[currentTableIndex].TokenContent = append(matchTables[currentTableIndex].TokenContent, token)
					}
					break
				}
			}

		case html.TextToken:
			if waitText {
				rz, _ := regexp.MatchString("(?i)\\bdetailed\\b.*\\bresults\\b", token.Data)

				if rz {
					waitHref = true
					continue
				}
				waitText = false
			}
			if mainStageSeek {
				currentMainStage = token.Data
				mainStageSeek = false
			}
			if subStageSeek {
				currentSubStage = token.Data
				subStageSeek = false
			}
			if grabMode == true {
				matchTables[currentTableIndex].TokenContent = append(matchTables[currentTableIndex].TokenContent, token)
			}
		}
	}

	mtMap := map[string]ClassifiedStage{}
	for _, v := range matchTables {
		if _, ok := mtMap[v.MainStage]; ok {
			l := mtMap[v.MainStage].MatchTables
			l = append(l, v)

			mtMap[v.MainStage] = ClassifiedStage{
				External:     v.ExternalLink != "",
				ExternalLink: v.ExternalLink,
				MatchTables:  l,
			}

			continue
		}
		l := make([]MatchTable, 0, 1)
		l = append(l, v)
		mtMap[v.MainStage] = ClassifiedStage{
			External:     v.ExternalLink != "",
			ExternalLink: v.ExternalLink,
			MatchTables:  l,
		}
	}

	return mtMap, nil
}
