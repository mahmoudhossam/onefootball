package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Player struct {
	Id        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Age       string `json:"age"`
}

type Team struct {
	Name       string   `json:"name"`
	IsNational bool     `json:"isNational"`
	Players    []Player `json:"players"`
}

func GetAttribute(dec *json.Decoder, attribute string) string {
	for dec.More() {
		token, err := dec.Token()
		if err != nil {
			panic(err)
		}
		if token == attribute {
			token, err = dec.Token()
			return token.(string)
		}
	}
	// attribute not found
	return ""
}

func GetTeamData(dec *json.Decoder) (players []Player) {
	for dec.More() {
		token := NextToken(dec)
		if token == "team" {
			var team Team
			err := dec.Decode(&team)
			if err != nil {
				panic(err)
			}
			return team.Players
		}
	}
	return
}

func IsFinished(dec *json.Decoder) bool {
	for dec.More() {
		token := NextToken(dec)
		if token == "code" {
			token := NextToken(dec)
			if token != 0.0 {
				return true
			}
		}
	}
	return false
}

func NextToken(dec *json.Decoder) (token json.Token) {
	token, err := dec.Token()
	if err != nil {
		panic(err)
	}
	return
}

func main() {
	teams := map[string][]Player{
		"Germany":          []Player{},
		"England":          []Player{},
		"France":           []Player{},
		"Spain":            []Player{},
		"Arsenal":          []Player{},
		"Chelsea":          []Player{},
		"Barcelona":        []Player{},
		"Real Madrid":      []Player{},
		"Manchester Utd":   []Player{},
		"FC Bayern Munich": []Player{},
	}
	finished := false
	intId := 1
	for !finished {
		url := fmt.Sprintf("https://vintagemonster.onefootball.com/api/teams/en/%v.json", intId)
		resp, err := http.Get(url)
		if err != nil {
			panic(err)
		}
		dec := json.NewDecoder(resp.Body)
		if err != nil {
			panic(err)
		}
		if team, ok := teams[GetAttribute(dec, "name")]; ok {
			fmt.Println(team)
		}
		finished = IsFinished(dec)
		intId++
	}

}
