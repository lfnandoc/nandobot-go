package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type DiscordMessage struct {
	Username   string        `json:"username,omitempty"`
	AvatarURL  string        `json:"avatar_url,omitempty"`
	Content    string        `json:"content,omitempty"`
	Embeds     []Embed       `json:"embeds,omitempty"`
	Components []interface{} `json:"components,omitempty"`
}

type Embed struct {
	Title       string    `json:"title"`
	Color       int       `json:"color"`
	Description string    `json:"description"`
	Timestamp   time.Time `json:"timestamp"`
	URL         string    `json:"url"`
	Author      struct {
	} `json:"author,omitempty"`
	Image     Image     `json:"image"`
	Thumbnail Thumbnail `json:"thumbnail"`
	Footer    Footer    `json:"footer,omitempty"`
	Fields    []Field   `json:"fields"`
}

type Footer struct {
	Text string `json:"text"`
}

type Image struct {
	URL string `json:"url"`
}

type Thumbnail struct {
	URL string `json:"url"`
}

type Field struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

func GenerateDiscordMessage(matchInfo PlayerMatchInfo) DiscordMessage {

	color := 2329825
	if !matchInfo.Win {
		color = 13247783
	}

	matchResult := "ganhou"
	if !matchInfo.Win {
		matchResult = "perdeu"
	}

	embed := Embed{
		Color:       color,
		Title:       matchInfo.PlayerName,
		Description: fmt.Sprintf("Jogou de %s e %s!", matchInfo.Champion, matchResult),
		Thumbnail:   Thumbnail{URL: fmt.Sprintf("http://ddragon.leagueoflegends.com/cdn/%s/img/champion/%s.png", GetGameVersion(), matchInfo.Champion)},
		Timestamp:   matchInfo.EndTime,
		URL:         fmt.Sprintf("https://blitz.gg/lol/match/br1/%s/%s", url.QueryEscape(matchInfo.PlayerName), strings.Replace(matchInfo.MatchID, "BR1_", "", -1)),
		Image:       Image{URL: matchInfo.ItemListImageLink},
	}

	embed.Fields = []Field{
		{
			Name:   "KDA",
			Value:  fmt.Sprintf("%d/%d/%d", matchInfo.Kills, matchInfo.Deaths, matchInfo.Assists),
			Inline: true,
		},
		{
			Name:   "Duração",
			Value:  fmt.Sprintf("%d minutos", matchInfo.Duration/60),
			Inline: true,
		},
		{
			Name:   "Modo",
			Value:  matchInfo.GameMode,
			Inline: true,
		},
		{
			Name:   "Dano Total",
			Value:  fmt.Sprint(matchInfo.TotalDamageDealt),
			Inline: true,
		},
		{
			Name:   "Cura Total",
			Value:  fmt.Sprint(matchInfo.TotalHeal),
			Inline: true,
		}}

	if matchInfo.GameMode != "ARAM" {
		embed.Fields = append(embed.Fields, Field{
			Name:   "Posição",
			Value:  matchInfo.TeamPosition,
			Inline: true})
	}

	discordMessage := DiscordMessage{
		Embeds: []Embed{embed},
	}

	return discordMessage
}

func GenerateCommentDiscordMessage(matchInfo PlayerMatchInfo) DiscordMessage {

	comment, err := GeneratePlayerMatchCommentUsingOpenAi(matchInfo)
	if err != nil {
		comment = "Erro ao gerar comentário. Maldita inteligência artificial."
	}

	discordMessage := DiscordMessage{
		Content: comment,
	}

	return discordMessage
}

func SendDiscordMessageToWebhook(discordMessage DiscordMessage) (err error) {

	data, err := json.Marshal(discordMessage)
	if err != nil {
		return err
	}

	response, err := http.Post(Configs.DiscordWebhook, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	if response.StatusCode != 204 {
		return fmt.Errorf("discord webhook returned %d", response.StatusCode)
	}

	return nil
}
