package cmd

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/david22573/openrouter-api-go/internal/app"
	"github.com/david22573/openrouter-api-go/pkg/openrouter"
	"github.com/spf13/cobra"
)

var defaultModel string

var modelID string

// chatCmd represents the chat command
var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Start an interactive streaming chat session with OpenRouter",
	Long: `Start an interactive chat session. The API key must be set via the
OPENROUTER_API_KEY environment variable. You can specify a model using the --model flag.

Example:
  export OPENROUTER_API_KEY="sk-..."
  ./mycli chat --model nousresearch/nous-hermes-2-mixtral-8x7b-dpo`,
	Run: runChat,
}

func init() {
	// will be replaced in PersistentPreRunE
	defaultModel = "tngtech/tng-r1t-chimera:free"

	rootCmd.AddCommand(chatCmd)
	chatCmd.Flags().StringVarP(&modelID, "model", "m", "", "Model ID")
}

func runChat(cmd *cobra.Command, args []string) {
	// 1. Determine model by priority
	//    flag → config → hard default
	if cmd.Flags().Changed("model") {
		// modelID already set
	} else if app.A.Config.Model != "" {
		modelID = app.A.Config.Model
	} else {
		modelID = "tngtech/tng-r1t-chimera:free"
	}

	fmt.Println("Using model:", modelID)

	client := app.A.Client // Already initialized in root.go
	// Initialize chat history

	messages := make([]openrouter.ChatMessage, 0)

	fmt.Printf("Starting chat with model: %s\n", modelID)
	fmt.Println("Type 'exit' or 'quit' to end the session.")
	fmt.Print("-------------------------------------------------------------------\n")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("\n> You: ")
		if !scanner.Scan() {
			break // EOF or error
		}
		userInput := scanner.Text()

		if strings.ToLower(userInput) == "exit" || strings.ToLower(userInput) == "quit" {
			break
		}
		if strings.TrimSpace(userInput) == "" {
			continue
		}

		// 1. Add user message to history
		messages = append(messages, openrouter.ChatMessage{
			Role:    "user",
			Content: userInput,
		})

		// 2. Prepare request for streaming
		req := openrouter.ChatCompletionRequest{
			Model:       modelID,
			Messages:    messages,
			Temperature: floatPtr(0.7),
			MaxTokens:   4096,
		}

		ctx := context.Background()
		stream, err := client.CreateChatCompletionStream(ctx, req)
		if err != nil {
			fmt.Printf("\n[API Error]: %v\n", err)
			// Remove the failed user message from history
			messages = messages[:len(messages)-1]
			continue
		}
		defer stream.Close()

		fmt.Print("\n< AI: ")

		var fullResponseContent strings.Builder

		// 3. Process the streaming response
		for {
			resp, err := stream.Recv()
			if err != nil {
				if err != io.EOF {
					fmt.Printf("\n[Stream Error]: %v\n", err)
				}
				break
			}

			// In a streaming response, the content is in the Delta field
			if len(resp.Choices) > 0 && resp.Choices[0].Delta != nil {
				content := resp.Choices[0].Delta.Content // Content is string for simple text

				// Handle both string and array content if necessary, but typically
				// streamed text is a string from OpenRouter/OpenAI compatible APIs.
				if contentStr, ok := content.(string); ok && contentStr != "" {
					fmt.Print(contentStr)
					fullResponseContent.WriteString(contentStr)
				}
			}
		}

		// 4. Add the full AI response to history
		aiContent := fullResponseContent.String()
		if aiContent != "" {
			messages = append(messages, openrouter.ChatMessage{
				Role:    "assistant",
				Content: aiContent,
			})
		}

		fmt.Println() // Newline after AI response is complete
	}

	fmt.Println("\nChat session ended.")
}

// Helper function to get a pointer to a float32
func floatPtr(f float32) *float32 {
	return &f
}
