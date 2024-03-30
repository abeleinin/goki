package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/charmbracelet/bubbles/list"
)

type ChatCompletion struct {
	ID                string   `json:"id"`
	Object            string   `json:"object"`
	Created           int      `json:"created"`
	Model             string   `json:"model"`
	Choices           []Choice `json:"choices"`
	Usage             Usage    `json:"usage"`
	SystemFingerprint string   `json:"system_fingerprint"`
}

type Choice struct {
	Index        int              `json:"index"`
	Message      Message          `json:"message"`
	Logprobs     *json.RawMessage `json:"logprobs"`
	FinishReason string           `json:"finish_reason"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type Response struct {
	DeckTitle  string     `json:"title"`
	FlashCards []CardInfo `json:"flashcards"`
}

type CardInfo struct {
	Front string `json:"front"`
	Back  string `json:"back"`
}

func gptClient(prompt string) *Deck {
	apiKey := os.Getenv("OPENAI_API_KEY")
	url := "https://api.openai.com/v1/chat/completions"

	systemPrompt := "You are a helpful flashcard making assistant. Given topic, category, or concept generate a JSON object with 'title' (string) and 'flashcards' (object) with 'front' and 'back' values containing flashcard data."
	message := [...]map[string]interface{}{
		{"role": "system", "content": systemPrompt},
		{"role": "user", "content": prompt},
	}

	data := map[string]interface{}{
		"model":           "gpt-4-turbo-preview",
		"messages":        message,
		"response_format": map[string]string{"type": "json_object"},
		"temperature":     0.5,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}

	req.Header.Set("OpenAI-Beta", "assistants=v1")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// Test saved requests
	// body := ``

	// Save output from api
	// os.WriteFile("gpt.json", body, 0644)

	var chatCompletion ChatCompletion
	err = json.Unmarshal([]byte(body), &chatCompletion)
	if err != nil {
		fmt.Println("error unmarshaling JSON:", err)
	}

	var content Response
	err = json.Unmarshal([]byte(chatCompletion.Choices[0].Message.Content), &content)
	if err != nil {
		fmt.Println("error unmarshaling JSON:", err)
	}

	cards := []list.Item{}
	for _, cardInfo := range content.FlashCards {
		card := NewCard(cardInfo.Front, cardInfo.Back)
		cards = append(cards, card)
	}

	return NewDeck(content.DeckTitle, cards)
}
