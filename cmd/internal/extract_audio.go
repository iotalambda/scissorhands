package internal

import (
	"scissorhands/ffmpeg"

	"github.com/spf13/cobra"
)

var ExtractAudioCmd = &cobra.Command{
	Use:   "extract-audio",
	Short: "Extract audio from an input file",
	RunE: func(cmd *cobra.Command, args []string) error {
		return ffmpeg.ExtractAudio(input, output)
	},
}

func init() {
	ExtractAudioCmd.Flags().StringVarP(&input, "input", "i", "", "Input file path.")
	ExtractAudioCmd.MarkFlagRequired("input")

	ExtractAudioCmd.Flags().StringVarP(&output, "output", "o", "", "Output file path.")
	ExtractAudioCmd.MarkFlagRequired("output")
}
