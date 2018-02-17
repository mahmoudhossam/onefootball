package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"sync"
)

type Player struct {
	Id    string   `json:"id"`
	Name  string   `json:"name"`
	Age   string   `json:"age"`
	Teams []string `json:"teams"`
}

type Team struct {
	Name       string   `json:"name"`
	IsNational bool     `json:"isNational"`
	Players    []Player `json:"players"`
}

var roster map[string]Player

func GetAttribute(dec *json.Decoder, attribute string) string {
	for dec.More() {
		token := NextToken(dec)
		if token == attribute {
			token := NextToken(dec)
			return token.(string)
		}
	}
	// attribute not found
	return ""
}

func GetTeamData(team string, dec *json.Decoder) {
	// look for the players object
	for token := NextToken(dec); token != "players"; {
		token = NextToken(dec)
	}
	// parse start of the players object
	NextToken(dec)
	for dec.More() {
		var player Player
		err := dec.Decode(&player)
		if err != nil {
			panic(err)
		}
		if player.Teams == nil {
			player.Teams = make([]string, 2)
		}
		player.Teams = []string{team}
		if _, found := roster[player.Id]; found {
			newTeams := append(roster[player.Id].Teams, team)
			player.Teams = newTeams
		}
		roster[player.Id] = player
	}
}

func NextToken(dec *json.Decoder) (token json.Token) {
	token, err := dec.Token()
	if err != nil {
		panic(err)
	}
	return
}

func GetURL(url string, ch chan io.ReadCloser) {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		fmt.Println(url)
		panic(err)
	}
	ch <- resp.Body
}

func ProcessJSONData(in chan io.ReadCloser, wg *sync.WaitGroup) {
	teams := []string{"Germany", "England", "France", "Spain", "Arsenal", "Chelsea",
		"Barcelona", "Real Madrid", "Manchester Utd", "FC Bayern Munich"}
	defer wg.Done()
	body := <-in
	dec := json.NewDecoder(body)
	teamName := GetAttribute(dec, "name")
	for _, t := range teams {
		if teamName == t {
			GetTeamData(teamName, dec)
		}
	}
	if err := body.Close(); err != nil {
		panic(err)
	}
}

func main() {
	httpData := make(chan io.ReadCloser, 10)
	roster = make(map[string]Player)
	var wg sync.WaitGroup
	for id := 1; id < 100; id++ {
		url := fmt.Sprintf("https://vintagemonster.onefootball.com/api/teams/en/%v.json", id)
		go GetURL(url, httpData)
		wg.Add(1)
		go ProcessJSONData(httpData, &wg)
	}
	wg.Wait()
	close(httpData)
	players := make([]Player, 0)
	for _, p := range roster {
		players = append(players, p)
	}
	sort.Slice(players, func(i int, j int) bool {
		return players[i].Name < players[j].Name
	})
	for i, p := range players {
		fmt.Printf("%v. %v; %v; %v\n", i+1, p.Name, p.Age, strings.Join(p.Teams, ","))
	}
}
