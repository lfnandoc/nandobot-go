package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
)

type CompletionResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Usage   struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
		Index        int    `json:"index"`
	} `json:"choices"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GptRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

func GeneratePlayerMatchCommentUsingOpenAi(matchInfo PlayerMatchInfo) (result string, err error) {

	url := "https://api.openai.com/v1/chat/completions"

	tones := []string{"debochada", "sarcástica", "irônica", "encorajadora", "perversa", "malvada", "inspiracional", "desmotivacional", "cruel"}
	randomTone := rand.Intn(len(tones))

	prompt := fmt.Sprintf("Escreva uma pequena história %s da performance deste jogador numa partida de League of Legends com estes dados: Nome do jogador: %s, Campeão: %s, Kills: %d, Deaths: %d, Assists: %d, Venceu: %t, Modo de jogo: %s, Dano total: %d, Cura total: %d", tones[randomTone], matchInfo.PlayerName, matchInfo.Champion, matchInfo.Kills, matchInfo.Deaths, matchInfo.Assists, matchInfo.Win, matchInfo.GameMode, matchInfo.TotalDamageDealt, matchInfo.TotalHeal)

	if matchInfo.GameMode != "ARAM" {
		prompt = fmt.Sprintf("%s Posição: %s", prompt, matchInfo.TeamPosition)
	}

	promptReq := GptRequest{
		Model: "gpt-3.5-turbo",
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonReq, err := json.Marshal(promptReq)

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonReq))
	if err != nil {
		return "", err
	}

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Set("Authorization", "Bearer "+Configs.OpenAiApiKey)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}

	response_body, err := ioutil.ReadAll(response.Body)
	println(string(response_body))

	var responseObject CompletionResponse
	json.Unmarshal(response_body, &responseObject)

	return responseObject.Choices[0].Message.Content, nil
}
