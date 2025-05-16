package liquiwrap

import (
	"fmt"
	"github.com/mnlght/liquiwrap/internal"
	decomposition_tools2 "github.com/mnlght/liquiwrap/internal/decomposition_tools"
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

func (g *GetCurrentStateOfTheTournament) Action() ([]CompletedMatch, error) {
	p := internal.NewLiquipediaPageClient(fmt.Sprintf("/%s/%s", g.Game, g.Url))
	pageContent, err := p.Do()
	if err != nil {
		return nil, err
	}

	mts, err := decomposition_tools2.ClassifyTypeOfTournamentStage(strings.NewReader(string(pageContent)))
	if err != nil {
		return nil, err
	}

	var result []CompletedMatch
	for _, v := range mts {
		if v.External {
			pz := internal.NewLiquipediaPageClient(v.ExternalLink)
			pageContentExternal, err := pz.Do()
			if err != nil {
				return nil, err
			}

			stagesExternal, err := decomposition_tools2.ClassifyTypeOfTournamentStage(strings.NewReader(string(pageContentExternal)))
			var allExternalMatchBlocks []decomposition_tools2.MatchTable
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

func (g *GetCurrentStateOfTheTournament) seekMatchesInPage(mt []decomposition_tools2.MatchTable) ([]CompletedMatch, error) {
	var matches []CompletedMatch

	//1 приоритет
	for _, v := range mt {
		//это отдельные таблицы с пиками/банами https://liquipedia.net/dota2/DreamLeague/Season_24/Group_Stage_1#Matches
		//плюс спаренные таблицы по кс - https://liquipedia.net/counterstrike/BLAST/Premier/2024/World_Final (Group stage)
		if v.Type == "bm" {
			fmt.Println("bm found")
			for _, iv := range mt {
				if iv.Type == "bm" {
					bm, err := decomposition_tools2.BrktsMatchlistMatchesPickOut(iv, g.Game)
					fmt.Println(len(bm))
					if err != nil {
						fmt.Println(err)
						panic("")
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
			fmt.Println("brfound")
			for _, iv := range mt {
				bm, err := decomposition_tools2.BrktsBracketMatchesPickOut(iv, g.Game)

				if err != nil {
					fmt.Println(err)
					panic("")
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
			fmt.Println("mc")
			mr, err := decomposition_tools2.WikitableMatchcardMatchesPickOut(v)
			if err != nil {
				fmt.Println(err)
				panic("")
			}

			//fmt.Println("MR MATCHES", mr)
			matches = append(matches, mr...)

			return matches, nil
		}
	}

	//4 приоритет
	for _, v := range mt {
		//швейцарская система
		if v.Type == "sw" {
			sw, err := decomposition_tools2.WikitableSwisstableMatchesPickOut(v)
			if err != nil {
				fmt.Println(err)
				panic("")
			}

			//fmt.Println("SWISS MATCHES", sw)
			matches = append(matches, sw...)

			return matches, nil
		}
	}

	//5 приоритет
	for _, v := range mt {
		//грубая попытка спарсить таблицу с пересечениями https://liquipedia.net/dota2/DreamLeague/Season_24 (group A)
		if v.Type == "cr" {
			cr, err := decomposition_tools2.TableCrosstableMatchesPickOut(v)
			if err != nil {
				fmt.Println(err)
				panic("")
			}

			//fmt.Println("CROSS MATCHES", cr)
			matches = append(matches, cr...)

			return matches, nil
		}
	}

	return nil, nil
}
