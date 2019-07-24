package persist

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type response struct {
	NextChangeID string        `json:"next_change_id"`
	Stashes      []interface{} `json:"stashes"`
}

type Item struct {
	ItemID        string `json:"item_id"`
	StashID       string `json:"stash_id"`
	Account       string `json:"account"`
	LastCharacter string `json:"last_character"`
	ItemData      string `json:"item_data"`
}

func GetStashData(ctx context.Context, league string) ([]Item, error) {
	res, err := http.Get("https://www.pathofexile.com/api/public-stash-tabs")
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	rs := &response{}
	if err := json.Unmarshal(body, rs); err != nil {
		return nil, err
	}
	fmt.Println(rs.NextChangeID)
	ret := make([]Item, 0)
	for _, stash := range rs.Stashes {
		var (
			itemID        string
			stashID       string
			account       string
			lastCharacter string
			itemData      string
		)
		stash := stash.(map[string]interface{})
		switch name := stash["accountName"].(type) {
		case string:
			account = name
		default:
			account = ""
		}
		switch char := stash["lastCharacterName"].(type) {
		case string:
			lastCharacter = char
		default:
			lastCharacter = ""
		}
		stashID = stash["id"].(string)
		items := stash["items"].([]interface{})
		if len(items) < 1 && !stash["public"].(bool) {
			continue
		}
		for _, item := range items {
			asserted := item.(map[string]interface{})
			note := asserted["note"]
			if asserted["league"] != league || note == "" {
				continue
			}
			itemID = asserted["id"].(string)
			itemBytes, err := json.Marshal(item)
			if err != nil {
				return nil, err
			}
			itemData = string(itemBytes)
		}
		item := Item{
			ItemID:        itemID,
			StashID:       stashID,
			Account:       account,
			LastCharacter: lastCharacter,
			ItemData:      itemData,
		}
		ret = append(ret, item)
	}
	return ret, nil
}
