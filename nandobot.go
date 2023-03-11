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

			// Log message sent
			fmt.Println("Sending message to discord - Match ID: " + matchInfo.MatchID + " - Player ID: " + fmt.Sprint(matchInfo.PlayerID) + " - Player Name: " + matchInfo.PlayerName)

			// Randomly send a comment
			rand.Seed(time.Now().UnixNano())
			sendComment := rand.Intn(100) < 25
			var commentMessage string
			var commentError error

			if sendComment {
				commentMessage, commentError = GeneratePlayerMatchCommentUsingOpenAi(*matchInfo)
				if commentError == nil {
					discordMessage.Embeds[0].Footer = Footer{Text: commentMessage}
				}
			}

			err = SendDiscordMessageToWebhook(discordMessage)
			if err == nil {
				matchInfo.MessageSent = true
				DB.Save(&matchInfo)
			}
		}
	}
}
