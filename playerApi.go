package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupApi() {
	r := gin.Default()
	r.POST("/player", AddPlayer)
	r.GET("/player/:id", GetLastMatchInfoOfPlayer)
	r.Run()
}

func GetLastMatchInfoOfPlayer(c *gin.Context) {
	id := c.Param("id")
	matchInfo, err := GetLastMatchInfo(id)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, matchInfo)
}

func AddPlayer(c *gin.Context) {
	var player Player

	if err := c.ShouldBindJSON(&player); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uuid, err := GetPlayerUuid(player.Name)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	player.UserId = *uuid
	DB.Create(&player)
	c.IndentedJSON(http.StatusOK, player)
}

type SummonerResponse struct {
	ID            string `json:"id"`
	AccountID     string `json:"accountId"`
	Puuid         string `json:"puuid"`
	Name          string `json:"name"`
	ProfileIconID int    `json:"profileIconId"`
	RevisionDate  int64  `json:"revisionDate"`
	SummonerLevel int    `json:"summonerLevel"`
}

func GetPlayerUuid(name string) (uuid *string, err error) {
	url := fmt.Sprintf("https://br1.api.riotgames.com/lol/summoner/v4/summoners/by-name/%s?api_key=%s", name, Configs.RiotApiKey)
	response, err := Get(url)
	if err != nil {
		return nil, err
	}

	var summoner SummonerResponse
	json.Unmarshal(response, &summoner)
	if summoner.Puuid == "" {
		return nil, fmt.Errorf("Player does not exist on Riot")
	}

	return &summoner.Puuid, nil
}
