package main

import (
	"encoding/json"
	"io/ioutil"
)

var Configs Configuration

type Configuration struct {
	RiotApiKey        string
	GameVersion       string
	ImgurClientID     string
	ImgurClientSecret string
	DiscordWebhook    string
}

func SetupConfiguration() {
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	var config Configuration
	err = json.Unmarshal(file, &config)
	if err != nil {
		panic(err)
	}

	Configs = config
}
