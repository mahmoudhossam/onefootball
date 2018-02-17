package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

var teams = []string{"Germany", "England", "France", "Spain", "Arsenal", "Chelsea", "Barcelona", "Real Madrid", "Manchester Utd", "FC Bayern Munich"}
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
	defer wg.Done()
	body := <-in
	defer body.Close()
	dec := json.NewDecoder(body)
	teamName := GetAttribute(dec, "name")
	for _, t := range teams {
		if teamName == t {
			GetTeamData(teamName, dec)
		}
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
	for _, p := range roster {
		fmt.Printf("%v; %v; %v.\n", p.Name, p.Age, strings.Join(p.Teams, ","))
	}
}
