package liquiwrap

import (
	"errors"
	"fmt"
	"github.com/mnlght/liquiwrap/internal"
	"github.com/mnlght/liquiwrap/internal/decomposition_tools"
	"github.com/tdewolff/minify/v2/minify"
	"strings"
)

type GetOngoingMatchesByGame struct {
	Game string
	Tier int
}

func NewGetOngoingMatchesByGame(game string, tier int) *GetOngoingMatchesByGame {
	return &GetOngoingMatchesByGame{Game: game, Tier: tier}
}

func (g *GetOngoingMatchesByGame) Action() ([]OngoingMatch, error) {
	if g.Game == "counterstrike" {
		return g.getForCS()
	}

	filter := g.formLQFilter()
	if filter == "" {
		return nil, errors.New("game is not supported")
	}

	return g.getForParsebleGame(filter)
}

func (g *GetOngoingMatchesByGame) formLQFilter() string {
	switch g.Game {
	case "dota2":
		return fmt.Sprintf("{{MainPageMatches/Upcoming|filterbuttons-liquipediatier=%d|filterbuttons-liquipediatiertype=monthly,weekly,qualifier,misc,showmatch,national}}", g.Tier)
	case "leagueoflegends":
		return fmt.Sprintf("{{MainPageMatches/Upcoming|filterbuttons-liquipediatier=%d|filterbuttons-region=Korea,China,Europe,Turkey,Arab States,Brazil,Latin America South,North America,Latin America North,Oceania,Vietnam,Japan,Taiwan,Other}}", g.Tier)
	}
	return ""
}

func (g *GetOngoingMatchesByGame) getForCS() ([]OngoingMatch, error) {
	pz := internal.NewLiquipediaPageClient("/counterstrike/Liquipedia:Matches")
	pageContent, err := pz.Do()
	if err != nil {
		return nil, err
	}
	minifiedPageContent, err := minify.HTML(string(pageContent))
	if err != nil {
		fmt.Println(err)
	}

	res, err := decomposition_tools.OngoingMatchesCSPickOut(strings.NewReader(minifiedPageContent))
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (g *GetOngoingMatchesByGame) getForParsebleGame(filter string) ([]OngoingMatch, error) {
	client := internal.NewLiquipediaApiClient(internal.LiquipediaApiClientParams{
		Query: map[string]string{
			"action":       "parse",
			"format":       "json",
			"contentmodel": "wikitext",
			"maxage":       "600",
			"smaxage":      "600",
			"uselang":      "content",
			"prop":         "text",
			"text":         filter,
		},
		Game: g.Game,
	})

	resp, err := client.Do()
	ongoingMatchesRespCh := LiquipediaResponse{
		Body:  resp,
		Error: err,
	}
	ongoingMatchesResp := ongoingMatchesRespCh
	if ongoingMatchesResp.Error != nil {
		return nil, fmt.Errorf("%s, %e", "get ongoing matches error", ongoingMatchesResp.Error)
	}

	formattedContent, err := internal.BuildContent(ongoingMatchesResp.Body)
	if err != nil {
		return nil, err
	}

	matches, err := decomposition_tools.OngoingMatchesPickOut(strings.NewReader(formattedContent), g.Game)
	if err != nil {
		return nil, err
	}

	return matches, nil
}
