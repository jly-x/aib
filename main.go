package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

func initClient(ctx context.Context) *genai.Client {
	key := os.Getenv("GOOGLE_API_KEY")

	if key == "" {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("The environment variable GOOGLE_API_KEY is not set, and there was an unexpected error loading .env file: %v", err)
		}
		key = os.Getenv("GOOGLE_API_KEY")
	}

	if key == "" {
		log.Fatalf("GOOGLE_API_KEY environment variable not set")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(key))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	return client
}

func generateAnswer(client *genai.Client, ctx context.Context, s bool, q string) string {
	model := client.GenerativeModel("gemini-2.0-flash")

	var systemPrompt string
	if s {
		systemPrompt = "You are a helpful and versatile assistant, especially skilled in programming and devops. The user is asking for a bash terminal command. Respond ONLY with the exact, executable bash command that accomplishes what they're asking for. Do not include explanations, markdown formatting, unnecessary quotation marks, or any text other than the command itself. The command will be automatically executed in the user's terminal."
	} else {
		systemPrompt = "You are a helpful and versatile assistant, especially skilled in programming and devops. Answer the user's question(s) in a very concise manner, unless it requires a long response or the user specifically asks you for one."
	}

	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(systemPrompt)},
		Role:  "system",
	}

	resp, err := model.GenerateContent(ctx, genai.Text(q))
	if err != nil {
		log.Fatalf("Failed to generate content: %v", err)
	}

	if len(resp.Candidates) > 0 && resp.Candidates[0].Content != nil && len(resp.Candidates[0].Content.Parts) > 0 {
		result := fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0])

		if s {
			// remove trailing whitespace
			result = strings.TrimSpace(result)
			err := clipboard.WriteAll(result)
			if err != nil {
				log.Fatalf("Failed to copy to clipboard: %v", err)
			}
			result = fmt.Sprintf("%v\nCopied to clipboard", result)
		}

		return result
	}

	log.Fatal("No valid response")
	return ""
}

func main() {
	ctx := context.Background()
	client := initClient(ctx)

	sFlag := flag.Bool("s", false, "Returns a valid bash prompt that matches the request.")

	flag.Parse()

	args := flag.Args()

	if len(args) == 0 {
		err := fmt.Errorf("no arguments provided")
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	question := strings.Join(args, " ")

	answer := generateAnswer(client, ctx, *sFlag, question)
	fmt.Println(answer)
}
