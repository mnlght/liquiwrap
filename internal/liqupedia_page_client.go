package internal

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

type LiquipediaPageClient struct {
	Page string
}

func NewLiquipediaPageClient(page string) *LiquipediaPageClient {
	return &LiquipediaPageClient{Page: page}
}

func (lpc *LiquipediaPageClient) Do() ([]byte, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://liquipedia.net%s", lpc.Page), nil)
	if err != nil {
		return nil, err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 300 {
		return nil, errors.New(fmt.Sprintf("http query error - %d", resp.StatusCode))
	}

	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
