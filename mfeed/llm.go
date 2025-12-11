package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"text/template"

	openai "github.com/sashabaranov/go-openai"
)

// LLMClient wraps the OpenAI API client
type LLMClient struct {
	client *openai.Client
	config *Config
}

// NewLLMClient creates a new OpenAI client
func NewLLMClient(cfg Config) *LLMClient {
	apiKey, err := fetchAPIKey()
	if err != nil {
		log.Fatalf("failed to get API key: %v", err)
	}
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = cfg.OpenAIEndpoint
	return &LLMClient{
		client: openai.NewClientWithConfig(config),
		config: &cfg,
	}
}

type Scoring struct {
	Articles []ItemScore `json:"articles"`
}

type ItemScore struct {
	ID    int `json:"id"`
	Score int `json:"score"`
}

// ScoreArticles sends item titles to OpenAI API for scoring
func (c *LLMClient) ScoreArticles(items []FeedItem) ([]FeedItem, error) {
	ctx := context.Background()

	rankPrompt, err := preparePrompt("prompts/rank.tmpl", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare prompt: %w", err)
	}

	data := struct {
		Items []FeedItem
	}{
		Items: items,
	}
	itemsPrompt, err := preparePrompt("prompts/titles.tmpl", data)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare prompt: %w", err)
	}

	req := openai.ChatCompletionRequest{
		Model: c.config.LLMModel,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    "system",
				Content: rankPrompt,
			},
			{
				Role:    "user",
				Content: itemsPrompt,
			},
		},
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type: "json_object",
		},
		Temperature: 0.0,
		MaxTokens:   1500,
	}

	resp, err := c.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create chat completion: %w", err)
	}

	var scores Scoring
	err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), &scores)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal scoring response: %w", err)
	}

	for i := range items {
		for _, score := range scores.Articles {
			if items[i].ID == score.ID {
				items[i].Score = score.Score
				break
			}
		}
	}

	return items, nil
}

func (c *LLMClient) SummarizeContent(content string) (string, error) {
	ctx := context.Background()

	data := struct {
		Content string
	}{
		Content: content,
	}
	summaryPrompt, err := preparePrompt("prompts/summarize.tmpl", data)
	if err != nil {
		return "", fmt.Errorf("failed to prepare prompt: %w", err)
	}

	req := openai.ChatCompletionRequest{
		Model: c.config.LLMModel,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    "system",
				Content: summaryPrompt,
			},
			{
				Role:    "user",
				Content: content,
			},
		},
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type: "json_object",
		},
		Temperature: 0.5,
		MaxTokens:   800,
	}

	resp, err := c.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to create chat completion: %w", err)
	}

	var summary struct {
		Summary string `json:"summary"`
	}
	err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), &summary)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal scoring response: %w", err)
	}

	return summary.Summary, nil
}

func preparePrompt(promptFile string, data interface{}) (string, error) {
	tmpl, err := template.ParseFiles(promptFile)
	if err != nil {
		return "", fmt.Errorf("failed to read prompt template: %w", err)
	}
	tmplBuffer := bytes.Buffer{}
	err = tmpl.Execute(&tmplBuffer, data)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}
	return tmplBuffer.String(), nil
}

// fetchAPIKey retrieves the API key from environment variables
func fetchAPIKey() (string, error) {
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("GROQ_API_KEY environment variable not set")
	}
	return apiKey, nil
}
