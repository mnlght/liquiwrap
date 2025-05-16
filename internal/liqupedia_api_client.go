package internal

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const BaseLiquipediaUrl = "https://liquipedia.net/%s/api.php"

type LiquipediaApiClientParams struct {
	Query map[string]string
	Game  string
}

type LiquipediaApiClient struct {
	Params LiquipediaApiClientParams
}

func NewLiquipediaApiClient(p LiquipediaApiClientParams) *LiquipediaApiClient {
	return &LiquipediaApiClient{Params: p}
}

func (lc *LiquipediaApiClient) Do() ([]byte, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(BaseLiquipediaUrl, lc.Params.Game), nil)
	if err != nil {
		return nil, err
	}

	queryParams := req.URL.Query()
	for k, v := range lc.Params.Query {
		queryParams.Add(k, v)
	}
	req.URL.RawQuery = queryParams.Encode()
	req.Header.Add("Accept", "application/json; charset=utf-8")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("User-Agent", "GORN_PW/0.1")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 300 {
		return nil, errors.New(fmt.Sprintf("http query error - %d", resp.StatusCode))
	}

	defer resp.Body.Close()

	enc, err := gzip.NewReader(resp.Body)

	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(enc)
	if err != nil {
		return nil, err
	}

	return body, nil
}
