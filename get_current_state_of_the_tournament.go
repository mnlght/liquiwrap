package liquiwrap

import (
	"fmt"
	"github.com/mnlght/liquiwrap/content"
	"github.com/mnlght/liquiwrap/internal"
	"github.com/mnlght/liquiwrap/internal/decomposition_tools"
	"strings"
)

type GetCurrentStateOfTheTournament struct {
	Game string
	Url  string
}

func NewGetCurrentStateOfTheTournament(game string, url string) *GetCurrentStateOfTheTournament {
	return &GetCurrentStateOfTheTournament{
		Game: game,
		Url:  url,
	}
}

func (g *GetCurrentStateOfTheTournament) Action() ([]content.CompletedMatch, error) {
	p := internal.NewLiquipediaPageClient(fmt.Sprintf("/%s/%s", g.Game, g.Url))
	pageContent, err := p.Do()
	if err != nil {
		return nil, err
	}

	mts, err := decomposition_tools.ClassifyTypeOfTournamentStage(strings.NewReader(string(pageContent)))
	if err != nil {
		return nil, err
	}

	var result []content.CompletedMatch
	for _, v := range mts {
		if v.External {
			pz := internal.NewLiquipediaPageClient(v.ExternalLink)
			pageContentExternal, err := pz.Do()
			if err != nil {
				return nil, err
			}

			stagesExternal, err := decomposition_tools.ClassifyTypeOfTournamentStage(strings.NewReader(string(pageContentExternal)))
			var allExternalMatchBlocks []decomposition_tools.MatchTable
			for _, se := range stagesExternal {
				allExternalMatchBlocks = append(allExternalMatchBlocks, se.MatchTables...)
			}
			matches, err := g.seekMatchesInPage(allExternalMatchBlocks)
			if err != nil {
				return nil, err
			}
			result = append(result, matches...)
			continue
		}
		matches, err := g.seekMatchesInPage(v.MatchTables)
		if err != nil {
			return nil, err
		}
		result = append(result, matches...)
	}
	return result, nil
}

func (g *GetCurrentStateOfTheTournament) seekMatchesInPage(mt []decomposition_tools.MatchTable) ([]content.CompletedMatch, error) {
	var matches []content.CompletedMatch

	//1 приоритет
	for _, v := range mt {
		//это отдельные таблицы с пиками/банами https://liquipedia.net/dota2/DreamLeague/Season_24/Group_Stage_1#Matches
		//плюс спаренные таблицы по кс - https://liquipedia.net/counterstrike/BLAST/Premier/2024/World_Final (Group stage)
		if v.Type == "bm" {
			for _, iv := range mt {
				if iv.Type == "bm" {
					bm, err := decomposition_tools.BrktsMatchlistMatchesPickOut(iv, g.Game)
					if err != nil {
						return nil, err
					}
					matches = append(matches, bm...)
				}
			}

			return matches, nil
		}
	}

	//2 приоритет
	for _, v := range mt {
		//это контент турнирной сетки в виде подтаблицы с пиками/банами https://liquipedia.net/dota2/DreamLeague/Season_24 (плейофф)
		if v.Type == "br" {
			for _, iv := range mt {
				bm, err := decomposition_tools.BrktsBracketMatchesPickOut(iv, g.Game)

				if err != nil {
					return nil, err
				}
				if iv.Type == "br" {
					matches = append(matches, bm...)
				}
			}

			return matches, nil
		}
	}
	//3 приоритет
	for _, v := range mt {
		//грубая попытка спарсить турнирную таблицу по таблице https://liquipedia.net/dota2/DreamLeague/Season_24 (playoff show shedule)
		if v.Type == "mc" {
			mr, err := decomposition_tools.WikitableMatchcardMatchesPickOut(v)
			if err != nil {
				return nil, err
			}

			matches = append(matches, mr...)

			return matches, nil
		}
	}

	//4 приоритет
	for _, v := range mt {
		//швейцарская система
		if v.Type == "sw" {
			sw, err := decomposition_tools.WikitableSwisstableMatchesPickOut(v)
			if err != nil {
				return nil, err
			}

			matches = append(matches, sw...)

			return matches, nil
		}
	}

	//5 приоритет
	for _, v := range mt {
		//грубая попытка спарсить таблицу с пересечениями https://liquipedia.net/dota2/DreamLeague/Season_24 (group A)
		if v.Type == "cr" {
			cr, err := decomposition_tools.TableCrosstableMatchesPickOut(v)
			if err != nil {
				return nil, err
			}

			matches = append(matches, cr...)

			return matches, nil
		}
	}

	return nil, nil
}
