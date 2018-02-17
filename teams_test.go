package main

import (
	"encoding/json"
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

func TestNextToken(t *testing.T) {
	url := "https://vintagemonster.onefootball.com/api/teams/en/10000.json"
	resp, err := http.Get(url)
	if err != nil {
		t.Error(err)
	}
	dec := json.NewDecoder(resp.Body)
	token := NextToken(dec)
	if token.(json.Delim) != '{' {
		t.FailNow()
	}
}
