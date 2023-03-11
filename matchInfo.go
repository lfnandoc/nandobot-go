package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/lib/pq"
)

func GetLastMatchInfo(id string) (matchInfo *PlayerMatchInfo, err error) {
	var player Player
	DB.Preload("Matches").Find(&player, "id = ?", id)
	if player.ID == 0 {
		return nil, fmt.Errorf("Player does not exist")
	}

	matchId, err := GetLastMatchId(player.UserId)
	if err != nil {
		return nil, err
	}

	existingMatch := CheckIfMatchIsAlreadyRecordered(player.Matches, *matchId)
	if existingMatch != nil {
		return existingMatch, nil
	}

	url := fmt.Sprintf("https://americas.api.riotgames.com/lol/match/v5/matches/%s?api_key=%s", *matchId, Configs.RiotApiKey)
	response, err := Get(url)
	if err != nil {
		return nil, err
	}

	var match MatchInfo
	json.Unmarshal(response, &match)
	playerData := GetParticipantData(match, player.UserId)
	items := [6]int64{int64(playerData.Item0), int64(playerData.Item1), int64(playerData.Item2), int64(playerData.Item3), int64(playerData.Item4), int64(playerData.Item5)}

	playerMatch := PlayerMatchInfo{
		MatchID:          *matchId,
		PlayerID:         player.ID,
		PlayerName:       playerData.SummonerName,
		Champion:         playerData.ChampionName,
		Kills:            playerData.Kills,
		Deaths:           playerData.Deaths,
		Assists:          playerData.Assists,
		GameMode:         strings.Replace(match.Info.GameMode, "CLASSIC", "Normal", -1),
		Win:              playerData.Win,
		EndTime:          time.UnixMilli(match.Info.GameEndTimestamp),
		Duration:         match.Info.GameDuration,
		TotalDamageDealt: playerData.TotalDamageDealt,
		TotalHeal:        playerData.TotalHeal,
		MythicItem:       playerData.Challenges.MythicItemUsed,
		Items:            pq.Int64Array{items[0], items[1], items[2], items[3], items[4], items[5]},
	}

	if playerMatch.GameMode == "Normal" {
		position := strings.Replace(playerData.TeamPosition, "UTILITY", "SUPORTE", -1)
		position = strings.Replace(position, "MIDDLE", "MID", -1)
		playerMatch.TeamPosition = cases.Title(language.Und).String(position)
	}

	itemListImage, err := MergeItemImages(playerMatch.Items)
	if err == nil {
		playerMatch.ItemListImageLink = *itemListImage
	}

	DB.Create(&playerMatch)

	return &playerMatch, nil
}

func GetParticipantData(match MatchInfo, uuid string) (participant *ParticipantInfo) {

	for _, participant := range match.Info.Participants {
		if participant.Puuid == uuid {
			return &participant
		}
	}

	return nil
}

func CheckIfMatchIsAlreadyRecordered(matches []PlayerMatchInfo, id string) *PlayerMatchInfo {
	for _, match := range matches {
		if match.MatchID == id {
			return &match
		}
	}
	return nil
}

func GetLastMatchId(uuid string) (matchId *string, err error) {
	url := fmt.Sprintf("https://americas.api.riotgames.com/lol/match/v5/matches/by-puuid/%s/ids?start=0&count=1&api_key=%s", uuid, Configs.RiotApiKey)
	response, err := Get(url)
	if err != nil {
		return nil, err
	}
	var matchIds []string
	json.Unmarshal(response, &matchIds)

	return &matchIds[0], nil
}
