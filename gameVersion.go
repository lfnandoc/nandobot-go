package main

import (
	"encoding/json"
)

func GetGameVersion() string {
	url := "https://ddragon.leagueoflegends.com/api/versions.json"
	response, err := Get(url)
	if err != nil {
		return Configs.GameVersion
	}

	var gameVersions []string
	json.Unmarshal(response, &gameVersions)
	firstGameVersion := gameVersions[0]

	return firstGameVersion
}
