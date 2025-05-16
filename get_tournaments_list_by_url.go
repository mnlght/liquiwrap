package liquiwrap

import (
	"fmt"
	"github.com/mnlght/liquiwrap/content"
	"github.com/mnlght/liquiwrap/internal"
	"github.com/mnlght/liquiwrap/internal/decomposition_tools"
	"strings"
)

type GetTournamentListByUrl struct {
	Url string
}

func NewGetTournamentListByUrl(url string) *GetTournamentListByUrl {
	return &GetTournamentListByUrl{Url: url}
}

func (g *GetTournamentListByUrl) Action() ([]content.Tournament, error) {
	meta, err := decomposition_tools.GetMetaFromUrl(g.Url)
	if err != nil {
		return nil, err
	}

	pz := internal.NewLiquipediaPageClient(fmt.Sprintf("/%s/%s", meta.Game, meta.PageToParse))
	pageContentExternal, err := pz.Do()
	if err != nil {
		return nil, err
	}

	tournaments := decomposition_tools.CollectTournamentMetaList(strings.NewReader(string(pageContentExternal)), meta.Game)
	return tournaments, nil
}
