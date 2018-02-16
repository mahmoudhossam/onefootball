package main

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestGetTeamName(t *testing.T) {
	url := "https://vintagemonster.onefootball.com/api/teams/en/1.json"
	resp, err := http.Get(url)
	if err != nil {
		t.Error(err)
	}
	team := GetTeamName(json.NewDecoder(resp.Body))
	if team != "Apoel FC" {
		t.Fail()
	}
}
