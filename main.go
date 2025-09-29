package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/joho/godotenv"
)

const openRouterURL = "https://openrouter.ai/api/v1/chat/completions"

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type reasoningParams struct {
	Exclude bool `json:"exclude"`
}

type chatCompletionRequest struct {
	Model     string           `json:"model"`
	Messages  []chatMessage    `json:"messages"`
	Reasoning *reasoningParams `json:"reasoning,omitempty"`
}

type chatCompletionResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func loadAPIKey() string {
	key := os.Getenv("OPENROUTER_API_KEY")
	if key == "" {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("The environment variable OPENROUTER_API_KEY is not set, and there was an unexpected error loading .env file: %v", err)
		}
		key = os.Getenv("OPENROUTER_API_KEY")
	}

	if key == "" {
		log.Fatalf("OPENROUTER_API_KEY environment variable not set")
	}

	return key
}

func generateAnswer(ctx context.Context, apiKey string, s bool, q string) string {
	var systemPrompt string
	if s {
		systemPrompt = "You are a helpful and versatile assistant, especially skilled in programming and devops. The user is asking for a bash terminal command. Respond ONLY with the exact, executable bash command that accomplishes what they're asking for. Do not include explanations, markdown formatting, unnecessary quotation marks, or any text other than the command itself. The command will be automatically executed in the user's terminal."
	} else {
		systemPrompt = "You are a helpful and versatile assistant, especially skilled in programming and devops. Answer the user's question(s) in a very concise manner, unless it requires a long response or the user specifically asks you for one."
	}

	payload := chatCompletionRequest{
		Model: "x-ai/grok-4-fast",
		Messages: []chatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: q},
		},
		Reasoning: &reasoningParams{Exclude: true},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("Failed to encode request: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, openRouterURL, bytes.NewReader(body))
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Failed to call OpenRouter: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}

	if resp.StatusCode >= 300 {
		log.Fatalf("OpenRouter request failed with status %d: %s", resp.StatusCode, strings.TrimSpace(string(respBody)))
	}

	var parsed chatCompletionResponse
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		log.Fatalf("Failed to decode response: %v", err)
	}

	if parsed.Error != nil {
		log.Fatalf("OpenRouter error: %s", parsed.Error.Message)
	}

	if len(parsed.Choices) > 0 {
		result := parsed.Choices[0].Message.Content
		if s {
			result = strings.TrimSpace(result)
			if err := clipboard.WriteAll(result); err != nil {
				log.Fatalf("Failed to copy to clipboard: %v", err)
			}
			return fmt.Sprintf("%s\nCopied to clipboard", result)
		}
		return result
	}

	log.Fatal("No valid response")
	return ""
}

func main() {
	ctx := context.Background()
	apiKey := loadAPIKey()

	sFlag := flag.Bool("s", false, "Returns a valid bash prompt that matches the request.")

	flag.Parse()

	args := flag.Args()

	if len(args) == 0 {
		err := fmt.Errorf("no arguments provided")
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	question := strings.Join(args, " ")

	answer := generateAnswer(ctx, apiKey, *sFlag, question)
	fmt.Println(answer)
}
