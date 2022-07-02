package main

import (
	"fmt"
	"time"

	"github.com/robfig/cron"
)

func main() {
	SetupConfiguration()
	SetupDatabase()
	c := cron.New()
	c.AddFunc("@every 2m", func() {
		UpdatePlayerMatches()
	})
	c.Start()

	SetupApi()
}

func UpdatePlayerMatches() {
	var players []Player
	DB.Find(&players)
	for _, player := range players {
		time.Sleep(time.Second * 2)
		matchInfo, err := GetLastMatchInfo(fmt.Sprint(player.ID))
		if err != nil {
			fmt.Println(err)
			continue
		}

		if !matchInfo.MessageSent {
			discordMessage := GenerateDiscordMessage(*matchInfo)
			err = SendDiscordMessageToWebhook(discordMessage)
			if err == nil {
				matchInfo.MessageSent = true
				DB.Save(&matchInfo)
			}
		}
	}
}
