package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetTeamName(dec *json.Decoder) string {
	for dec.More() {
		token, err := dec.Token()
		if err != nil {
			panic(err)
		}
		if token == "name" {
			token, err = dec.Token()
			return token.(string)
		}
	}
	return ""
}

func main() {
	teams := map[string]bool{
		"Germany":          false,
		"England":          false,
		"France":           false,
		"Spain":            false,
		"Arsenal":          false,
		"Chelsea":          false,
		"Barcelona":        false,
		"Real Madrid":      false,
		"Manchester Utd":   false,
		"FC Bayern Munich": false,
	}
	for i := 1; i < 100; i++ {
		url := fmt.Sprintf("https://vintagemonster.onefootball.com/api/teams/en/%v.json", i)
		resp, err := http.Get(url)
		if err != nil {
			panic(err)
		}
		dec := json.NewDecoder(resp.Body)
		if err != nil {
			panic(err)
		}
		if team, ok := teams[GetTeamName(dec)]; ok {
			fmt.Println(team)
		}
	}

}
