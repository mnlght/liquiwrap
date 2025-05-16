package liquiwrap

import (
	"errors"
	"fmt"
	"github.com/mnlght/liquiwrap/content"
	"github.com/mnlght/liquiwrap/internal"
	"github.com/mnlght/liquiwrap/internal/decomposition_tools"
	"github.com/tdewolff/minify/v2/minify"
	"strings"
)

type GetOngoingMatchesByGameWithBus struct {
	Game string
	Tier int
	bus  *LiquipediaBus
}

func NewGetOngoingMatchesByGameWithBus(game string, tier int, bus *LiquipediaBus) *GetOngoingMatchesByGameWithBus {
	return &GetOngoingMatchesByGameWithBus{Game: game, Tier: tier, bus: bus}
}

func (g *GetOngoingMatchesByGameWithBus) Action() ([]content.OngoingMatch, error) {
	if g.Game == "counterstrike" {
		return g.getForCS()
	}

	filter := g.formLQFilter()
	if filter == "" {
		return nil, errors.New("game is not supported")
	}

	return g.getForParsebleGame(filter)
}

func (g *GetOngoingMatchesByGameWithBus) formLQFilter() string {
	switch g.Game {
	case "dota2":
		return fmt.Sprintf("{{MainPageMatches/Upcoming|filterbuttons-liquipediatier=%d|filterbuttons-liquipediatiertype=monthly,weekly,qualifier,misc,showmatch,national}}", g.Tier)
	case "leagueoflegends":
		return fmt.Sprintf("{{MainPageMatches/Upcoming|filterbuttons-liquipediatier=%d|filterbuttons-region=Korea,China,Europe,Turkey,Arab States,Brazil,Latin America South,North America,Latin America North,Oceania,Vietnam,Japan,Taiwan,Other}}", g.Tier)
	}
	return ""
}

func (g *GetOngoingMatchesByGameWithBus) getForCS() ([]content.OngoingMatch, error) {
	pz := internal.NewLiquipediaPageClient("/counterstrike/Liquipedia:Matches")
	pageContent, err := pz.Do()
	if err != nil {
		return nil, err
	}
	minifiedPageContent, err := minify.HTML(string(pageContent))
	if err != nil {
		return nil, err
	}

	res, err := decomposition_tools.OngoingMatchesCSPickOut(strings.NewReader(minifiedPageContent))
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (g *GetOngoingMatchesByGameWithBus) getForParsebleGame(filter string) ([]content.OngoingMatch, error) {
	ongoingMatchesRespCh := make(chan LiquipediaResponse)
	g.bus.AddRequest(&LiquipediaRequest{
		Game: g.Game,
		Params: map[string]string{
			"action":       "parse",
			"format":       "json",
			"contentmodel": "wikitext",
			"maxage":       "600",
			"smaxage":      "600",
			"uselang":      "content",
			"prop":         "text",
			"text":         filter,
		},
		ResponseCh: ongoingMatchesRespCh,
	})
	ongoingMatchesResp := <-ongoingMatchesRespCh
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
