package decomposition_tools

import (
	"errors"
	"fmt"
	"strings"
)

type MetaFromUrl struct {
	Game        string
	Tier        string
	PageToParse string
}

func GetMetaFromUrl(url string) (*MetaFromUrl, error) {
	ms := strings.Split(url, "liquipedia.net")
	if len(ms) > 1 {
		if ms[1] != "" {
			mm := strings.Split(ms[1], "/")
			if len(mm) > 1 {
				t := strings.Replace(mm[2], "_", " ", -1)
				if len(mm) < 4 {
					return &MetaFromUrl{
						Game:        mm[1],
						Tier:        t,
						PageToParse: mm[2],
					}, nil
				}

				if len(mm) == 4 {
					return &MetaFromUrl{
						Game:        mm[1],
						Tier:        t,
						PageToParse: fmt.Sprintf("%s/%s", mm[2], mm[3]),
					}, nil
				}
			}

		}
	}

	return nil, errors.New("no meta info")
}
