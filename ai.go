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
	Choices []struct {
		Text         string `json:"text"`
		Index        int    `json:"index"`
		Logprobs     string `json:"logprobs"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

type PromptRequest struct {
	Model            string      `json:"model"`
	Prompt           string      `json:"prompt"`
	MaxTokens        int         `json:"max_tokens"`
	Temperature      float32     `json:"temperature"`
	TopP             float32     `json:"top_p"`
	N                int         `json:"n"`
	Stream           bool        `json:"stream"`
	Logprobs         interface{} `json:"logprobs"`
	FrequencyPenalty float32     `json:"frequency_penalty"`
	PresencePenalty  float32     `json:"presence_penalty"`
}

func GeneratePlayerMatchCommentUsingOpenAi(matchInfo PlayerMatchInfo) (result string, err error) {

	url := "https://api.openai.com/v1/completions"

	tones := []string{"debochada", "sarcástica", "irônica", "cínica", "satírica", "contente", "encorajadora", "sassy", "interpretando o jogador como um gênio"}
	randomTone := rand.Intn(len(tones))

	prompt := fmt.Sprintf("Escreva uma pequena avaliação (no pretérito) %s deste jogador considerando sua performance numa partida de League of Legends com estes dados. Nome do jogador: %sCampeão: %s	Kills: %d Deaths: %d Assists: %d Venceu: %t	Modo de jogo: %s Dano total: %d	Cura total: %d", tones[randomTone], matchInfo.PlayerName, matchInfo.Champion, matchInfo.Kills, matchInfo.Deaths, matchInfo.Assists, matchInfo.Win, matchInfo.GameMode, matchInfo.TotalDamageDealt, matchInfo.TotalHeal)

	if matchInfo.GameMode != "ARAM" {
		prompt = fmt.Sprintf("%s Posição: %s", prompt, matchInfo.TeamPosition)
	}

	promptReq := PromptRequest{
		Model:            "text-davinci-003",
		Prompt:           prompt,
		MaxTokens:        700,
		Temperature:      0.9,
		TopP:             1,
		N:                1,
		Stream:           false,
		Logprobs:         nil,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
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

	return responseObject.Choices[0].Text, nil
}
