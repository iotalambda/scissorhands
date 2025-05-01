package internal

import (
	"scissorhands/stuff"

	"github.com/spf13/cobra"
)

var AzureSpeechDiarizeCmd = &cobra.Command{
	Use:   "azure-speech-diarize",
	Short: "Diarize an audio file by speakers",
	RunE: func(cmd *cobra.Command, args []string) error {
		return stuff.AzureSpeechDiarize(input, output, maxSpeakers)
	},
}

func init() {
	AzureSpeechDiarizeCmd.Flags().StringVarP(&input, "input", "i", "", "Input file path.")
	AzureSpeechDiarizeCmd.MarkFlagRequired("input")

	AzureSpeechDiarizeCmd.Flags().StringVarP(&output, "output", "o", "", "Output file path.")
	AzureSpeechDiarizeCmd.MarkFlagRequired("output")

	AzureSpeechDiarizeCmd.Flags().IntVarP(&maxSpeakers, "max-speakers", "m", 5, "Maximum number of possible speakers to evaluate.")
}
