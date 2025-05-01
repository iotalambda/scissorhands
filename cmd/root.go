package cmd

import (
	"log"
	internal "scissorhands/cmd/internal"
	"scissorhands/stuff"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "scissorhands",
	Short: "A CLI tool for smart video editing",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return stuff.InitConfig()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}

func init() {
	rootCmd.AddCommand(inferDialogueCmd)

	rootCmd.AddCommand(internal.AzureSpeechDiarizeCmd)
	rootCmd.AddCommand(internal.ExtractAudioCmd)
	rootCmd.AddCommand(internal.OpenAIWhisperSegmentCmd)
	rootCmd.AddCommand(internal.PromptCmd)
	rootCmd.AddCommand(internal.ScreenshotCmd)
}
