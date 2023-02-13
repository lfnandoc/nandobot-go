package main

import (
	"fmt"
	"github.com/robfig/cron"
	"math/rand"
	"time"
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

			// Randomly send a comment
			rand.Seed(time.Now().UnixNano())
			sendComment := rand.Intn(100) < 25
			var commentMessage DiscordMessage

			if sendComment {
				commentMessage = GenerateCommentDiscordMessage(*matchInfo)
			}

			err = SendDiscordMessageToWebhook(discordMessage)
			if err == nil {
				matchInfo.MessageSent = true
				DB.Save(&matchInfo)
				if sendComment {
					SendDiscordMessageToWebhook(commentMessage)
				}
			}
		}
	}
}
