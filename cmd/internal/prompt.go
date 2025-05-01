package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"scissorhands/config"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/spf13/cobra"
)

var PromptCmd = &cobra.Command{
	Use:   "prompt",
	Short: "Let an LLM decide the actions to take",
	RunE: func(cmd *cobra.Command, args []string) error {
		return prompt()
	},
}

func prompt() error {
	switch service {
	case "openai-gpt4o":
		if err := promptWithOpenAIGPT4o(); err != nil {
			return fmt.Errorf("promp with OpenAI GPT-4o: %v", err)
		}
	default:
		return fmt.Errorf("unrecognized service: %v", service)
	}
	return nil
}

func promptWithOpenAIGPT4o() error {
	client := openai.NewClient(option.WithAPIKey(config.Global.OpenAIApiKey))
	ctx := context.Background()
	question := "Can you list the files and directories in the `/workspaces` directory, please."
	fmt.Println("YOU> " + question)
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage("You are a helpful agent that can call `ls` locally and tell the user information about their local files and directories."),
		openai.UserMessage(question),
	}

	newParams := func() openai.ChatCompletionNewParams {
		return openai.ChatCompletionNewParams{
			Messages: messages,
			Tools: []openai.ChatCompletionToolParam{
				{
					Function: openai.FunctionDefinitionParam{
						Name:        "ls",
						Description: openai.String("Runs the plain old `ls` in a given location. The result is the stdout."),
						Parameters: openai.FunctionParameters{
							"type": "object",
							"properties": map[string]any{
								"location": map[string]any{
									"type": "string",
								},
							},
							"required": []string{"location"},
						},
					},
				},
			},
			Model: openai.ChatModelGPT4o,
			Seed:  openai.Int(0),
		}
	}

	// Initial chat completion
	res1, err := client.Chat.Completions.New(ctx, newParams())
	if err != nil {
		return fmt.Errorf("initial chat completion: %v", err)
	}

	// Return early if no tool calls
	res1Msg := res1.Choices[0].Message
	messages = append(messages, res1Msg.ToParam())
	if len(res1Msg.ToolCalls) == 0 {
		fmt.Println("No function calls.")
		return nil
	}

	// Handle tool calls
	for _, toolCall := range res1Msg.ToolCalls {
		if toolCall.Function.Name == "ls" {

			// Parse ls tool call args
			var args map[string]any
			err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args)
			if err != nil {
				return fmt.Errorf("unmarshal ls tool call args: %v", err)
			}
			location := args["location"].(string)

			// Exec ls tool call
			cmd := exec.Command("ls", location)
			b, err := cmd.Output()
			if err != nil {
				return fmt.Errorf("ls command exec: %v", err)
			}

			messages = append(messages, openai.ToolMessage(string(b), toolCall.ID))

			// Send tool call result back
			res2, err := client.Chat.Completions.New(ctx, newParams())
			if err != nil {
				return fmt.Errorf("send ls tool call results back: %v", err)
			}
			res2Msg := res2.Choices[0].Message
			messages = append(messages, res2Msg.ToParam())

			fmt.Println("LLM> " + res2Msg.Content)
		}
	}

	return nil
}

func init() {
	PromptCmd.Flags().StringVarP(&input, "input", "i", "", "Input file path.")
	PromptCmd.MarkFlagRequired("input")

	PromptCmd.Flags().StringVarP(&output, "output", "o", "", "Output file path.")
	PromptCmd.MarkFlagRequired("output")

	PromptCmd.Flags().StringVarP(&message, "message", "m", "", "Prompt message.")
	PromptCmd.MarkFlagRequired("message")

	PromptCmd.Flags().StringVarP(&service, "service", "s", "", "The LLM to use. Allowed values: openai-gpt4o.")
	PromptCmd.MarkFlagRequired("service")
}
