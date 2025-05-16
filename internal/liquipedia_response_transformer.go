package internal

import (
	"encoding/json"
)

type SimpleLQAnswer struct {
	Parse struct {
		Title  string `json:"title"`
		PageId int    `json:"pageid"`
		Text   struct {
			Content string `json:"*"`
		} `json:"text"`
	} `json:"parse"`
}

func BuildContent(data []byte) (string, error) {
	r := &SimpleLQAnswer{}

	err := json.Unmarshal(data, r)

	if err != nil {
		return "", err
	}

	return r.Parse.Text.Content, nil
}
