package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

func TestGetAttribute(t *testing.T) {
	url := "https://vintagemonster.onefootball.com/api/teams/en/1.json"
	resp, err := http.Get(url)
	if err != nil {
		t.Error(err)
	}
	teamName := GetAttribute(json.NewDecoder(resp.Body), "name")
	if teamName != "Apoel FC" {
		t.FailNow()
	}
}

func TestIsFinished(t *testing.T) {
	url := "https://vintagemonster.onefootball.com/api/teams/en/10000.json"
	resp, err := http.Get(url)
	if err != nil {
		t.Error(err)
	}
	finished := IsFinished(json.NewDecoder(resp.Body))
	if !finished {
		t.FailNow()
	}
}

func TestGetTeamData(t *testing.T) {
	url := "https://vintagemonster.onefootball.com/api/teams/en/1.json"
	resp, err := http.Get(url)
	if err != nil {
		t.Error(err)
	}
	teamPlayers := GetTeamData(json.NewDecoder(resp.Body))
	if len(teamPlayers) != 34 && fmt.Sprintf("%T", teamPlayers[0]) == "Player" {
		t.FailNow()
	}
}
